package core

import (
	"math/rand"
	"strings"
)

// Challenge définit un mini-jeu vérifiable.
type Challenge interface {
	Instructions() string
	Check(attempt string) (bool, error)
	ResetFor(r Rarity, word Word)
}

// AnagramChallenge: réussir une anagramme (réarrangement des lettres, différent de l'original).
type AnagramChallenge struct {
	secret Word
}

func (a *AnagramChallenge) Instructions() string {
	return "Donne un anagramme valide de \"" + a.secret.Text + "\""
}

func (a *AnagramChallenge) ResetFor(r Rarity, w Word) { a.secret = w }

func (a *AnagramChallenge) Check(attempt string) (bool, error) {
	attempt = strings.ToLower(strings.TrimSpace(attempt))
	if attempt == "" {
		return false, &InvalidAttemptError{Input: attempt, Reason: "entrée vide"}
	}
	sec := strings.ToLower(a.secret.Text)
	if attempt == sec {
		return false, &InvalidAttemptError{Input: attempt, Reason: "identique au mot"}
	}
	if !sameMultiset(sec, attempt) {
		return false, nil
	}
	return true, nil
}

func sameMultiset(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	m := map[rune]int{}
	for _, r := range a {
		m[r]++
	}
	for _, r := range b {
		m[r]--
		if m[r] < 0 {
			return false
		}
	}
	for _, v := range m {
		if v != 0 {
			return false
		}
	}
	return true
}

// AutoAttemptFor génère une tentative plausible pour démo (anagramme par shuffle).
func AutoAttemptFor(e Encounter) string {
	w := []rune(e.Word.Text)
	for i := range w {
		j := rand.Intn(i + 1)
		w[i], w[j] = w[j], w[i]
	}
	cand := string(w)
	if strings.EqualFold(cand, e.Word.Text) {
		cand = cand + "x"
	} // garantit différent
	return cand
}
