package core

import (
	"testing"
)

func TestNewEncounter(t *testing.T) {
	enc := NewEncounter()

	if enc.Phase != StateIdle {
		t.Errorf("Nouvelle rencontre devrait être dans l'état IDLE, got %s", enc.Phase)
	}
}

func TestEncounter_State(t *testing.T) {
	enc := Encounter{Phase: StateInBattle}

	state := enc.State()
	if state != StateInBattle {
		t.Errorf("State() devrait retourner %s, got %s", StateInBattle, state)
	}
}

func TestEncounter_WordMon(t *testing.T) {
	word := Word{Text: "test", Points: 10}
	enc := Encounter{Word: word}

	wordMon := enc.WordMon()
	if wordMon.Text != "test" {
		t.Errorf("WordMon() devrait retourner %s, got %s", "test", wordMon.Text)
	}
}

func TestEncounter_CurrentChallenge(t *testing.T) {
	challenge := &AnagramChallenge{}
	enc := Encounter{Challenge: challenge}

	currentChallenge := enc.CurrentChallenge()
	if currentChallenge != challenge {
		t.Error("CurrentChallenge() devrait retourner le challenge actuel")
	}
}

func TestEncounter_BeginBattle_ValidTransition(t *testing.T) {
	enc := Encounter{Phase: StateEncounter}
	word := Word{Text: "test", Rarity: Common}
	enc.Word = word

	err := enc.BeginBattle()

	if err != nil {
		t.Errorf("BeginBattle ne devrait pas retourner d'erreur: %v", err)
	}
	if enc.Phase != StateInBattle {
		t.Errorf("Phase devrait être IN_BATTLE, got %s", enc.Phase)
	}
	if enc.Challenge == nil {
		t.Error("Challenge devrait être initialisé")
	}
}

func TestEncounter_BeginBattle_InvalidTransition(t *testing.T) {
	tests := []struct {
		name         string
		initialState State
	}{
		{"Depuis IDLE", StateIdle},
		{"Depuis IN_BATTLE", StateInBattle},
		{"Depuis WON", StateWon},
		{"Depuis LOST", StateLost},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := Encounter{Phase: tt.initialState}

			err := enc.BeginBattle()

			if err == nil {
				t.Error("BeginBattle devrait retourner une erreur pour une transition invalide")
			}
			if enc.Phase != tt.initialState {
				t.Errorf("Phase ne devrait pas changer, attendu %s, got %s", tt.initialState, enc.Phase)
			}
		})
	}
}

func TestEncounter_SubmitAttempt_ValidTransition(t *testing.T) {
	enc := Encounter{Phase: StateInBattle}
	challenge := &AnagramChallenge{}
	challenge.ResetFor(Common, Word{Text: "test"})
	enc.Challenge = challenge

	// Test avec une tentative valide
	valid, err := enc.SubmitAttempt("tset")

	if err != nil {
		t.Errorf("SubmitAttempt ne devrait pas retourner d'erreur: %v", err)
	}
	if !valid {
		t.Error("SubmitAttempt devrait retourner true pour une anagramme valide")
	}
	if enc.Phase != StateWon {
		t.Errorf("Phase devrait être WON, got %s", enc.Phase)
	}
}

func TestEncounter_SubmitAttempt_InvalidTransition(t *testing.T) {
	tests := []struct {
		name         string
		initialState State
	}{
		{"Depuis IDLE", StateIdle},
		{"Depuis ENCOUNTERED", StateEncounter},
		{"Depuis WON", StateWon},
		{"Depuis LOST", StateLost},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := Encounter{Phase: tt.initialState}

			_, err := enc.SubmitAttempt("test")

			if err == nil {
				t.Error("SubmitAttempt devrait retourner une erreur pour une transition invalide")
			}
		})
	}
}

func TestEncounter_SubmitAttempt_InvalidAttempt(t *testing.T) {
	enc := Encounter{Phase: StateInBattle}
	challenge := &AnagramChallenge{}
	challenge.ResetFor(Common, Word{Text: "test"})
	enc.Challenge = challenge

	// Test avec une tentative invalide
	_, err := enc.SubmitAttempt("test") // Mot identique

	if err == nil {
		t.Error("SubmitAttempt devrait retourner une erreur pour une tentative invalide")
	}
	if enc.Phase != StateInBattle {
		t.Errorf("Phase ne devrait pas changer, attendu IN_BATTLE, got %s", enc.Phase)
	}
}

func TestEncounter_Resolve_Won(t *testing.T) {
	player := &Player{Inventory: make(map[string]int)}
	word := Word{Text: "test", Points: 10}
	enc := Encounter{
		Phase:  StateWon,
		Player: player,
		Word:   word,
	}

	err := enc.Resolve()

	if err != nil {
		t.Errorf("Resolve ne devrait pas retourner d'erreur: %v", err)
	}
	if enc.Phase != StateEncounter {
		t.Errorf("Phase devrait être ENCOUNTERED, got %s", enc.Phase)
	}
	if player.Inventory["test"] != 1 {
		t.Error("Mot devrait être ajouté à l'inventaire")
	}
}

func TestEncounter_Resolve_Lost(t *testing.T) {
	enc := Encounter{Phase: StateLost}

	err := enc.Resolve()

	if err != nil {
		t.Errorf("Resolve ne devrait pas retourner d'erreur: %v", err)
	}
	if enc.Phase != StateEncounter {
		t.Errorf("Phase devrait être ENCOUNTERED, got %s", enc.Phase)
	}
}

func TestEncounter_Resolve_InvalidState(t *testing.T) {
	tests := []struct {
		name  string
		state State
	}{
		{"IDLE", StateIdle},
		{"ENCOUNTERED", StateEncounter},
		{"IN_BATTLE", StateInBattle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := Encounter{Phase: tt.state}

			err := enc.Resolve()

			if err == nil {
				t.Error("Resolve devrait retourner une erreur pour un état invalide")
			}
		})
	}
}
