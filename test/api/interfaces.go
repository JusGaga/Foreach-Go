package api

import (
	"time"

	"github.com/jusgaga/wordmon-go/internal/core"
)

// PlayerStore définit l'interface pour le stockage des joueurs
type PlayerStore interface {
	CreatePlayer(name string) (*PlayerResponse, error)
	GetPlayer(id string) (*PlayerResponse, error)
	GetAllPlayers() []*PlayerResponse
	UpdatePlayer(player *PlayerResponse) error
	GetPlayerCount() int
	GetStartTime() time.Time
	UpdateXP(id string, newXP int, newLevel int) error
}

// SpawnStore définit l'interface pour le stockage des spawns
type SpawnStore interface {
	AddSpawn(spawn interface{}) error
	GetCurrentSpawn() interface{}
}

// WordStore définit l'interface pour le stockage des mots
type WordStore interface {
	Seed(words []core.Word) error
	Get(id string) (*core.Word, error)
	RandomByRarity(rarity string) (*core.Word, error)
}

// CaptureStore définit l'interface pour le stockage des captures
type CaptureStore interface {
	Add(playerId, wordId string) error
	ListByPlayer(playerId string) ([]core.Word, error)
}

// LeaderboardStore définit l'interface pour le leaderboard
type LeaderboardStore interface {
	GetLeaderboard(limit int) ([]*PlayerResponse, error)
}
