// Package core contient la logique métier du jeu WordMon.
// Il gère les joueurs, leurs interactions avec les mots et les mécaniques de progression.
package core

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// LevelFromXP calcule un niveau simple: 1 + XP/100.
// Le niveau augmente de 1 tous les 100 points d'expérience.
func LevelFromXP(xp int) int {
	if xp < 0 {
		return 1
	}
	return 1 + xp/100
}

// AwardXP ajoute des points et recalcule le niveau.
// Met à jour l'expérience du joueur et recalcule automatiquement son niveau.
func AwardXP(p *Player, points int) error {
	if points < 0 {
		return &NegativePointsError{Points: points}
	}
	p.XP += points
	p.Level = LevelFromXP(p.XP)
	return nil
}

// Capture ajoute le mot à l'inventaire et retourne les points gagnés.
// Vérifie que le mot n'est pas vide avant de l'ajouter à l'inventaire.
func Capture(p *Player, w Word) (int, error) {
	if w.Text == "" {
		return 0, &CaptureError{Word: w.Text, Reason: "mot vide"}
	}
	p.Inventory[w.Text]++
	return w.Points, nil
}

// NewPlayer crée un nouveau joueur avec les valeurs par défaut.
// Initialise un joueur avec 0 XP, niveau 1 et un inventaire vide.
func NewPlayer(name string) Player {
	return Player{
		ID:        "ID",
		Name:      name,
		XP:        0,
		Level:     1,
		Inventory: make(map[string]int),
	}
}

// PrintPlayer affiche les informations d'un joueur.
// Affiche le nom, l'XP, le niveau et la taille de l'inventaire.
func PrintPlayer(p Player) {
	fmt.Printf("Joueur: %s | XP: %d | Level: %d | Inventaire: %d mot(s)\n",
		p.Name, p.XP, p.Level, len(p.Inventory))
}

// StartListen orchestre les rounds en utilisant la machine d'états Encounter.
// Gère le cycle de jeu principal avec les apparitions de mots et les tentatives de capture.
func StartListen(ctx context.Context, p *Player, spawns chan SpawnEvent, attempts chan Attempts, interval time.Duration, enc Encounter) {
	if interval <= 0 {
		interval = time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Timeout de combat par WordMon
	battleTimeout := 5 * time.Second

	for {
		select {
		case <-ctx.Done():
			return

		case ev, ok := <-spawns:
			if !ok {
				return
			}

			// Démarrage du combat (IN_BATTLE)
			if err := enc.BeginBattle(); err != nil {
				fmt.Printf("[Round %d] Erreur BeginBattle pour %q: %v\n", ev.Round, enc.Word.Text, err)
				// Impossible de combattre, on passe à FLED et on revient à IDLE
				fmt.Printf("[Annonce] Rencontre interrompue: %q → état=%s\n", enc.Word.Text, enc.State())
				continue
			}
			fmt.Printf("[Round %d] Combat lancé contre %q → état=%s\n", ev.Round, enc.Word.Text, enc.State())

			// Le joueur décide éventuellement de tenter la capture
			if rand.Intn(100) < 80 {
				fmt.Printf("[%s] tente: %q (round %d)\n", p.Name, ev.Word.Text, ev.Round)

				test := AutoAttemptFor(enc)
				if won, err := enc.SubmitAttempt(test); err != nil {
					fmt.Println("erreur:", err)
					return
				} else {
					// On envoie la tentative avec le résultat (Won) dans le canal
					attempts <- Attempts{
						Player: *p,
						Word:   ev.Word,
						Round:  ev.Round,
						Won:    won,
					}
				}
			} else {
				enc.Phase = StateEncounter
			}

			// Fenêtre de combat: première tentative reçue ⇒ victoire/défaite, sinon timeout
			select {
			case att := <-attempts:
				// Résolution (CAPTURED/FLED) puis retour à IDLE
				if err := enc.Resolve(); err != nil {
					fmt.Printf("[Round %d] Erreur de résolution pour %q: %v\n", att.Round, enc.Word.Text, err)
				}
				issue := map[bool]string{true: "VICTOIRE", false: "DEFAITE"}[att.Won]
				if att.Won {
					if points, err := Capture(p, att.Word); err != nil {
						fmt.Printf("[Round %d] Erreur: %q \n", att.Round, err)
					} else if points > 0 {
						fmt.Printf("[Round %d] Capture de %q\n", att.Round, att.Word.Text)
					}
				}
				fmt.Printf("[Annonce] Issue du combat (round %d) pour %q: joueur %s → %s\n",
					att.Round, att.Word.Text, att.Player.Name, issue)

			case <-time.After(battleTimeout):
				// Aucune tentative dans le délai → le WordMon s’enfuit
				fmt.Printf("[Round %d] %q s’est enfui (timeout %s) → état=%s\n",
					ev.Round, enc.Word.Text, battleTimeout, enc.State())
				// Annonce globale
				fmt.Printf("[Annonce] Aucun vainqueur (round %d) pour %q, expiration du délai.\n",
					ev.Round, enc.Word.Text)
			}
		}
	}
}
