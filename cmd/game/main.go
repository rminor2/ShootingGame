package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/ryanminor/GO_Practice/pkg/database"
	"github.com/ryanminor/GO_Practice/pkg/player"
)

var (
	leaderboardMutex sync.Mutex
	leaderboard      map[string]player.Player
)

// Function to connect to the database
func connectToDatabase() (*sql.DB, error) {
	// Replace these with your actual database credentials and configuration
	const (
		host     = "your-database-host.amazonaws.com"
		port     = 5432 // Default port for PostgreSQL
		user     = "yourDatabaseUser"
		password = "yourDatabasePassword"
		dbname   = "yourDatabaseName"
	)

	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Open a connection to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Player struct represents a player in the game with a username, age, and score.
type GameRequest struct {
	Players   []player.Player `json:"players"`
	NumRounds int             `json:"numRounds"`
}

// LeaderboardManager struct to manage the leaderboard
type LeaderboardManager struct {
	scores map[string]player.Player
	mutex  sync.Mutex
}

// NewLeaderboardManager creates and initializes a new LeaderboardManager
func NewLeaderboardManager() *LeaderboardManager {
	return &LeaderboardManager{
		scores: make(map[string]player.Player),
	}
}

// UpdateLeaderboard updates the leaderboard with a slice of players
// UpdateLeaderboard updates the leaderboard with a slice of players
// UpdateLeaderboard updates the leaderboard with a slice of players
func (lm *LeaderboardManager) UpdateLeaderboard(db *sql.DB, players []player.Player) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	for _, p := range players {
		// Update in-memory leaderboard
		lm.scores[p.UserName] = p

		// Update database leaderboard
		_, err := db.Exec("INSERT INTO leaderboard (username, score) VALUES ($1, $2) ON CONFLICT (username) DO UPDATE SET score = EXCLUDED.score", p.UserName, p.Score)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetLeaderboard returns the current state of the leaderboard
func (lm *LeaderboardManager) GetLeaderboard() []player.Player {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	var leaderboard []player.Player
	for _, p := range lm.scores {
		leaderboard = append(leaderboard, p)
	}
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Score > leaderboard[j].Score
	})
	return leaderboard
}

// takeShot simulates a shooting action in the game, randomly returning 0 (miss) or 3 (hit).
func takeShot() int {
	if rand.Intn(2) == 0 {
		return 0 // Missed shot
	}
	return 3 // Successful shot
}

// clearBuffer clears the input buffer in case of invalid input.
// This is necessary to handle cases where Scan functions leave a newline character in the buffer.
func clearBuffer() {
	var discard string
	fmt.Scanln(&discard)
}

// playGame is a goroutine function that handles an individual game
func playGame(gameReq GameRequest, doneChan chan<- []player.Player) {
	players := playRound(gameReq)
	calculateScores(players)
	doneChan <- players
}

// initializePlayers sets up the initial state of each player.
// It takes the number of players as input, and prompts for each player's name and age.
// Returns a slice of initialized Player structs.
func initializePlayers(numPlayers int) []player.Player {
	players := make([]player.Player, numPlayers)
	reader := bufio.NewReader(os.Stdin)

	for i := 0; i < numPlayers; i++ {
		var name string
		// Loop until a valid name is entered.
		for name == "" {
			fmt.Printf("Enter the name of player %d: ", i+1)
			nameInput, _ := reader.ReadString('\n')
			name = strings.TrimSpace(nameInput)
			if name == "" {
				fmt.Println("Name cannot be empty. Please enter a valid name.")
			}
		}

		var age int
		// Loop until a valid age is entered.
		for {
			fmt.Printf("Enter the age of player %s: ", name)
			_, err := fmt.Scan(&age)
			if err != nil || age <= 0 {
				fmt.Println("Invalid age. Please enter a valid positive integer.")
				clearBuffer()
				continue
			}
			break
		}

		players[i] = player.Player{UserName: name, Age: age}
	}

	return players
}

// playRound handles the logic for each round of the game.
// It loops through the specified number of rounds, and within each round,
// it allows each player to take shots and updates their scores.
// Returns the updated slice of players with their final scores.
func playRound(gameReq GameRequest) []player.Player {
	reader := bufio.NewReader(os.Stdin)

	for round := 1; round <= gameReq.NumRounds; round++ {
		fmt.Printf("\n--- Round %d ---\n", round)
		for miniRound := 1; miniRound <= 3; miniRound++ {
			fmt.Printf("\nMini Round %d\n", miniRound)
			for shot := 1; shot <= 3; shot++ {
				for i := range gameReq.Players {
					fmt.Printf("\n%s, it's your turn to shoot! Press Enter to shoot...\n", gameReq.Players[i].UserName)
					reader.ReadString('\n')
					points := takeShot()
					gameReq.Players[i].Score += points
					if points > 0 {
						fmt.Println("\nNice shot! You scored 3 points!")
					} else {
						fmt.Println("\nDarn, you missed. Sorry!")
					}
				}
			}
		}
	}
	return gameReq.Players
}

// calculateScores sorts the players based on their scores in descending order.
func calculateScores(players []player.Player) {
	// Add logic to calculate and update player scores
	// For now, it's just the sum of shots

	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})
}

// displayScores prints the leaderboard showing the players and their scores.
func displayScores(players []player.Player) {
	fmt.Println("\nLeader Board: ")
	for i, player := range players {
		fmt.Printf("\nPlayer %d: Name: %s Age: %d Total Score: %d \n", i+1, player.UserName, player.Age, player.Score)
	}
}

func startGameHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, lm *LeaderboardManager) {

	if r.Method != "POST" {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var gameReq GameRequest
	if err := json.NewDecoder(r.Body).Decode(&gameReq); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	players := playRound(gameReq)
	for _, p := range players {
		if err := player.UpdatePlayer(db, p); err != nil {
			http.Error(w, "Error updating player", http.StatusInternalServerError)
			return
		}
	}

	if err := lm.UpdateLeaderboard(db, players); err != nil {
		http.Error(w, "Error updating leaderboard", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Game started and data recorded")
}

func updateLeaderboard(players []player.Player) {
	leaderboardMutex.Lock()
	defer leaderboardMutex.Unlock()

	for _, player := range players {
		// Assuming UserName is unique

		leaderboard[player.UserName] = player
	}
}

// main is the entry point of the program.
// It seeds the random number generator, prompts for the number of players and rounds,
// and orchestrates the flow of the game.
func main() {
	rand.Seed(time.Now().UnixNano())

	db, err := database.ConnectToDatabase()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	lm := NewLeaderboardManager()

	http.HandleFunc("/start-game", func(w http.ResponseWriter, r *http.Request) {
		startGameHandler(w, r, db, lm)
	})

	http.HandleFunc("/leaderboard", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
			return
		}
		leaderboard, err := player.GetLeaderboard(db)
		if err != nil {
			http.Error(w, "Error fetching leaderboard", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(leaderboard)
	})

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
