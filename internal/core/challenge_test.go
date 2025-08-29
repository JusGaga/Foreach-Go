package core

import (
	"strings"
	"testing"
)

func TestAnagramChallenge_Check(t *testing.T) {
	tests := []struct {
		name        string
		secret      string
		attempt     string
		expectValid bool
		expectError bool
	}{
		{"Anagramme valide", "hello", "olleh", true, false},
		{"Anagramme valide avec majuscules", "Hello", "OLLEH", true, false},
		{"Mot identique", "test", "test", false, true},
		{"Mot identique avec espaces", "test", " test ", false, true},
		{"Entrée vide", "word", "", false, true},
		{"Entrée avec espaces", "word", "  ", false, true},
		{"Lettres différentes", "hello", "world", false, false},
		{"Longueur différente", "hello", "hell", false, false},
		{"Anagramme complexe", "listen", "silent", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			challenge := &AnagramChallenge{
				secret: Word{Text: tt.secret},
			}

			valid, err := challenge.Check(tt.attempt)

			if tt.expectError {
				if err == nil {
					t.Error("Check devrait retourner une erreur")
				}
			} else {
				if err != nil {
					t.Errorf("Check ne devrait pas retourner d'erreur: %v", err)
				}
				if valid != tt.expectValid {
					t.Errorf("Validité = %v, attendu %v", valid, tt.expectValid)
				}
			}
		})
	}
}

func TestAnagramChallenge_Instructions(t *testing.T) {
	challenge := &AnagramChallenge{
		secret: Word{Text: "test"},
	}

	instructions := challenge.Instructions()
	expected := "Donne un anagramme valide de \"test\""

	if instructions != expected {
		t.Errorf("Instructions = %q, attendu %q", instructions, expected)
	}
}

func TestAnagramChallenge_ResetFor(t *testing.T) {
	challenge := &AnagramChallenge{}
	word := Word{Text: "new", Rarity: Rare}

	challenge.ResetFor(Rare, word)

	if challenge.secret.Text != "new" {
		t.Errorf("Secret = %q, attendu %q", challenge.secret.Text, "new")
	}
}

func TestSameMultiset(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{"Chaînes identiques", "hello", "hello", true},
		{"Anagrammes", "hello", "olleh", true},
		{"Longueurs différentes", "hello", "hell", false},
		{"Lettres différentes", "hello", "world", false},
		{"Chaînes vides", "", "", true},
		{"Un caractère", "a", "a", true},
		{"Caractères répétés", "aab", "aba", true},
		{"Caractères répétés différents", "aab", "abb", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sameMultiset(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("sameMultiset(%q, %q) = %v, attendu %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestAutoAttemptFor(t *testing.T) {
	encounter := Encounter{
		Word: Word{Text: "test"},
	}

	attempt := AutoAttemptFor(encounter)

	// Vérifier que l'attempt n'est pas vide
	if attempt == "" {
		t.Error("AutoAttemptFor devrait retourner une tentative non vide")
	}

	// Vérifier que l'attempt est différent du mot original
	if strings.EqualFold(attempt, "test") {
		t.Error("AutoAttemptFor devrait retourner un mot différent")
	}

	// Vérifier que l'attempt a la même longueur (sauf si suffixé avec 'x')
	if len(attempt) != len("test") && !strings.HasSuffix(attempt, "x") {
		t.Errorf("Attempt %q devrait avoir la même longueur que %q", attempt, "test")
	}
}
