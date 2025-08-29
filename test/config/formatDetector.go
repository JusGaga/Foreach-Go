package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type fileFormat int

const (
	fmtUnknown fileFormat = iota
	fmtYAML
	fmtTOML
	fmtJSON
)

func detectFormat(path string) fileFormat {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		return fmtYAML
	case ".toml":
		return fmtTOML
	case ".json":
		return fmtJSON
	default:
		return fmtUnknown
	}
}

func decodeFile(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("impossible de lire le fichier %s: %w", path, err)
	}
	switch detectFormat(path) {
	case fmtYAML:
		if err := yaml.Unmarshal(data, v); err != nil {
			return fmt.Errorf("erreur de parsing YAML (%s): %w", path, err)
		}
	case fmtTOML:
		if err := toml.Unmarshal(data, v); err != nil {
			return fmt.Errorf("erreur de parsing TOML (%s): %w", path, err)
		}
	case fmtJSON:
		if err := json.Unmarshal(data, v); err != nil {
			return fmt.Errorf("erreur de parsing JSON (%s): %w", path, err)
		}
	default:
		return &UnsupportedFormatError{Path: path}
	}
	return nil
}

func mustBeYAMLorTOML(path string) error {
	switch detectFormat(path) {
	case fmtYAML, fmtTOML:
		return nil
	default:
		if detectFormat(path) == fmtJSON {
			return errors.New("JSON non supporté pour ce type de configuration, utilisez YAML ou TOML")
		}
		return &UnsupportedFormatError{Path: path}
	}
}

func mustBeJSON(path string) error {
	if detectFormat(path) != fmtJSON {
		return errors.New("le dictionnaire de mots doit être au format JSON (.json)")
	}
	return nil
}

func isAllowedRarity(r string) bool {
	_, ok := AllowedRarities[r]
	return ok
}

func getenvOrDefault(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
