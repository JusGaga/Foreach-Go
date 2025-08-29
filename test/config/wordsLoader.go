package config

import (
	"os"
)

const (
	envWordsPath = "WORDMON_WORDS_PATH"
)

func LoadWords(path string) ([]WordEntry, error) {
	if env := os.Getenv(envWordsPath); env != "" {
		path = env
	}
	if path == "" {
		return nil, &ValidationError{Section: "words", Problems: []string{"aucun chemin fourni (WORDMON_WORDS_PATH ou argument requis)"}}
	}
	if err := mustBeJSON(path); err != nil {
		return nil, err
	}

	var words []WordEntry
	if err := decodeFile(path, &words); err != nil {
		return nil, err
	}
	if err := validateWords(words); err != nil {
		return nil, err
	}
	return words, nil
}

func validateWords(words []WordEntry) error {
	e := newValidationError("words")

	if len(words) == 0 {
		e.addf("le dictionnaire est vide")
	}

	for i, w := range words {
		if strings := stringsTrim(w.ID); strings == "" {
			e.addf("mot #%d: id manquant", i+1)
		}
		if stringsTrim(w.Text) == "" {
			e.addf("mot #%d: text manquant", i+1)
		}
		if !isAllowedRarity(w.Rarity) {
			e.addf("mot #%d: rareté inconnue '%s'", i+1, w.Rarity)
		}
	}

	if e.ok() {
		return nil
	}
	return e
}

// petite fonction utilitaire locale pour éviter d'importer strings au niveau fichier
func stringsTrim(s string) string {
	// implémentation minimale sans coût: on évite les espaces simples
	// si besoin de plus, importer "strings" et utiliser strings.TrimSpace
	var start, end = 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
