package core

import (
	"context"
	"testing"
	"time"
)

func TestNewPlayer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Joueur avec nom simple"},
		{"Joueur avec nom vide"},
		{"Joueur avec nom long"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(tt.name)

			if player.Name != tt.name {
				t.Errorf("Nom = %q, attendu %q", player.Name, tt.name)
			}

			if player.XP != 0 {
				t.Errorf("XP initial = %d, attendu 0", player.XP)
			}

			if player.Level != 1 {
				t.Errorf("Niveau initial = %d, attendu 1", player.Level)
			}

			if player.Inventory == nil {
				t.Error("Inventaire non initialisé")
			}

			if len(player.Inventory) != 0 {
				t.Errorf("Inventaire non vide: %d éléments", len(player.Inventory))
			}
		})
	}
}

func TestPrintPlayer(t *testing.T) {
	player := Player{
		Name:  "TestPlayer",
		XP:    150,
		Level: 3,
		Inventory: map[string]int{
			"hello": 2,
			"world": 1,
		},
	}

	// Test que la fonction ne panique pas
	PrintPlayer(player)
}

func TestSpawnEvent(t *testing.T) {
	word := Word{Text: "test", Points: 10, Rarity: Common}
	event := SpawnEvent{
		Round: 5,
		Word:  word,
	}

	if event.Round != 5 {
		t.Errorf("Round = %d, attendu 5", event.Round)
	}

	if event.Word.Text != "test" {
		t.Errorf("Word.Text = %q, attendu %q", event.Word.Text, "test")
	}

	if event.Word.Points != 10 {
		t.Errorf("Word.Points = %d, attendu 10", event.Word.Points)
	}
}

func TestAttempts(t *testing.T) {
	player := Player{Name: "TestPlayer"}
	word := Word{Text: "test", Points: 10}

	tests := []struct {
		name   string
		round  int
		player Player
		word   Word
		won    bool
	}{
		{"Victoire", 1, player, word, true},
		{"Défaite", 2, player, word, false},
		{"Round zéro", 0, player, word, true},
		{"Round négatif", -1, player, word, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attempt := Attempts{
				Round:  tt.round,
				Player: tt.player,
				Word:   tt.word,
				Won:    tt.won,
			}

			if attempt.Round != tt.round {
				t.Errorf("Round = %d, attendu %d", attempt.Round, tt.round)
			}

			if attempt.Player.Name != tt.player.Name {
				t.Errorf("Player.Name = %q, attendu %q", attempt.Player.Name, tt.player.Name)
			}

			if attempt.Word.Text != tt.word.Text {
				t.Errorf("Word.Text = %q, attendu %q", attempt.Word.Text, tt.word.Text)
			}

			if attempt.Won != tt.won {
				t.Errorf("Won = %v, attendu %v", attempt.Won, tt.won)
			}
		})
	}
}

func TestRarityConstants(t *testing.T) {
	if Common != "Common" {
		t.Errorf("Common = %q, attendu %q", Common, "Common")
	}

	if Rare != "Rare" {
		t.Errorf("Rare = %q, attendu %q", Rare, "Rare")
	}

	if Legendary != "Legendary" {
		t.Errorf("Legendary = %q, attendu %q", Legendary, "Legendary")
	}
}

func TestWordStructure(t *testing.T) {
	word := Word{
		ID:     "word_001",
		Text:   "hello",
		Rarity: Rare,
		Points: 25,
	}

	if word.ID != "word_001" {
		t.Errorf("ID = %q, attendu %q", word.ID, "word_001")
	}

	if word.Text != "hello" {
		t.Errorf("Text = %q, attendu %q", word.Text, "hello")
	}

	if word.Rarity != Rare {
		t.Errorf("Rarity = %q, attendu %q", word.Rarity, Rare)
	}

	if word.Points != 25 {
		t.Errorf("Points = %d, attendu 25", word.Points)
	}
}

func TestSpawnWord(t *testing.T) {
	// Test que SpawnWord retourne toujours un mot valide
	for i := 0; i < 100; i++ {
		word := SpawnWord()

		if word.Text == "" {
			t.Error("SpawnWord a retourné un mot vide")
		}

		if word.Points <= 0 {
			t.Error("SpawnWord a retourné un mot avec des points négatifs ou nuls")
		}

		// Vérifier que la rareté est valide
		validRarity := false
		switch word.Rarity {
		case Common, Rare, Legendary:
			validRarity = true
		}

		if !validRarity {
			t.Errorf("SpawnWord a retourné une rareté invalide: %s", word.Rarity)
		}
	}
}

func TestWordPresentation(t *testing.T) {
	tests := []struct {
		word     Word
		expected string
	}{
		{
			Word{Text: "test", Rarity: Common, Points: 10},
			"Common — \"test\" (+10 XP)",
		},
		{
			Word{Text: "hello", Rarity: Rare, Points: 25},
			"Rare — \"hello\" (+25 XP)",
		},
		{
			Word{Text: "legend", Rarity: Legendary, Points: 100},
			"Legendary — \"legend\" (+100 XP)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.word.Text, func(t *testing.T) {
			result := tt.word.Presentation()
			if result != tt.expected {
				t.Errorf("Presentation() = %q, attendu %q", result, tt.expected)
			}
		})
	}
}

func TestStartSpawner(t *testing.T) {
	// Test que StartSpawner fonctionne correctement
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan SpawnEvent, 10)

	// Démarrer le spawner avec un intervalle court
	go StartSpawner(ctx, ch, 10*time.Millisecond)

	// Attendre un peu pour qu'un spawn se produise
	time.Sleep(50 * time.Millisecond)

	// Vérifier qu'au moins un événement a été généré
	select {
	case event := <-ch:
		if event.Round <= 0 {
			t.Error("Round devrait être positif")
		}
		if event.Word.Text == "" {
			t.Error("Word.Text ne devrait pas être vide")
		}
	default:
		t.Error("Aucun événement de spawn généré")
	}

	// Annuler le contexte pour arrêter le spawner
	cancel()

	// Le test se termine ici, le contexte annulé arrêtera le spawner
}
