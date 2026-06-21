package main

import (
	"sort"
	"sync"
	"time"
)

type UploadEntry struct {
	FileName   string    `json:"file_name"`
	Folder     string    `json:"folder"`
	Size       int64     `json:"size"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type UploadStore struct {
	mu      sync.RWMutex
	entries []UploadEntry
}

func NewUploadStore() *UploadStore {
	return &UploadStore{}
}

func (s *UploadStore) Add(entry UploadEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, entry)
}

func (s *UploadStore) List() []UploadEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]UploadEntry, len(s.entries))
	copy(result, s.entries)
	sort.Slice(result, func(i, j int) bool {
		return result[i].UploadedAt.After(result[j].UploadedAt)
	})
	return result
}
