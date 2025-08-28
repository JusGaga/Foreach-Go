package config

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Section string
	Problems []string
}

func (e *ValidationError) Error() string {
	if len(e.Problems) == 0 {
		return fmt.Sprintf("configuration invalide dans %s", e.Section)
	}
	return fmt.Sprintf("configuration invalide dans %s: %s", e.Section, strings.Join(e.Problems, "; "))
}

func newValidationError(section string) *ValidationError {
	return &ValidationError{Section: section, Problems: []string{}}
}

func (e *ValidationError) addf(format string, args ...any) {
	e.Problems = append(e.Problems, fmt.Sprintf(format, args...))
}

func (e *ValidationError) ok() bool {
	return len(e.Problems) == 0
}

type UnsupportedFormatError struct {
	Path string
}

func (e *UnsupportedFormatError) Error() string {
	return fmt.Sprintf("format non supporté pour le fichier: %s (extensions supportées: .yaml, .yml, .toml, .json uniquement pour les mots)", e.Path)
}

type EnvOverrideError struct {
	Var   string
	Value string
	Reason string
}

func (e *EnvOverrideError) Error() string {
	return fmt.Sprintf("valeur invalide pour %s='%s': %s", e.Var, e.Value, e.Reason)
}
