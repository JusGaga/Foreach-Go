package api

import "time"

// StatusResponse représente la réponse de l'endpoint /status
type StatusResponse struct {
	Game          string     `json:"game"`
	Version       string     `json:"version"`
	UptimeSeconds int64      `json:"uptimeSeconds"`
	ActivePlayers int        `json:"activePlayers"`
	CurrentSpawn  *SpawnInfo `json:"currentSpawn"`
}

// SpawnInfo représente les informations d'un spawn actif
type SpawnInfo struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Rarity string `json:"rarity"`
	Points int    `json:"points"`
}

// CreatePlayerRequest représente la requête pour créer un joueur
type CreatePlayerRequest struct {
	Name string `json:"name" binding:"required"`
}

// PlayerResponse représente la réponse pour un joueur
type PlayerResponse struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	XP        int            `json:"xp"`
	Level     int            `json:"level"`
	Inventory map[string]int `json:"inventory"`
}

// CaptureAttemptRequest représente la requête pour tenter une capture
type CaptureAttemptRequest struct {
	PlayerID string `json:"playerId" binding:"required"`
	Attempt  string `json:"attempt" binding:"required"`
}

// CaptureResultResponse représente le résultat d'une tentative de capture
type CaptureResultResponse struct {
	Status   string `json:"status"`
	Word     string `json:"word,omitempty"`
	Rarity   string `json:"rarity,omitempty"`
	XP       int    `json:"xp,omitempty"`
	NewLevel int    `json:"newLevel,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

// LeaderboardEntry représente une entrée du leaderboard
type LeaderboardEntry struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	XP    int    `json:"xp"`
	Level int    `json:"level"`
}

// ErrorResponse représente une réponse d'erreur
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// GameState représente l'état global du jeu
type GameState struct {
	StartTime     time.Time
	Players       map[string]*PlayerResponse
	CurrentSpawn  *SpawnInfo
	PlayerCounter int
}
