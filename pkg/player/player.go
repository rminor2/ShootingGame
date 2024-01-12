package player

import (
	"database/sql"
)

type Player struct {
	UserName string
	Age      int
	Score    int
}

func UpdatePlayer(db *sql.DB, player Player) error {
	_, err := db.Exec("INSERT INTO players (username, age, score) VALUES ($1, $2, $3) ON CONFLICT (username) DO UPDATE SET age = EXCLUDED.age, score = EXCLUDED.score", player.UserName, player.Age, player.Score)
	return err
}

func GetLeaderboard(db *sql.DB) ([]Player, error) {
	rows, err := db.Query("SELECT username, age, score FROM players ORDER BY score DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.UserName, &p.Age, &p.Score); err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, rows.Err()
}
