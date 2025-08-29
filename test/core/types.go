// Package core contient les types et structures fondamentales du jeu WordMon.
// Il définit les entités principales comme les joueurs, les mots, les rencontres
// et les mécaniques de base du jeu.
package core

import "context"

// Rarity représente la rareté d'un WordMon.
type Rarity string

const (
	Common    Rarity = "Common"
	Rare      Rarity = "Rare"
	Legendary Rarity = "Legendary"
)

// Word représente une créature-mot capturable.
// Chaque Word a un identifiant unique, un texte, une rareté et des points d'expérience.
type Word struct {
	ID     string
	Text   string
	Rarity Rarity
	Points int
}

// Player représente le dresseur de mots.
// Un joueur a un identifiant, un nom, de l'expérience, un niveau et un inventaire.
type Player struct {
	ID        string
	Name      string
	XP        int
	Level     int
	Inventory map[string]int // mot -> quantité
}

// SpawnEvent représente l'apparition d'un mot dans le jeu.
// Il contient le numéro du round et le mot qui apparaît.
type SpawnEvent struct {
	Round int
	Word  Word
}

// Encounter représente une rencontre entre un joueur et un mot.
// Elle gère l'état de la rencontre, le joueur, le mot et le défi associé.
type Encounter struct {
	Phase     State
	Player    *Player
	Word      Word
	Challenge Challenge
	Cancel    context.CancelFunc
}

// Attempts représente les tentatives d'un joueur pour capturer un mot.
// Il enregistre le round, le joueur, le mot et le résultat de la tentative.
type Attempts struct {
	Round  int
	Player Player
	Word   Word
	Won    bool
}
