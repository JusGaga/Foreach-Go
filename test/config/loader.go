package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// GameData contient toutes les configurations chargées
type GameData struct {
	Game           *GameConfig
	Challenges     *ChallengesConfig
	Words          []WordEntry
	ConfigPath     string
	WordsPath      string
	ChallengesPath string
}

// LoadAll charge toutes les configurations nécessaires au jeu
func LoadAll() (*GameData, error) {
	log.Println("[config] Chargement des configurations...")

	// Déterminer les chemins (avec fallback sur les valeurs par défaut)
	configPath := getenvOrDefault("WORDMON_CONFIG_PATH", "configs/game.yaml")
	wordsPath := getenvOrDefault("WORDMON_WORDS_PATH", "configs/words.json")
	challengesPath := getenvOrDefault("WORDMON_CHALLENGES_PATH", "configs/challenges.yaml")

	// Charger la configuration du jeu
	log.Printf("[config] Chargement de la configuration depuis: %s", configPath)
	game, err := LoadGameConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("échec du chargement de la configuration: %w", err)
	}

	// Log de la configuration du jeu
	log.Printf("[config] game: %s v%s", game.Game.Name, game.Game.Version)

	// Validation et log des pondérations de rareté
	sum := 0
	for _, weight := range game.RarityWeights {
		sum += weight
	}
	log.Printf("[config] rarity weights: C=%d R=%d L=%d (OK sum=%d)",
		game.RarityWeights[RarityCommon],
		game.RarityWeights[RarityRare],
		game.RarityWeights[RarityLegendary],
		sum)

	// Log des récompenses XP
	log.Printf("[config] xp rewards: C=%d R=%d L=%d",
		game.XPRewards[RarityCommon],
		game.XPRewards[RarityRare],
		game.XPRewards[RarityLegendary])

	// Log des paramètres du spawner
	if envInterval := os.Getenv("WORDMON_SPAWN_INTERVAL"); envInterval != "" {
		log.Printf("[config] spawner.intervalSeconds=%s (override ENV)", envInterval)
	} else {
		log.Printf("[config] spawner.intervalSeconds=%d", game.Spawner.IntervalSeconds)
	}

	// Charger les défis
	log.Printf("[config] Chargement des défis depuis: %s", challengesPath)
	challenges, err := LoadChallenges(challengesPath)
	if err != nil {
		return nil, fmt.Errorf("échec du chargement des défis: %w", err)
	}
	log.Println("[config] challenges: anagram + a-trou chargés")

	// Charger le dictionnaire de mots
	log.Printf("[config] Chargement du dictionnaire depuis: %s", wordsPath)
	words, err := LoadWords(wordsPath)
	if err != nil {
		return nil, fmt.Errorf("échec du chargement du dictionnaire: %w", err)
	}

	// Compter les mots par rareté
	counts := make(map[string]int)
	for _, word := range words {
		counts[word.Rarity]++
	}
	log.Printf("[config] words: Common=%d Rare=%d Legendary=%d (total=%d)",
		counts[RarityCommon], counts[RarityRare], counts[RarityLegendary], len(words))

	// Vérifier qu'il y a au moins 5 mots par rareté
	for rarityType, count := range counts {
		if count < 5 {
			log.Printf("[config] ⚠️  Attention: seulement %d mots pour la rareté %s (minimum recommandé: 5)", count, rarityType)
		}
	}

	return &GameData{
		Game:           game,
		Challenges:     challenges,
		Words:          words,
		ConfigPath:     configPath,
		WordsPath:      wordsPath,
		ChallengesPath: challengesPath,
	}, nil
}

// ResolvePath résout un chemin relatif en chemin absolu pour les logs
func ResolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return path // retourner le chemin original si erreur
	}
	return absPath
}
