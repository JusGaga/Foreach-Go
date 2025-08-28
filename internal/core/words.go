package core

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var poolCommon = []Word{
	{ID: "c1", Text: "chat", Rarity: Common, Points: 5},
	{ID: "c2", Text: "chien", Rarity: Common, Points: 5},
	{ID: "c3", Text: "pomme", Rarity: Common, Points: 5},
	{ID: "c4", Text: "table", Rarity: Common, Points: 5},
	{ID: "c5", Text: "route", Rarity: Common, Points: 5},
}

var poolRare = []Word{
	{ID: "r1", Text: "carafe", Rarity: Rare, Points: 20},
	{ID: "r2", Text: "trace", Rarity: Rare, Points: 20},
	{ID: "r3", Text: "atelier", Rarity: Rare, Points: 20},
	{ID: "r4", Text: "portail", Rarity: Rare, Points: 20},
	{ID: "r5", Text: "lingerie", Rarity: Rare, Points: 20},
}

var poolLegendary = []Word{
	{ID: "l1", Text: "polyglotte", Rarity: Legendary, Points: 100},
	{ID: "l2", Text: "intergalaxie", Rarity: Legendary, Points: 100},
	{ID: "l3", Text: "mythologie", Rarity: Legendary, Points: 100},
	{ID: "l4", Text: "clairvoyance", Rarity: Legendary, Points: 100},
	{ID: "l5", Text: "transcendant", Rarity: Legendary, Points: 100},
}

// SpawnWord choisit une rareté selon des poids (~80/18/2), puis un mot dans la pool.
func SpawnWord() Word {
	x := rand.Intn(100)
	var pool []Word
	switch {
	case x < 80:
		pool = poolCommon
	case x < 98:
		pool = poolRare
	default:
		pool = poolLegendary
	}
	return pool[rand.Intn(len(pool))]
}

// Presentation retourne une fiche textuelle.
func (w Word) Presentation() string {
	return string(w.Rarity) + " — \"" + w.Text + "\" (+" + strconv.Itoa(w.Points) + " XP)"
}

func StartSpawner(ctx context.Context, ch chan<- SpawnEvent, interval time.Duration) {
	defer close(ch)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	round := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			round++
			w := SpawnWord()
			fmt.Printf("Un Pokémon a spawn ! Round %d : %s\n", round, w.Text)
			ch <- SpawnEvent{Round: round, Word: w}
		}
	}
}
