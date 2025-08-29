package config

import "os"

const (
	envChallengesPath = "WORDMON_CHALLENGES_PATH"
)

func LoadChallenges(path string) (*ChallengesConfig, error) {
	if env := os.Getenv(envChallengesPath); env != "" {
		path = env
	}
	if path == "" {
		return nil, &ValidationError{Section: "challenges", Problems: []string{"aucun chemin fourni (WORDMON_CHALLENGES_PATH ou argument requis)"}}
	}
	// YAML ou TOML
	if err := mustBeYAMLorTOML(path); err != nil {
		return nil, err
	}

	var cfg ChallengesConfig
	if err := decodeFile(path, &cfg); err != nil {
		return nil, err
	}
	if err := validateChallenges(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validateChallenges(c *ChallengesConfig) error {
	e := newValidationError("challenges")

	// Anagram: minLenByRarity > 0, clés valides
	if len(c.Anagram.MinLenByRarity) == 0 {
		e.addf("anagram.minLenByRarity manquant ou vide")
	} else {
		for k, v := range c.Anagram.MinLenByRarity {
			if !isAllowedRarity(k) {
				e.addf("anagram.minLenByRarity: rareté inconnue '%s'", k)
			}
			if v <= 0 {
				e.addf("anagram.minLenByRarity[%s] doit être > 0 (actuel %d)", k, v)
			}
		}
	}

	// aTrou: revealedLetters >= 0 par rareté, maxAttempts >= 1
	if len(c.ATrou.RevealedLetters) == 0 {
		e.addf("aTrou.revealedLetters manquant ou vide")
	} else {
		for k, v := range c.ATrou.RevealedLetters {
			if !isAllowedRarity(k) {
				e.addf("aTrou.revealedLetters: rareté inconnue '%s'", k)
			}
			if v < 0 {
				e.addf("aTrou.revealedLetters[%s] doit être >= 0 (actuel %d)", k, v)
			}
		}
	}
	if c.ATrou.MaxAttempts <= 0 {
		e.addf("aTrou.maxAttempts doit être >= 1 (actuel %d)", c.ATrou.MaxAttempts)
	}

	if e.ok() {
		return nil
	}
	return e
}
