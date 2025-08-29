// Package config contient la configuration du jeu WordMon.
// Il définit les paramètres de jeu, les poids de rareté et la configuration des défis.
package config

import "time"

const (
	RarityCommon    = "Common"
	RarityRare      = "Rare"
	RarityLegendary = "Legendary"
)

// AllowedRarities définit les raretés autorisées dans le jeu.
var AllowedRarities = map[string]struct{}{
	RarityCommon:    {},
	RarityRare:      {},
	RarityLegendary: {},
}

// GameConfig contient la configuration principale du jeu.
// Définit les paramètres de spawn, les récompenses XP et la progression des niveaux.
type GameConfig struct {
	Game struct {
		Name    string `yaml:"name" toml:"name" json:"name"`
		Version string `yaml:"version" toml:"version" json:"version"`
	} `yaml:"game" toml:"game" json:"game"`

	RarityWeights map[string]int `yaml:"rarityWeights" toml:"rarityWeights" json:"rarityWeights"`
	XPRewards     map[string]int `yaml:"xpRewards" toml:"xpRewards" json:"xpRewards"`

	Spawner struct {
		IntervalSeconds   int `yaml:"intervalSeconds" toml:"intervalSeconds" json:"intervalSeconds"`
		AutoFleeAfterSecs int `yaml:"autoFleeAfterSeconds" toml:"autoFleeAfterSeconds" json:"autoFleeAfterSeconds"`
	} `yaml:"spawner" toml:"spawner" json:"spawner"`

	Level struct {
		Base       int `yaml:"base" toml:"base" json:"base"`
		XPPerLevel int `yaml:"xpPerLevel" toml:"xpPerLevel" json:"xpPerLevel"`
	} `yaml:"level" toml:"level" json:"level"`
}

// SpawnInterval retourne l'intervalle de spawn en tant que Duration.
// Convertit les secondes de configuration en time.Duration.
func (g GameConfig) SpawnInterval() time.Duration {
	secs := g.Spawner.IntervalSeconds
	if secs <= 0 {
		return 0
	}
	return time.Duration(secs) * time.Second
}

// ChallengesConfig définit la configuration des différents types de défis.
// Contient les paramètres pour les anagrammes et les mots à trous.
type ChallengesConfig struct {
	Anagram struct {
		MinLenByRarity       map[string]int `yaml:"minLenByRarity" toml:"minLenByRarity" json:"minLenByRarity"`
		MustDifferFromSource bool           `yaml:"mustDifferFromSource" toml:"mustDifferFromSource" json:"mustDifferFromSource"`
	} `yaml:"anagram" toml:"anagram" json:"anagram"`

	ATrou struct {
		RevealedLetters map[string]int `yaml:"revealedLetters" toml:"revealedLetters" json:"revealedLetters"`
		MaxAttempts     int            `yaml:"maxAttempts" toml:"maxAttempts" json:"maxAttempts"`
	} `yaml:"aTrou" toml:"aTrou" json:"aTrou"`
}

// WordEntry représente une entrée de mot dans la base de données.
// Contient l'identifiant, le texte et la rareté d'un mot.
type WordEntry struct {
	ID     string `json:"id" yaml:"id" toml:"id"`
	Text   string `json:"text" yaml:"text" toml:"text"`
	Rarity string `json:"rarity" yaml:"rarity" toml:"rarity"`
}
