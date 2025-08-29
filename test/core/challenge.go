// Package core contient les défis et mini-jeux du système WordMon.
// Il définit les interfaces et implémentations des challenges que les joueurs doivent résoudre.
package core

import (
	"math/rand"
	"strings"
)

// Challenge définit un mini-jeu vérifiable.
// Chaque challenge doit fournir des instructions et une méthode de vérification.
type Challenge interface {
	Instructions() string
	Check(attempt string) (bool, error)
	ResetFor(r Rarity, word Word)
}

// AnagramChallenge: réussir une anagramme (réarrangement des lettres, différent de l'original).
// Le joueur doit créer un mot en réarrangeant les lettres du mot secret.
type AnagramChallenge struct {
	secret Word
}

// Instructions retourne les instructions du défi anagramme.
// Explique au joueur ce qu'il doit faire.
func (a *AnagramChallenge) Instructions() string {
	return "Donne un anagramme valide de \"" + a.secret.Text + "\""
}

// ResetFor initialise le challenge avec un nouveau mot et une rareté.
// Prépare le challenge pour une nouvelle rencontre.
func (a *AnagramChallenge) ResetFor(r Rarity, w Word) { a.secret = w }

// Check vérifie si la tentative est valide.
// Valide que l'anagramme est correct et différent du mot original.
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

// sameMultiset vérifie si deux chaînes contiennent les mêmes caractères.
// Utilise un algorithme de comptage pour comparer les multiset de caractères.
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
// Crée automatiquement une anagramme valide pour les tests et démonstrations.
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
