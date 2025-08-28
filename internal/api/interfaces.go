package api

import "time"

// PlayerStore définit l'interface pour le stockage des joueurs
type PlayerStore interface {
	CreatePlayer(name string) (*PlayerResponse, error)
	GetPlayer(id string) (*PlayerResponse, error)
	GetAllPlayers() []*PlayerResponse
	UpdatePlayer(player *PlayerResponse) error
	GetPlayerCount() int
	GetStartTime() time.Time
}

// SpawnStore définit l'interface pour le stockage des spawns
type SpawnStore interface {
	AddSpawn(spawn interface{}) error
	GetCurrentSpawn() interface{}
}
