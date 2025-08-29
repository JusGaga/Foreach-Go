package api

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jusgaga/wordmon-go/internal/core"
	_ "github.com/lib/pq"
)

// SQLStore implémente le stockage SQL pour PostgreSQL
type SQLStore struct {
	db *sql.DB
}

// NewSQLStore crée un nouveau store SQL
func NewSQLStore(databaseURL string) (*SQLStore, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("erreur connexion DB: %w", err)
	}

	// Tester la connexion
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erreur ping DB: %w", err)
	}

	log.Printf("[db] Connected to Postgres")
	return &SQLStore{db: db}, nil
}

// Close ferme la connexion à la base de données
func (s *SQLStore) Close() error {
	return s.db.Close()
}

// CreatePlayer crée un nouveau joueur
func (s *SQLStore) CreatePlayer(name string) (*PlayerResponse, error) {
	playerID := uuid.New().String()

	query := `INSERT INTO players (id, name, xp, level) VALUES ($1, $2, $3, $4)`
	_, err := s.db.Exec(query, playerID, name, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("erreur création joueur: %w", err)
	}

	player := &PlayerResponse{
		ID:        playerID,
		Name:      name,
		XP:        0,
		Level:     1,
		Inventory: make(map[string]int),
	}

	log.Printf("[api] Player created: %s (id=%s)", name, playerID)
	return player, nil
}

// GetPlayer récupère un joueur par son ID
func (s *SQLStore) GetPlayer(id string) (*PlayerResponse, error) {
	query := `SELECT id, name, xp, level FROM players WHERE id = $1`

	var player PlayerResponse
	err := s.db.QueryRow(query, id).Scan(&player.ID, &player.Name, &player.XP, &player.Level)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &PlayerNotFoundError{ID: id}
		}
		return nil, fmt.Errorf("erreur récupération joueur: %w", err)
	}

	// Récupérer l'inventaire
	inventory, err := s.ListByPlayer(id)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération inventaire: %w", err)
	}

	player.Inventory = make(map[string]int)
	for _, word := range inventory {
		player.Inventory[word.Text]++
	}

	return &player, nil
}

// GetAllPlayers récupère tous les joueurs
func (s *SQLStore) GetAllPlayers() []*PlayerResponse {
	query := `SELECT id, name, xp, level FROM players ORDER BY xp DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("erreur récupération joueurs: %v", err)
		return []*PlayerResponse{}
	}
	defer rows.Close()

	var players []*PlayerResponse
	for rows.Next() {
		var player PlayerResponse
		if err := rows.Scan(&player.ID, &player.Name, &player.XP, &player.Level); err != nil {
			log.Printf("erreur scan joueur: %v", err)
			continue
		}
		player.Inventory = make(map[string]int)
		players = append(players, &player)
	}

	return players
}

// UpdatePlayer met à jour un joueur
func (s *SQLStore) UpdatePlayer(player *PlayerResponse) error {
	query := `UPDATE players SET name = $1, xp = $2, level = $3 WHERE id = $4`

	result, err := s.db.Exec(query, player.Name, player.XP, player.Level, player.ID)
	if err != nil {
		return fmt.Errorf("erreur mise à jour joueur: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur vérification rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return &PlayerNotFoundError{ID: player.ID}
	}

	return nil
}

// UpdateXP met à jour l'XP et le niveau d'un joueur
func (s *SQLStore) UpdateXP(id string, newXP int, newLevel int) error {
	query := `UPDATE players SET xp = $1, level = $2 WHERE id = $3`

	result, err := s.db.Exec(query, newXP, newLevel, id)
	if err != nil {
		return fmt.Errorf("erreur mise à jour XP: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur vérification rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return &PlayerNotFoundError{ID: id}
	}

	return nil
}

// GetPlayerCount retourne le nombre de joueurs
func (s *SQLStore) GetPlayerCount() int {
	query := `SELECT COUNT(*) FROM players`

	var count int
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("erreur comptage joueurs: %v", err)
		return 0
	}

	return count
}

// GetStartTime retourne l'heure de démarrage (toujours maintenant pour SQL)
func (s *SQLStore) GetStartTime() time.Time {
	return time.Now()
}

// AddSpawn ajoute un spawn à l'historique (non persisté en SQL pour l'instant)
func (s *SQLStore) AddSpawn(spawn interface{}) error {
	// Pour l'instant, on ne persiste pas les spawns
	return nil
}

// GetCurrentSpawn récupère le spawn actuel (généré dynamiquement)
func (s *SQLStore) GetCurrentSpawn() interface{} {
	// Générer un spawn aléatoire depuis la base
	word, err := s.RandomByRarity("Common")
	if err != nil {
		log.Printf("erreur génération spawn: %v", err)
		return nil
	}

	return core.SpawnEvent{
		Round: 1,
		Word:  *word,
	}
}

// Seed insère les mots dans la base de données
func (s *SQLStore) Seed(words []core.Word) error {
	// Vider la table words d'abord
	_, err := s.db.Exec("DELETE FROM words")
	if err != nil {
		return fmt.Errorf("erreur vidage table words: %w", err)
	}

	// Insérer les nouveaux mots
	query := `INSERT INTO words (id, text, rarity, points) VALUES ($1, $2, $3, $4)`

	for _, word := range words {
		_, err := s.db.Exec(query, word.ID, word.Text, word.Rarity, word.Points)
		if err != nil {
			return fmt.Errorf("erreur insertion mot %s: %w", word.Text, err)
		}
	}

	log.Printf("[seed] %d words loaded into DB", len(words))
	return nil
}

// Get récupère un mot par son ID
func (s *SQLStore) Get(id string) (*core.Word, error) {
	query := `SELECT id, text, rarity, points FROM words WHERE id = $1`

	var word core.Word
	err := s.db.QueryRow(query, id).Scan(&word.ID, &word.Text, &word.Rarity, &word.Points)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("mot non trouvé: %s", id)
		}
		return nil, fmt.Errorf("erreur récupération mot: %w", err)
	}

	return &word, nil
}

// RandomByRarity récupère un mot aléatoire par rareté
func (s *SQLStore) RandomByRarity(rarity string) (*core.Word, error) {
	query := `SELECT id, text, rarity, points FROM words WHERE rarity = $1 ORDER BY RANDOM() LIMIT 1`

	var word core.Word
	err := s.db.QueryRow(query, rarity).Scan(&word.ID, &word.Text, &word.Rarity, &word.Points)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("aucun mot trouvé pour la rareté: %s", rarity)
		}
		return nil, fmt.Errorf("erreur récupération mot aléatoire: %w", err)
	}

	log.Printf("[api] Spawned WordMon %s: %s", word.Rarity, word.Text)
	return &word, nil
}

// Add ajoute une capture
func (s *SQLStore) Add(playerId, wordId string) error {
	captureID := uuid.New().String()

	query := `INSERT INTO captures (id, player_id, word_id) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, captureID, playerId, wordId)
	if err != nil {
		return fmt.Errorf("erreur ajout capture: %w", err)
	}

	return nil
}

// ListByPlayer récupère tous les mots capturés par un joueur
func (s *SQLStore) ListByPlayer(playerId string) ([]core.Word, error) {
	query := `
		SELECT w.id, w.text, w.rarity, w.points 
		FROM words w 
		JOIN captures c ON w.id = c.word_id 
		WHERE c.player_id = $1
		ORDER BY c.captured_at DESC
	`

	rows, err := s.db.Query(query, playerId)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération captures: %w", err)
	}
	defer rows.Close()

	var words []core.Word
	for rows.Next() {
		var word core.Word
		if err := rows.Scan(&word.ID, &word.Text, &word.Rarity, &word.Points); err != nil {
			log.Printf("erreur scan capture: %v", err)
			continue
		}
		words = append(words, word)
	}

	return words, nil
}

// GetLeaderboard récupère le leaderboard des joueurs
func (s *SQLStore) GetLeaderboard(limit int) ([]*PlayerResponse, error) {
	query := `SELECT id, name, xp, level FROM players ORDER BY xp DESC LIMIT $1`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération leaderboard: %w", err)
	}
	defer rows.Close()

	var players []*PlayerResponse
	for rows.Next() {
		var player PlayerResponse
		if err := rows.Scan(&player.ID, &player.Name, &player.XP, &player.Level); err != nil {
			log.Printf("erreur scan leaderboard: %v", err)
			continue
		}
		player.Inventory = make(map[string]int)
		players = append(players, &player)
	}

	return players, nil
}
