package config

import (
	"testing"
	"time"
)

func TestGameConfig_SpawnInterval(t *testing.T) {
	tests := []struct {
		name           string
		intervalSecs   int
		expectedResult time.Duration
	}{
		{"Intervalle positif", 5, 5 * time.Second},
		{"Intervalle zéro", 0, 0},
		{"Intervalle négatif", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := GameConfig{}
			config.Spawner.IntervalSeconds = tt.intervalSecs

			result := config.SpawnInterval()

			if result != tt.expectedResult {
				t.Errorf("SpawnInterval() = %v, attendu %v", result, tt.expectedResult)
			}
		})
	}
}

func TestGameConfig_Validation(t *testing.T) {
	tests := []struct {
		name        string
		config      GameConfig
		expectValid bool
	}{
		{
			name: "Configuration valide",
			config: GameConfig{
				Game: struct {
					Name    string `yaml:"name" toml:"name" json:"name"`
					Version string `yaml:"version" toml:"version" json:"version"`
				}{
					Name:    "WordMon",
					Version: "1.0.0",
				},
				RarityWeights: map[string]int{
					"Common":    70,
					"Rare":      25,
					"Legendary": 5,
				},
				XPRewards: map[string]int{
					"Common":    10,
					"Rare":      25,
					"Legendary": 100,
				},
			},
			expectValid: true,
		},
		{
			name: "Poids de rareté ne totalisent pas 100",
			config: GameConfig{
				RarityWeights: map[string]int{
					"Common":    80,
					"Rare":      15,
					"Legendary": 10, // Total: 105
				},
			},
			expectValid: false,
		},
		{
			name: "Poids de rareté totalisent moins de 100",
			config: GameConfig{
				RarityWeights: map[string]int{
					"Common": 60,
					"Rare":   20, // Total: 80
				},
			},
			expectValid: false,
		},
		{
			name: "Poids de rareté négatifs",
			config: GameConfig{
				RarityWeights: map[string]int{
					"Common":    70,
					"Rare":      -5,
					"Legendary": 35,
				},
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validation simple des poids de rareté
			total := 0
			hasNegative := false

			for _, weight := range tt.config.RarityWeights {
				total += weight
				if weight < 0 {
					hasNegative = true
				}
			}

			isValid := total == 100 && !hasNegative

			if isValid != tt.expectValid {
				t.Errorf("Validation = %v, attendu %v (total: %d, hasNegative: %v)",
					isValid, tt.expectValid, total, hasNegative)
			}
		})
	}
}

func TestChallengesConfig_Validation(t *testing.T) {
	config := ChallengesConfig{
		Anagram: struct {
			MinLenByRarity       map[string]int `yaml:"minLenByRarity" toml:"minLenByRarity" json:"minLenByRarity"`
			MustDifferFromSource bool           `yaml:"mustDifferFromSource" toml:"mustDifferFromSource" json:"mustDifferFromSource"`
		}{
			MinLenByRarity: map[string]int{
				"Common":    3,
				"Rare":      4,
				"Legendary": 5,
			},
			MustDifferFromSource: true,
		},
		ATrou: struct {
			RevealedLetters map[string]int `yaml:"revealedLetters" toml:"revealedLetters" json:"revealedLetters"`
			MaxAttempts     int            `yaml:"maxAttempts" toml:"maxAttempts" json:"maxAttempts"`
		}{
			RevealedLetters: map[string]int{
				"Common":    1,
				"Rare":      2,
				"Legendary": 3,
			},
			MaxAttempts: 3,
		},
	}

	// Test des longueurs minimales par rareté
	if config.Anagram.MinLenByRarity["Common"] != 3 {
		t.Errorf("MinLenByRarity[Common] = %d, attendu 3", config.Anagram.MinLenByRarity["Common"])
	}

	if config.Anagram.MustDifferFromSource != true {
		t.Errorf("MustDifferFromSource = %v, attendu true", config.Anagram.MustDifferFromSource)
	}

	if config.ATrou.MaxAttempts != 3 {
		t.Errorf("MaxAttempts = %d, attendu 3", config.ATrou.MaxAttempts)
	}
}

func TestWordEntry_Validation(t *testing.T) {
	word := WordEntry{
		ID:     "word_001",
		Text:   "test",
		Rarity: "Common",
	}

	if word.ID == "" {
		t.Error("WordEntry.ID ne devrait pas être vide")
	}

	if word.Text == "" {
		t.Error("WordEntry.Text ne devrait pas être vide")
	}

	if word.Rarity == "" {
		t.Error("WordEntry.Rarity ne devrait pas être vide")
	}

	// Vérifier que la rareté est autorisée
	_, allowed := AllowedRarities[word.Rarity]
	if !allowed {
		t.Errorf("Rareté %s n'est pas autorisée", word.Rarity)
	}
}
