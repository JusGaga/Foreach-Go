package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jusgaga/wordmon-go/internal/core"
)

const (
	envSnapshotPath = "WORDMON_SNAPSHOT_PATH"
	DefaultSnapshotPath = "data/snapshot.json"
)

// Snapshot représente l'état sauvegardé du jeu
type Snapshot struct {
	UpdatedAt time.Time           `json:"updatedAt"`
	Players   []PlayerSnapshot    `json:"players"`
}

// PlayerSnapshot représente l'état sauvegardé d'un joueur
type PlayerSnapshot struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	XP        int            `json:"xp"`
	Level     int            `json:"level"`
	Inventory map[string]int `json:"inventory"`
}

// SaveSnapshot sauvegarde l'état des joueurs dans un fichier JSON
func SaveSnapshot(players []core.Player, path string) error {
	if env := os.Getenv(envSnapshotPath); env != "" {
		path = env
	}
	if path == "" {
		path = DefaultSnapshotPath
	}

	// Créer le répertoire parent si nécessaire
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("impossible de créer le répertoire %s: %w", dir, err)
	}

	// Créer le snapshot
	snapshot := Snapshot{
		UpdatedAt: time.Now().UTC(),
		Players:   make([]PlayerSnapshot, 0, len(players)),
	}

	// Convertir les joueurs
	for _, p := range players {
		ps := PlayerSnapshot{
			ID:        p.ID,
			Name:      p.Name,
			XP:        p.XP,
			Level:     p.Level,
			Inventory: make(map[string]int),
		}
		
		// Copier l'inventaire
		for word, count := range p.Inventory {
			ps.Inventory[word] = count
		}
		
		snapshot.Players = append(snapshot.Players, ps)
	}

	// Écrire d'abord dans un fichier temporaire pour éviter la corruption
	tempPath := path + ".tmp"
	
	// Encoder en JSON avec indentation
	data, err := json.MarshalIndent(snapshot, "", " ")
	if err != nil {
		return fmt.Errorf("erreur d'encodage JSON: %w", err)
	}

	// Écrire le fichier temporaire
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("erreur d'écriture du fichier temporaire %s: %w", tempPath, err)
	}

	// Renommer atomiquement le fichier temporaire
	if err := os.Rename(tempPath, path); err != nil {
		// Nettoyer le fichier temporaire en cas d'échec
		os.Remove(tempPath)
		return fmt.Errorf("erreur de renommage atomique vers %s: %w", path, err)
	}

	return nil
}

// LoadSnapshot charge un snapshot depuis un fichier JSON
func LoadSnapshot(path string) (*Snapshot, error) {
	if env := os.Getenv(envSnapshotPath); env != "" {
		path = env
	}
	if path == "" {
		path = DefaultSnapshotPath
	}

	// Vérifier que le fichier existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("fichier snapshot non trouvé: %s", path)
	}

	// Lire le fichier
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("erreur de lecture du snapshot %s: %w", path, err)
	}

	// Décoder le JSON
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("erreur de décodage JSON du snapshot: %w", err)
	}

	return &snapshot, nil
}
