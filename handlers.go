package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Handler struct {
	config *Config
	store  *UploadStore
}

func NewHandler(cfg *Config, store *UploadStore) *Handler {
	return &Handler{config: cfg, store: store}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, h.config); err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}

func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := make([]FolderResponse, len(h.config.Folders))
	for i, f := range h.config.Folders {
		resp[i] = FolderResponse{Label: f.Label}
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("error encoding config response: %v", err)
	}
}

func (h *Handler) ListUploads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(h.store.List()); err != nil {
		log.Printf("error encoding uploads response: %v", err)
	}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	label := r.PathValue("label")
	if label == "" {
		http.Error(w, "Missing folder label", http.StatusBadRequest)
		return
	}

	folder := h.config.FindFolder(label)
	if folder == nil {
		http.Error(w, fmt.Sprintf("Unknown folder: %s", label), http.StatusNotFound)
		return
	}

	maxBytes := int64(h.config.Server.MaxFileSize)
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	if err := r.ParseMultipartForm(maxBytes); err != nil {
		http.Error(w, "File too large or invalid form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file in request", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing uploaded file: %v", err)
		}
	}()

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".torrent") {
		http.Error(w, "Only .torrent files are allowed", http.StatusBadRequest)
		return
	}

	filename := sanitizeFilename(header.Filename)
	destPath := filepath.Join(folder.Path, filename)

	absDest, err := filepath.Abs(destPath)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	absFolder, err := filepath.Abs(folder.Path)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if !strings.HasPrefix(absDest, absFolder) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	if err := os.MkdirAll(folder.Path, 0755); err != nil {
		http.Error(w, "Cannot create destination directory", http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Cannot create destination file", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := dst.Close(); err != nil {
			log.Printf("error closing destination file: %v", err)
		}
	}()

	written, err := io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	h.store.Add(UploadEntry{
		FileName:   filename,
		Folder:     label,
		Size:       written,
		UploadedAt: time.Now(),
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"success":  true,
		"filename": filename,
		"folder":   label,
		"size":     written,
	}); err != nil {
		log.Printf("error encoding upload response: %v", err)
	}
}

func sanitizeFilename(name string) string {
	name = path.Base(name)
	name = strings.ReplaceAll(name, "..", "")
	name = strings.TrimSpace(name)
	if name == "" || name == "." || name == "/" {
		return "unnamed.torrent"
	}
	return name
}
