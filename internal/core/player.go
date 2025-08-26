package core

import "fmt"

// LevelFromXP calcule un niveau simple: 1 + XP/100.
func LevelFromXP(xp int) int {
	if xp < 0 {
		return 1
	}
	return 1 + xp/100
}

// AwardXP ajoute des points et recalcule le niveau.
func AwardXP(p *Player, points int) error {
	if points < 0 {
		return &NegativePointsError{Points: points}
	}
	p.XP += points
	p.Level = LevelFromXP(p.XP)
	return nil
}

// Capture ajoute le mot à l'inventaire et retourne les points gagnés.
func Capture(p *Player, w Word) (int, error) {
	if w.Text == "" {
		return 0, &CaptureError{Word: w.Text, Reason: "mot vide"}
	}
	p.Inventory[w.Text]++
	return w.Points, nil
}

func NewPlayer(name string) *Player {
	return &Player{
		ID:        "p1",
		Name:      name,
		XP:        0,
		Level:     1,
		Inventory: make(map[string]int),
	}
}

func PrintPlayer(p *Player) {
	fmt.Printf("Joueur: %s | XP: %d | Level: %d | Inventaire: %d mot(s)\n",
		p.Name, p.XP, p.Level, len(p.Inventory))
}
