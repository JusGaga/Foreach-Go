package core

import (
	"fmt"
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
type Encounter struct {
	state     State
	player    *Player
	word      Word
	challenge Challenge
}

func NewEncounter() *Encounter { return &Encounter{state: StateIdle} }

func (e *Encounter) State() State                { return e.state }
func (e *Encounter) WordMon() Word               { return e.word }
func (e *Encounter) CurrentChallenge() Challenge { return e.challenge }

func (e *Encounter) Start(p *Player) error {
	if e.state != StateIdle {
		return &InvalidStateError{From: string(e.state), Expected: string(StateIdle)}
	}
	e.player = p
	e.word = SpawnWord()
	if e.word.Text == "" { // bug interne, cas exceptionnel → panic
		panic("WordMon invalide: mot vide")
	}
	e.state = StateEncounter
	return nil
}

func (e *Encounter) BeginBattle() error {
	if e.state != StateEncounter {
		return &InvalidStateError{From: string(e.state), Expected: string(StateEncounter)}
	}
	ch := &AnagramChallenge{}
	ch.ResetFor(e.word.Rarity, e.word)
	e.challenge = ch
	e.state = StateInBattle
	return nil
}

func (e *Encounter) SubmitAttempt(input string) (bool, error) {
	if e.state != StateInBattle {
		return false, &InvalidStateError{From: string(e.state), Expected: string(StateInBattle)}
	}
	ok, err := e.challenge.Check(input)
	if err != nil {
		return false, fmt.Errorf("erreur de tentative: %w", err)
	}
	if ok {
		e.state = StateWon
	} else {
		e.state = StateLost
	}
	return ok, nil
}

func (e *Encounter) Resolve() error {
	switch e.state {
	case StateWon:
		_, err := Capture(e.player, e.word)
		if err != nil {
			return fmt.Errorf("capture: %w", err)
		}
		if err := AwardXP(e.player, e.word.Points); err != nil {
			return fmt.Errorf("xp: %w", err)
		}
		e.state = StateCaptured
		// retour à IDLE
		e.state = StateIdle
		return nil
	case StateLost:
		e.state = StateFled
		e.state = StateIdle
		return nil
	default:
		return &InvalidStateError{From: string(e.state), Expected: "WON ou LOST"}
	}
}
