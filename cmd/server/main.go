package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jusgaga/wordmon-go/internal/api"
	"github.com/jusgaga/wordmon-go/internal/config"
	"github.com/jusgaga/wordmon-go/internal/core"
	_ "github.com/lib/pq"
)

var Version = "0.3.0" // peut être surchargée via -ldflags "-X 'main.Version=0.1.0'"

func main() {
	// Charger toutes les configurations
	gameData, err := config.LoadAll()
	if err != nil {
		log.Fatal("[main] Échec du chargement des configurations:", err)
	}

	// Afficher les informations de configuration
	log.Printf("-------------------------------------------------")
	log.Printf("WordMon Go v%s — Configuration chargée !", gameData.Game.Game.Version)
	log.Printf("-------------------------------------------------")

	var (
		showVersion bool
		port        string
	)
	flag.BoolVar(&showVersion, "version", false, "affiche la version")
	flag.BoolVar(&showVersion, "v", false, "affiche la version (abrégé)")
	flag.StringVar(&port, "port", "8080", "port du serveur HTTP")

	flag.Parse()

	if showVersion {
		log.Println("v" + Version)
		os.Exit(0)
	}

	// Créer le store SQL avec retry
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		"postgres", "postgres", "localhost", "5432", "wordmon", "disable")

	var sqlStore *api.SQLStore
	maxRetries := 30
	retryDelay := 2 * time.Second

	log.Printf("[main] Tentative de connexion à la base de données...")

	for i := 0; i < maxRetries; i++ {
		var err error
		sqlStore, err = api.NewSQLStore(databaseURL)
		if err == nil {
			log.Printf("[main] Connexion à la base de données réussie!")
			break
		}

		log.Printf("[main] Tentative %d/%d échouée: %v", i+1, maxRetries, err)

		if i < maxRetries-1 {
			log.Printf("[main] Nouvelle tentative dans %v...", retryDelay)
			time.Sleep(retryDelay)
		} else {
			log.Fatal("[main] Échec de la connexion à la base de données après", maxRetries, "tentatives")
		}
	}

	defer sqlStore.Close()

	// Charger les mots dans la base de données
	words := gameData.Words
	coreWords := make([]core.Word, len(words))
	for i, word := range words {
		coreWords[i] = core.Word{
			ID:     word.ID,
			Text:   word.Text,
			Rarity: core.Rarity(word.Rarity),
			Points: 10, // Points par défaut
		}
	}

	if err := sqlStore.Seed(coreWords); err != nil {
		log.Fatal("[main] Échec du seeding de la base de données:", err)
	}

	// Créer le serveur API avec le store SQL
	server := api.NewServer(sqlStore, sqlStore)

	// Gestion de l'arrêt propre
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configurer le spawner
	spawnInterval := gameData.Game.SpawnInterval()
	if spawnInterval == 0 {
		spawnInterval = 30 * time.Second // fallback
	}

	// Canal pour les événements de spawn
	spawnCh := make(chan core.SpawnEvent, 10)

	// Goroutine pour traiter les spawns et les envoyer au store
	go func() {
		ticker := time.NewTicker(spawnInterval)
		defer ticker.Stop()

		round := 0
		for {
			select {
			case <-ticker.C:
				round++
				word := core.SpawnWord()

				spawnEvent := core.SpawnEvent{
					Round: round,
					Word:  word,
				}

				log.Printf("[spawn] Nouveau WordMon: %q (%s)", word.Text, word.Rarity)

				// Envoyer l'événement de spawn
				select {
				case spawnCh <- spawnEvent:
					// Événement envoyé avec succès
				default:
					// Canal plein, ignorer ce spawn
					log.Printf("[spawn] Canal plein, spawn ignoré: %q", word.Text)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Goroutine pour traiter les spawns et les envoyer au store
	go func() {
		for {
			select {
			case spawnEvent := <-spawnCh:
				sqlStore.AddSpawn(spawnEvent)
			case <-ctx.Done():
				return
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine pour gérer l'arrêt propre
	go func() {
		<-sigChan
		log.Println("[main] Arrêt demandé, arrêt du serveur...")

		// Arrêter le serveur HTTP
		if err := server.Stop(); err != nil {
			log.Printf("[main] Erreur lors de l'arrêt du serveur: %v", err)
		}

		cancel()
	}()

	// Démarrer le serveur HTTP
	addr := ":" + port
	log.Printf("[main] Démarrage du serveur API sur %s", addr)
	log.Printf("[main] Spawner démarré avec intervalle: %v", spawnInterval)

	if err := server.Start(addr); err != nil {
		log.Fatal("[main] Erreur du serveur:", err)
	}
}
