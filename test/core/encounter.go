// Package core contient la machine à états des rencontres WordMon.
// Il gère les transitions d'état entre les différentes phases d'une rencontre.
package core

import (
	"context"
	"fmt"
	"time"
)

// State représente l'état actuel d'une rencontre.
type State string

const (
	StateIdle      State = "IDLE"
	StateEncounter State = "ENCOUNTERED"
	StateInBattle  State = "IN_BATTLE"
	StateWon       State = "WON"
	StateCaptured  State = "CAPTURED"
	StateLost      State = "LOST"
	StateFled      State = "FLED"
)

// NewEncounter crée une nouvelle rencontre dans l'état IDLE.
// Initialise une rencontre vide prête à démarrer.
func NewEncounter() Encounter { return Encounter{Phase: StateIdle} }

// State retourne l'état actuel de la rencontre.
func (e Encounter) State() State { return e.Phase }

// WordMon retourne le mot actuellement rencontré.
func (e Encounter) WordMon() Word { return e.Word }

// CurrentChallenge retourne le défi actuel de la rencontre.
func (e Encounter) CurrentChallenge() Challenge { return e.Challenge }

// Start démarre une nouvelle rencontre pour un joueur.
// Initialise le spawner et attend l'apparition d'un mot.
func (e *Encounter) Start(p *Player, spawnCh chan SpawnEvent, interval time.Duration) error {
	if e.Phase != StateIdle {
		return &InvalidStateError{From: string(e.Phase), Expected: string(StateIdle)}
	}
	e.Player = p
	ctx, cancel := context.WithCancel(context.Background())
	e.Cancel = cancel
	go StartSpawner(ctx, spawnCh, interval)
	ev := <-spawnCh
	e.Word = ev.Word
	if e.Word.Text == "" { // bug interne, cas exceptionnel → panic
		panic("WordMon invalide: mot vide")
	}
	e.Phase = StateEncounter
	return nil
}

// BeginBattle lance le combat en initialisant le défi.
// Passe de l'état ENCOUNTERED à IN_BATTLE.
func (e *Encounter) BeginBattle() error {
	if e.Phase != StateEncounter {
		return &InvalidStateError{From: string(e.Phase), Expected: string(StateEncounter)}
	}
	ch := &AnagramChallenge{}
	ch.ResetFor(e.Word.Rarity, e.Word)
	e.Challenge = ch
	e.Phase = StateInBattle
	return nil
}

// SubmitAttempt soumet une tentative de résolution du défi.
// Vérifie la réponse et met à jour l'état selon le résultat.
func (e *Encounter) SubmitAttempt(input string) (bool, error) {
	if e.Phase != StateInBattle {
		return false, &InvalidStateError{From: string(e.Phase), Expected: string(StateInBattle)}
	}
	ok, err := e.Challenge.Check(input)
	if err != nil {
		return false, fmt.Errorf("erreur de tentative: %w", err)
	}
	if ok {
		e.Phase = StateWon
	} else {
		e.Phase = StateLost
	}
	return ok, nil
}

// Resolve finalise la rencontre selon l'état actuel.
// Gère la capture en cas de victoire ou la fuite en cas de défaite.
func (e *Encounter) Resolve() error {
	switch e.Phase {
	case StateWon:
		_, err := Capture(e.Player, e.Word)
		if err != nil {
			return fmt.Errorf("capture: %w", err)
		}
		if err := AwardXP(e.Player, e.Word.Points); err != nil {
			return fmt.Errorf("xp: %w", err)
		}
		e.Phase = StateCaptured
		// retour à IDLE
		e.Phase = StateEncounter
		return nil
	case StateLost:
		e.Phase = StateFled
		e.Phase = StateEncounter
		return nil
	default:
		return &InvalidStateError{From: string(e.Phase), Expected: "WON ou LOST"}
	}
}
