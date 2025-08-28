package core

import (
	"context"
	"fmt"
	"time"
)

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

// Encounter orchestre une rencontre → combat → résolution.

func NewEncounter() Encounter { return Encounter{Phase: StateIdle} }

func (e Encounter) State() State                { return e.Phase }
func (e Encounter) WordMon() Word               { return e.Word }
func (e Encounter) CurrentChallenge() Challenge { return e.Challenge }

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
