package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type ByteSize int64

func (b *ByteSize) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}
	size, err := ParseByteSize(s)
	if err != nil {
		return err
	}
	*b = size
	return nil
}

func ParseByteSize(s string) (ByteSize, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty size")
	}

	upper := strings.ToUpper(s)

	units := []struct {
		suffix string
		mult   int64
	}{
		{"TIB", tebibyte}, {"GIB", gibibyte}, {"MIB", mebibyte}, {"KIB", kibibyte},
		{"TB", terabyte}, {"GB", gigabyte}, {"MB", megabyte}, {"KB", kilobyte},
		{"T", tebibyte}, {"G", gibibyte}, {"M", mebibyte}, {"K", kibibyte},
		{"B", 1},
	}

	for _, u := range units {
		if strings.HasSuffix(upper, u.suffix) {
			numStr := s[:len(s)-len(u.suffix)]
			numStr = strings.TrimSpace(numStr)
			if numStr == "" {
				return 0, fmt.Errorf("invalid size: %s", s)
			}
			val, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid size: %s", s)
			}
			return ByteSize(val * float64(u.mult)), nil
		}
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size: %s (expected format like 50MB, 1GiB, 500KB)", s)
	}
	return ByteSize(val), nil
}

const (
	kilobyte = 1000
	megabyte = 1000 * kilobyte
	gigabyte = 1000 * megabyte
	terabyte = 1000 * gigabyte
	kibibyte = 1024
	mebibyte = 1024 * kibibyte
	gibibyte = 1024 * mebibyte
	tebibyte = 1024 * gibibyte
)

type Config struct {
	Server  ServerConfig   `yaml:"server"`
	Folders []FolderConfig `yaml:"folders"`
}

type ServerConfig struct {
	Host        string   `yaml:"host"`
	Port        int      `yaml:"port"`
	MaxFileSize ByteSize `yaml:"max_file_size"`
}

type FolderConfig struct {
	Label string `yaml:"label"`
	Path  string `yaml:"path"`
}

type FolderResponse struct {
	Label string `json:"label"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config: %w", err)
	}

	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.MaxFileSize == 0 {
		cfg.Server.MaxFileSize = 50 * mebibyte
	}

	return &cfg, nil
}

func (c *Config) FindFolder(label string) *FolderConfig {
	for _, f := range c.Folders {
		if f.Label == label {
			return &f
		}
	}
	return nil
}
