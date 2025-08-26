package core

// Rarity représente la rareté d'un WordMon.
type Rarity string

const (
	Common    Rarity = "Common"
	Rare      Rarity = "Rare"
	Legendary Rarity = "Legendary"
)

// Word représente une créature-mot capturable.
type Word struct {
	ID     string
	Text   string
	Rarity Rarity
	Points int
}

// Player représente le dresseur.
type Player struct {
	ID        string
	Name      string
	XP        int
	Level     int
	Inventory map[string]int // mot -> quantité
}

type SpawnEvent struct {
	Round int
	Word  Word
}
