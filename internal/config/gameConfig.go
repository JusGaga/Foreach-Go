package config

import "time"

const (
	RarityCommon    = "Common"
	RarityRare      = "Rare"
	RarityLegendary = "Legendary"
)

var AllowedRarities = map[string]struct{}{
	RarityCommon:    {},
	RarityRare:      {},
	RarityLegendary: {},
}

type GameConfig struct {
	Game struct {
		Name    string `yaml:"name" toml:"name" json:"name"`
		Version string `yaml:"version" toml:"version" json:"version"`
	} `yaml:"game" toml:"game" json:"game"`

	RarityWeights map[string]int `yaml:"rarityWeights" toml:"rarityWeights" json:"rarityWeights"`
	XPRewards     map[string]int `yaml:"xpRewards" toml:"xpRewards" json:"xpRewards"`

	Spawner struct {
		IntervalSeconds    int `yaml:"intervalSeconds" toml:"intervalSeconds" json:"intervalSeconds"`
		AutoFleeAfterSecs  int `yaml:"autoFleeAfterSeconds" toml:"autoFleeAfterSeconds" json:"autoFleeAfterSeconds"`
	} `yaml:"spawner" toml:"spawner" json:"spawner"`

	Level struct {
		Base       int `yaml:"base" toml:"base" json:"base"`
		XPPerLevel int `yaml:"xpPerLevel" toml:"xpPerLevel" json:"xpPerLevel"`
	} `yaml:"level" toml:"level" json:"level"`
}

func (g GameConfig) SpawnInterval() time.Duration {
	secs := g.Spawner.IntervalSeconds
	if secs <= 0 {
		return 0
	}
	return time.Duration(secs) * time.Second
}

type ChallengesConfig struct {
	Anagram struct {
		MinLenByRarity     map[string]int `yaml:"minLenByRarity" toml:"minLenByRarity" json:"minLenByRarity"`
		MustDifferFromSource bool         `yaml:"mustDifferFromSource" toml:"mustDifferFromSource" json:"mustDifferFromSource"`
	} `yaml:"anagram" toml:"anagram" json:"anagram"`

	ATrou struct {
		RevealedLetters map[string]int `yaml:"revealedLetters" toml:"revealedLetters" json:"revealedLetters"`
		MaxAttempts     int            `yaml:"maxAttempts" toml:"maxAttempts" json:"maxAttempts"`
	} `yaml:"aTrou" toml:"aTrou" json:"aTrou"`
}

type WordEntry struct {
	ID     string `json:"id" yaml:"id" toml:"id"`
	Text   string `json:"text" yaml:"text" toml:"text"`
	Rarity string `json:"rarity" yaml:"rarity" toml:"rarity"`
}
