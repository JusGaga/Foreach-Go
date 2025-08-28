package config

import (
	"os"
	"strconv"
)

const (
	envConfigPath     = "WORDMON_CONFIG_PATH"
	envSpawnInterval  = "WORDMON_SPAWN_INTERVAL"
	DefaultXPPerLevel = 100
)

func LoadGameConfig(path string) (*GameConfig, error) {
	if env := os.Getenv(envConfigPath); env != "" {
		path = env
	}
	if path == "" {
		return nil, &ValidationError{Section: "game", Problems: []string{"aucun chemin fourni (WORDMON_CONFIG_PATH ou argument requis)"}}
	}
	if err := mustBeYAMLorTOML(path); err != nil {
		return nil, err
	}

	var cfg GameConfig
	if err := decodeFile(path, &cfg); err != nil {
		return nil, err
	}

	// Valeurs par défaut si manquantes
	if cfg.Level.XPPerLevel == 0 {
		cfg.Level.XPPerLevel = DefaultXPPerLevel
	}

	// Overrides d’environnement
	if v := os.Getenv(envSpawnInterval); v != "" {
		iv, err := strconv.Atoi(v)
		if err != nil || iv <= 0 {
			return nil, &EnvOverrideError{Var: envSpawnInterval, Value: v, Reason: "doit être un entier > 0 (secondes)"}
		}
		cfg.Spawner.IntervalSeconds = iv
	}

	// Validation
	if err := validateGameConfig(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validateGameConfig(c *GameConfig) error {
	e := newValidationError("game")

	// RarityWeights: clés autorisées, valeurs >0, somme=100
	sum := 0
	if len(c.RarityWeights) == 0 {
		e.addf("rarityWeights manquant ou vide")
	} else {
		for k, v := range c.RarityWeights {
			if !isAllowedRarity(k) {
				e.addf("rarityWeights: rareté inconnue '%s'", k)
			}
			if v <= 0 {
				e.addf("rarityWeights[%s] doit être > 0 (actuel %d)", k, v)
			}
			sum += v
		}
		if sum != 100 {
			e.addf("la somme de rarityWeights doit être 100 (actuelle %d)", sum)
		}
	}

	// XPRewards: clés autorisées, valeurs >0
	if len(c.XPRewards) == 0 {
		e.addf("xpRewards manquant ou vide")
	} else {
		for k, v := range c.XPRewards {
			if !isAllowedRarity(k) {
				e.addf("xpRewards: rareté inconnue '%s'", k)
			}
			if v <= 0 {
				e.addf("xpRewards[%s] doit être > 0 (actuel %d)", k, v)
			}
		}
	}

	// Spawner
	if c.Spawner.IntervalSeconds <= 0 {
		e.addf("spawner.intervalSeconds doit être > 0 (actuel %d)", c.Spawner.IntervalSeconds)
	}
	if c.Spawner.AutoFleeAfterSecs < 0 {
		e.addf("spawner.autoFleeAfterSeconds doit être >= 0 (actuel %d)", c.Spawner.AutoFleeAfterSecs)
	}

	// Level
	if c.Level.Base <= 0 {
		e.addf("level.base doit être >= 1 (actuel %d)", c.Level.Base)
	}
	if c.Level.XPPerLevel <= 0 {
		e.addf("level.xpPerLevel doit être > 0 (actuel %d)", c.Level.XPPerLevel)
	}

	if e.ok() {
		return nil
	}
	return e
}
