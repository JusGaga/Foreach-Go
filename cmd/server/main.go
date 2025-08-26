package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jusgaga/wordmon-go/internal/core"
)

var Version = "0.0.1" // peut être surchargée via -ldflags "-X 'main.Version=0.1.0'"

func main() {
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

	fmt.Printf("WordMon Go v%v — prêt !\n", Version)

	if playerName == "" {
		playerName = "Guest"
	}

	rand.Seed(time.Now().UnixNano())

	p := core.NewPlayer(playerName)

	// ===== Exo 02: une rencontre simple =====
	core.PrintPlayer(p)

	e := core.NewEncounter()
	if err := e.Start(p); err != nil {
		fmt.Println("erreur:", err)
		return
	}
	fmt.Println("Un WordMon apparaît:", e.WordMon().Presentation())

	if err := e.BeginBattle(); err != nil {
		fmt.Println("erreur:", err)
		return
	}
	fmt.Println("Défi:", e.CurrentChallenge().Instructions())

	// Simulation: on génère une tentative plausible pour l'anagramme.
	attempt := core.AutoAttemptFor(e)
	fmt.Printf("Tentative: %q\n", attempt)
	if won, err := e.SubmitAttempt(attempt); err != nil {
		fmt.Println("erreur:", err)
		return
	} else if won {
		fmt.Println("→ VICTOIRE !")
	} else {
		fmt.Println("→ ÉCHEC.")
	}

	if err := e.Resolve(); err != nil {
		fmt.Println("erreur:", err)
	}

	core.PrintPlayer(p)

}

func StartSpawner(ctx context.Context, ch chan<- SpawnEvent, interval time.Duration) {
	for ctx.Err() != nil {
		ch <- core.SpawnWord()
		interval := 10 * time.Second
		time.Sleep(interval)
	}
}
