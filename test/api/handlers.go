package api

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jusgaga/wordmon-go/internal/core"
)

// Handlers contient tous les gestionnaires d'endpoints
type Handlers struct {
	playerStore PlayerStore
	spawnStore  SpawnStore
	spawner     chan core.SpawnEvent
}

// NewHandlers crée une nouvelle instance de Handlers
func NewHandlers(playerStore PlayerStore, spawnStore SpawnStore) *Handlers {
	return &Handlers{
		playerStore: playerStore,
		spawnStore:  spawnStore,
		spawner:     make(chan core.SpawnEvent, 1),
	}
}

// SetSpawner définit le canal de spawn pour les Handlers
func (h *Handlers) SetSpawner(spawner chan core.SpawnEvent) {
	h.spawner = spawner
}

// GetStatus retourne le statut du serveur
func (h *Handlers) GetStatus(c *gin.Context) {
	uptime := time.Since(h.playerStore.GetStartTime()).Seconds()

	var currentSpawn *SpawnInfo
	if spawn := h.spawnStore.GetCurrentSpawn(); spawn != nil {
		if spawnEvent, ok := spawn.(core.SpawnEvent); ok {
			currentSpawn = &SpawnInfo{
				ID:     spawnEvent.Word.ID,
				Text:   spawnEvent.Word.Text,
				Rarity: string(spawnEvent.Word.Rarity),
				Points: spawnEvent.Word.Points,
			}
		}
	}

	response := StatusResponse{
		Game:          "WordMon Go",
		Version:       "0.3.0",
		UptimeSeconds: int64(uptime),
		ActivePlayers: h.playerStore.GetPlayerCount(),
		CurrentSpawn:  currentSpawn,
	}

	c.JSON(http.StatusOK, response)
}

// CreatePlayer crée un nouveau joueur
func (h *Handlers) CreatePlayer(c *gin.Context) {
	var req CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Nom du joueur requis",
		})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_name",
			Message: "Le nom du joueur ne peut pas être vide",
		})
		return
	}

	// Utiliser le store pour créer le joueur
	player, err := h.playerStore.CreatePlayer(req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "name_taken",
			Message: "Ce nom est déjà pris",
		})
		return
	}

	c.JSON(http.StatusOK, player)
}

// GetPlayer retourne les informations d'un joueur
func (h *Handlers) GetPlayer(c *gin.Context) {
	playerID := c.Param("id")

	player, err := h.playerStore.GetPlayer(playerID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "player_not_found",
			Message: "Joueur non trouvé",
		})
		return
	}

	c.JSON(http.StatusOK, player)
}

// GetCurrentSpawn retourne le spawn actuel
func (h *Handlers) GetCurrentSpawn(c *gin.Context) {
	spawn := h.spawnStore.GetCurrentSpawn()
	if spawn == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "no_spawn",
			Message: "Aucun WordMon actif",
		})
		return
	}

	// Convertir le spawn en SpawnInfo
	if spawnEvent, ok := spawn.(core.SpawnEvent); ok {
		spawnInfo := &SpawnInfo{
			ID:     spawnEvent.Word.ID,
			Text:   spawnEvent.Word.Text,
			Rarity: string(spawnEvent.Word.Rarity),
			Points: spawnEvent.Word.Points,
		}
		c.JSON(http.StatusOK, spawnInfo)
	} else {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "invalid_spawn",
			Message: "Erreur interne: spawn invalide",
		})
	}
}

// AttemptCapture tente de capturer un WordMon
func (h *Handlers) AttemptCapture(c *gin.Context) {
	var req CaptureAttemptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "playerId et attempt requis",
		})
		return
	}

	// Vérifier que le joueur existe
	player, err := h.playerStore.GetPlayer(req.PlayerID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "player_not_found",
			Message: "Joueur non trouvé",
		})
		return
	}

	// Vérifier qu'il y a un spawn actif
	spawn := h.spawnStore.GetCurrentSpawn()
	if spawn == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "no_spawn",
			Message: "Aucun WordMon actif",
		})
		return
	}

	// Convertir le spawn en SpawnEvent
	spawnEvent, ok := spawn.(core.SpawnEvent)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "invalid_spawn",
			Message: "Erreur interne: spawn invalide",
		})
		return
	}

	// Simuler la logique de capture
	// Pour l'instant, on considère que c'est une capture réussie si l'essai correspond
	if req.Attempt == spawnEvent.Word.Text {
		// Capture réussie
		player.XP += spawnEvent.Word.Points
		player.Level = 1 + player.XP/100
		player.Inventory[spawnEvent.Word.Text]++

		// Mettre à jour le joueur dans le store
		h.playerStore.UpdatePlayer(player)

		c.JSON(http.StatusOK, CaptureResultResponse{
			Status:   "captured",
			Word:     spawnEvent.Word.Text,
			Rarity:   string(spawnEvent.Word.Rarity),
			XP:       spawnEvent.Word.Points,
			NewLevel: player.Level,
		})
	} else {
		// Capture échouée
		c.JSON(http.StatusOK, CaptureResultResponse{
			Status: "fled",
			Word:   spawnEvent.Word.Text,
			Reason: "wrong attempt",
		})
	}
}

// GetLeaderboard retourne le classement des joueurs
func (h *Handlers) GetLeaderboard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	// Récupérer tous les joueurs depuis le store
	players := h.playerStore.GetAllPlayers()

	// Trier par XP décroissant
	sort.Slice(players, func(i, j int) bool {
		return players[i].XP > players[j].XP
	})

	// Limiter le nombre de résultats
	if len(players) > limit {
		players = players[:limit]
	}

	// Convertir en LeaderboardEntry
	entries := make([]LeaderboardEntry, len(players))
	for i, player := range players {
		entries[i] = LeaderboardEntry{
			ID:    player.ID,
			Name:  player.Name,
			XP:    player.XP,
			Level: player.Level,
		}
	}

	c.JSON(http.StatusOK, entries)
}

// UpdateCurrentSpawn met à jour le spawn actuel (appelé par le spawner)
func (h *Handlers) UpdateCurrentSpawn(spawn core.SpawnEvent) {
	// Mettre à jour le spawn dans le store
	h.spawnStore.AddSpawn(spawn)
}
