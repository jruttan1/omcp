package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Spec    string `json:"spec"`
	BaseURL string `json:"baseURL"`
	APIKey  string `json:"apiKey"`
	Tools   []Tool `json:"tools"`
	Info    Info   `json:"info"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".omcp", "config.json"), nil
}

func saveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}

func loadConfig() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()
	var cfg Config
	return cfg, json.NewDecoder(f).Decode(&cfg)
}
