package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jusgaga/wordmon-go/internal/core"
)

// Server représente le serveur HTTP de l'API
type Server struct {
	router   *gin.Engine
	handlers *Handlers
	server   *http.Server
}

// NewServer crée une nouvelle instance du serveur
func NewServer(playerStore PlayerStore, spawnStore SpawnStore) *Server {
	// Configurer Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlers := NewHandlers(playerStore, spawnStore)

	server := &Server{
		router:   router,
		handlers: handlers,
	}

	server.setupRoutes()
	return server
}

// setupRoutes configure toutes les routes de l'API
func (s *Server) setupRoutes() {
	// Groupe API (optionnel pour la versioning)
	api := s.router.Group("/api")
	{
		// Status
		api.GET("/status", s.handlers.GetStatus)

		// Players
		api.POST("/players", s.handlers.CreatePlayer)
		api.GET("/players/:id", s.handlers.GetPlayer)

		// Spawn
		api.GET("/spawn/current", s.handlers.GetCurrentSpawn)

		// Encounter
		api.POST("/encounter/attempt", s.handlers.AttemptCapture)

		// Leaderboard
		api.GET("/leaderboard", s.handlers.GetLeaderboard)
	}

	// Routes racine pour la compatibilité
	s.router.GET("/status", s.handlers.GetStatus)
	s.router.POST("/players", s.handlers.CreatePlayer)
	s.router.GET("/players/:id", s.handlers.GetPlayer)
	s.router.GET("/spawn/current", s.handlers.GetCurrentSpawn)
	s.router.POST("/encounter/attempt", s.handlers.AttemptCapture)
	s.router.GET("/leaderboard", s.handlers.GetLeaderboard)
}

// Start démarre le serveur sur le port spécifié
func (s *Server) Start(addr string) error {
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	log.Printf("[API] Serveur démarré sur %s", addr)
	return s.server.ListenAndServe()
}

// Stop arrête proprement le serveur
func (s *Server) Stop() error {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.server.Shutdown(ctx)
	}
	return nil
}

// SetSpawner configure le spawner pour les handlers
func (s *Server) SetSpawner(spawner chan core.SpawnEvent) {
	s.handlers.SetSpawner(spawner)
}

// GetHandlers retourne les handlers pour l'intégration
func (s *Server) GetHandlers() *Handlers {
	return s.handlers
}
