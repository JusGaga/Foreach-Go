package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *Handlers) {
	gin.SetMode(gin.TestMode)

	// Créer un store de test
	playerStore := &MockPlayerStore{
		players: make(map[string]*PlayerResponse),
	}
	spawnStore := &MockSpawnStore{
		spawns: make([]interface{}, 0),
	}

	handlers := NewHandlers(playerStore, spawnStore)

	router := gin.New()
	router.POST("/players", handlers.CreatePlayer)
	router.GET("/players/:id", handlers.GetPlayer)
	router.GET("/spawn/current", handlers.GetCurrentSpawn)
	router.POST("/encounter/attempt", handlers.AttemptCapture)
	router.GET("/leaderboard", handlers.GetLeaderboard)
	router.GET("/status", handlers.GetStatus)

	return router, handlers
}

// MockPlayerStore pour les tests
type MockPlayerStore struct {
	players map[string]*PlayerResponse
}

func (m *MockPlayerStore) CreatePlayer(name string) (*PlayerResponse, error) {
	player := &PlayerResponse{
		ID:        "p1",
		Name:      name,
		XP:        0,
		Level:     1,
		Inventory: make(map[string]int),
	}
	m.players["p1"] = player
	return player, nil
}

func (m *MockPlayerStore) GetPlayer(id string) (*PlayerResponse, error) {
	if player, exists := m.players[id]; exists {
		return player, nil
	}
	return nil, &PlayerNotFoundError{ID: id}
}

func (m *MockPlayerStore) GetAllPlayers() []*PlayerResponse {
	players := make([]*PlayerResponse, 0, len(m.players))
	for _, player := range m.players {
		players = append(players, player)
	}
	return players
}

func (m *MockPlayerStore) UpdatePlayer(player *PlayerResponse) error {
	m.players[player.ID] = player
	return nil
}

func (m *MockPlayerStore) GetPlayerCount() int {
	return len(m.players)
}

func (m *MockPlayerStore) GetStartTime() time.Time {
	return time.Now()
}

// MockSpawnStore pour les tests
type MockSpawnStore struct {
	spawns []interface{}
}

func (m *MockSpawnStore) AddSpawn(spawn interface{}) error {
	m.spawns = append(m.spawns, spawn)
	return nil
}

func (m *MockSpawnStore) GetCurrentSpawn() interface{} {
	if len(m.spawns) == 0 {
		return nil
	}
	return m.spawns[len(m.spawns)-1]
}

func TestCreatePlayer(t *testing.T) {
	router, _ := setupTestRouter()

	// Test création réussie
	playerData := CreatePlayerRequest{Name: "Sacha"}
	jsonData, _ := json.Marshal(playerData)

	req, _ := http.NewRequest("POST", "/players", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PlayerResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Sacha", response.Name)
	assert.Equal(t, 1, response.Level)
}

func TestGetPlayer(t *testing.T) {
	router, _ := setupTestRouter()

	// Créer un joueur d'abord
	playerData := CreatePlayerRequest{Name: "Sacha"}
	jsonData, _ := json.Marshal(playerData)

	req, _ := http.NewRequest("POST", "/players", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Maintenant récupérer le joueur
	req, _ = http.NewRequest("GET", "/players/p1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PlayerResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Sacha", response.Name)
}

func TestGetStatus(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest("GET", "/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response StatusResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "WordMon Go", response.Game)
	assert.Equal(t, "0.3.0", response.Version)
}
