package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jusgaga/wordmon-go/internal/config"
	"github.com/jusgaga/wordmon-go/internal/core"
)

var Version = "0.0.1" // peut être surchargée via -ldflags "-X 'main.Version=0.1.0'"

func main() {
	// Charger toutes les configurations
	gameData, err := config.LoadAll()
	if err != nil {
		log.Fatal("[main] Échec du chargement des configurations:", err)
	}

	// Afficher les informations de configuration
	fmt.Printf("-------------------------------------------------\n")
	fmt.Printf("WordMon Go v%s — Configuration chargée !\n", gameData.Game.Game.Version)
	fmt.Printf("-------------------------------------------------\n")

	var (
		showVersion bool
		playerName  string
	)
	flag.BoolVar(&showVersion, "version", false, "affiche la version")
	flag.BoolVar(&showVersion, "v", false, "affiche la version (abrégé)")

	flag.StringVar(&playerName, "player", "", "nom du joueur")
	flag.StringVar(&playerName, "p", "", "nom du joueur (abrégé)")

	flag.Parse()

	if showVersion {
		fmt.Println("v" + Version)
		os.Exit(0)
	}

	if playerName == "" {
		playerName = "Guest"
	}

	rand.Seed(time.Now().UnixNano())

	p := core.NewPlayer(playerName)

	// ===== Exo 02: une rencontre simple =====
	core.PrintPlayer(p)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Gestion de l'arrêt propre avec sauvegarde
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine pour gérer l'arrêt propre
	go func() {
		<-sigChan
		log.Println("[main] Arrêt demandé, sauvegarde en cours...")

		// Sauvegarder le snapshot
		if err := config.SaveSnapshot([]core.Player{p}, ""); err != nil {
			log.Printf("[main] Erreur lors de la sauvegarde: %v", err)
		} else {
			log.Println("[main] Snapshot sauvegardé avec succès")
		}

		cancel()
	}()

	enc := core.NewEncounter()

	spawnsCh := make(chan core.SpawnEvent, 1)
	attemptsCh := make(chan core.Attempts, 1)

	// Utiliser l'intervalle de spawn depuis la configuration
	spawnInterval := gameData.Game.SpawnInterval()
	if spawnInterval == 0 {
		spawnInterval = 30 * time.Second // fallback
	}

	if err := enc.Start(&p, spawnsCh, spawnInterval); err != nil {
		fmt.Println("erreur:", err)
		return
	}

	log.Printf("[main] Démarrage du jeu avec intervalle de spawn: %v", spawnInterval)
	core.StartListen(ctx, &p, spawnsCh, attemptsCh, spawnInterval, enc)

	// Sauvegarde finale si pas d'arrêt par signal
	log.Println("[main] Sauvegarde finale...")
	if err := config.SaveSnapshot([]core.Player{p}, ""); err != nil {
		log.Printf("[main] Erreur lors de la sauvegarde finale: %v", err)
	} else {
		log.Println("[main] Sauvegarde finale réussie")
	}
}
