package core

import (
	"testing"
)

func TestLevelFromXP(t *testing.T) {
	tests := []struct {
		name     string
		xp       int
		expected int
	}{
		{"XP négatif", -50, 1},
		{"XP zéro", 0, 1},
		{"XP 50", 50, 1},
		{"XP 100", 100, 2},
		{"XP 150", 150, 2},
		{"XP 200", 200, 3},
		{"XP 250", 250, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LevelFromXP(tt.xp)
			if result != tt.expected {
				t.Errorf("LevelFromXP(%d) = %d, attendu %d", tt.xp, result, tt.expected)
			}
		})
	}
}

func TestAwardXP(t *testing.T) {
	tests := []struct {
		name          string
		initialXP     int
		points        int
		expectError   bool
		expectedXP    int
		expectedLevel int
	}{
		{"Points positifs", 0, 50, false, 50, 1},
		{"Points négatifs", 100, -25, true, 100, 2},
		{"Points zéro", 75, 0, false, 75, 1},
		{"Niveau supérieur", 150, 100, false, 250, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &Player{XP: tt.initialXP, Level: LevelFromXP(tt.initialXP)}

			err := AwardXP(player, tt.points)

			if tt.expectError {
				if err == nil {
					t.Error("AwardXP devrait retourner une erreur")
				}
				if player.XP != tt.expectedXP {
					t.Errorf("XP ne devrait pas changer: %d, attendu %d", player.XP, tt.expectedXP)
				}
			} else {
				if err != nil {
					t.Errorf("AwardXP ne devrait pas retourner d'erreur: %v", err)
				}
				if player.XP != tt.expectedXP {
					t.Errorf("XP = %d, attendu %d", player.XP, tt.expectedXP)
				}
				if player.Level != tt.expectedLevel {
					t.Errorf("Level = %d, attendu %d", player.Level, tt.expectedLevel)
				}
			}
		})
	}
}

func TestCapture(t *testing.T) {
	tests := []struct {
		name           string
		word           Word
		expectError    bool
		expectedPoints int
	}{
		{"Mot valide", Word{Text: "test", Points: 10}, false, 10},
		{"Mot vide", Word{Text: "", Points: 5}, true, 0},
		{"Mot avec points zéro", Word{Text: "hello", Points: 0}, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &Player{Inventory: make(map[string]int)}

			points, err := Capture(player, tt.word)

			if tt.expectError {
				if err == nil {
					t.Error("Capture devrait retourner une erreur")
				}
			} else {
				if err != nil {
					t.Errorf("Capture ne devrait pas retourner d'erreur: %v", err)
				}
				if points != tt.expectedPoints {
					t.Errorf("Points = %d, attendu %d", points, tt.expectedPoints)
				}
				if player.Inventory[tt.word.Text] != 1 {
					t.Errorf("Mot non ajouté à l'inventaire")
				}
			}
		})
	}
}
