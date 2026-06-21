package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	configPath := "config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	store := NewUploadStore()
	handler := NewHandler(cfg, store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handler.Index)
	mux.HandleFunc("GET /api/config", handler.GetConfig)
	mux.HandleFunc("GET /api/uploads", handler.ListUploads)
	mux.HandleFunc("POST /api/upload/{label}", handler.Upload)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
