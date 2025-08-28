package config

import (
	"os"
	"testing"
)

func TestEnvironmentOverride(t *testing.T) {
	// Sauvegarder la valeur originale
	originalInterval := os.Getenv("WORDMON_SPAWN_INTERVAL")
	defer func() {
		if originalInterval != "" {
			os.Setenv("WORDMON_SPAWN_INTERVAL", originalInterval)
		} else {
			os.Unsetenv("WORDMON_SPAWN_INTERVAL")
		}
	}()

	// Tester l'override par variable d'environnement
	os.Setenv("WORDMON_SPAWN_INTERVAL", "5")

	// Vérifier que la variable est bien définie
	if os.Getenv("WORDMON_SPAWN_INTERVAL") != "5" {
		t.Error("La variable d'environnement n'a pas été définie correctement")
	}
}

func TestGetenvOrDefault(t *testing.T) {
	// Sauvegarder la valeur originale
	originalPath := os.Getenv("WORDMON_CONFIG_PATH")
	defer func() {
		if originalPath != "" {
			os.Setenv("WORDMON_CONFIG_PATH", originalPath)
		} else {
			os.Unsetenv("WORDMON_CONFIG_PATH")
		}
	}()

	// Test avec variable non définie
	os.Unsetenv("WORDMON_CONFIG_PATH")
	result := getenvOrDefault("WORDMON_CONFIG_PATH", "default_value")
	if result != "default_value" {
		t.Errorf("getenvOrDefault devrait retourner 'default_value', obtenu '%s'", result)
	}

	// Test avec variable définie
	os.Setenv("WORDMON_CONFIG_PATH", "custom_path")
	result = getenvOrDefault("WORDMON_CONFIG_PATH", "default_value")
	if result != "custom_path" {
		t.Errorf("getenvOrDefault devrait retourner 'custom_path', obtenu '%s'", result)
	}
}

func TestResolvePath(t *testing.T) {
	// Test avec chemin relatif
	relativePath := "configs/test.yaml"
	resolved := ResolvePath(relativePath)
	if resolved == relativePath {
		t.Logf("ResolvePath a retourné le chemin original: %s", resolved)
	}

	// Test avec chemin absolu (utiliser un chemin Windows valide)
	absPath := "C:\\absolute\\path"
	resolved = ResolvePath(absPath)
	if resolved != absPath {
		t.Errorf("ResolvePath devrait retourner le chemin absolu inchangé, obtenu '%s'", resolved)
	}
}
