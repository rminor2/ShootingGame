package game

import (
	"database/sql"
)

type Game struct {
	ID             int
	PlayerUsername string
	Score          int
}

func RecordGame(db *sql.DB, game Game) error {
	_, err := db.Exec("INSERT INTO games (player_username, score) VALUES ($1, $2)", game.PlayerUsername, game.Score)
	return err
}

