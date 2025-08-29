package api

import (
	"sync"
	"time"

	"github.com/jusgaga/wordmon-go/internal/core"
)

// SimpleStore implémente un stockage en mémoire simple
type SimpleStore struct {
	mu            sync.RWMutex
	players       map[string]*PlayerResponse
	spawns        []core.SpawnEvent
	startTime     time.Time
	playerCounter int
}

// NewSimpleStore crée un nouveau store simple
func NewSimpleStore() *SimpleStore {
	return &SimpleStore{
		players:       make(map[string]*PlayerResponse),
		spawns:        make([]core.SpawnEvent, 0),
		startTime:     time.Now(),
		playerCounter: 0,
	}
}

// CreatePlayer crée un nouveau joueur
func (s *SimpleStore) CreatePlayer(name string) (*PlayerResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Vérifier si le nom est déjà pris
	for _, player := range s.players {
		if player.Name == name {
			return nil, &PlayerNameTakenError{Name: name}
		}
	}

	s.playerCounter++
	playerID := "p" + string(rune(s.playerCounter+'0'))

	player := &PlayerResponse{
		ID:        playerID,
		Name:      name,
		XP:        0,
		Level:     1,
		Inventory: make(map[string]int),
	}

	s.players[playerID] = player
	return player, nil
}

// GetPlayer récupère un joueur par son ID
func (s *SimpleStore) GetPlayer(id string) (*PlayerResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	player, exists := s.players[id]
	if !exists {
		return nil, &PlayerNotFoundError{ID: id}
	}

	return player, nil
}

// GetAllPlayers récupère tous les joueurs
func (s *SimpleStore) GetAllPlayers() []*PlayerResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	players := make([]*PlayerResponse, 0, len(s.players))
	for _, player := range s.players {
		players = append(players, player)
	}

	return players
}

// UpdatePlayer met à jour un joueur
func (s *SimpleStore) UpdatePlayer(player *PlayerResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.players[player.ID]; !exists {
		return &PlayerNotFoundError{ID: player.ID}
	}

	s.players[player.ID] = player
	return nil
}

// UpdateXP met à jour l'XP et le niveau d'un joueur
func (s *SimpleStore) UpdateXP(id string, newXP int, newLevel int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	player, exists := s.players[id]
	if !exists {
		return &PlayerNotFoundError{ID: id}
	}

	player.XP = newXP
	player.Level = newLevel
	return nil
}

// AddSpawn ajoute un spawn à l'historique
func (s *SimpleStore) AddSpawn(spawn interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if spawnEvent, ok := spawn.(core.SpawnEvent); ok {
		s.spawns = append(s.spawns, spawnEvent)
		return nil
	}
	return &InvalidSpawnError{}
}

// GetCurrentSpawn récupère le spawn actuel
func (s *SimpleStore) GetCurrentSpawn() interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.spawns) == 0 {
		return nil
	}

	return s.spawns[len(s.spawns)-1]
}

// GetStartTime retourne l'heure de démarrage
func (s *SimpleStore) GetStartTime() time.Time {
	return s.startTime
}

// GetPlayerCount retourne le nombre de joueurs
func (s *SimpleStore) GetPlayerCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.players)
}

// PlayerNameTakenError erreur quand le nom est déjà pris
type PlayerNameTakenError struct {
	Name string
}

func (e *PlayerNameTakenError) Error() string {
	return "nom déjà pris: " + e.Name
}

// PlayerNotFoundError erreur quand le joueur n'est pas trouvé
type PlayerNotFoundError struct {
	ID string
}

func (e *PlayerNotFoundError) Error() string {
	return "joueur non trouvé: " + e.ID
}

// InvalidSpawnError erreur quand le spawn est invalide
type InvalidSpawnError struct{}

func (e *InvalidSpawnError) Error() string {
	return "spawn invalide"
}
