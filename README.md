BBGame - A Simple Shooting Game
Welcome to BBGame, a simple, interactive shooting game written in Go. This game allows multiple players to participate in a series of shooting rounds, accumulating scores based on their performance.

Installation
Before you run the game, ensure you have Go installed on your system. If you don't have Go installed, you can download it from the official Go website.

To set up the game, clone the repository or download the source code to your local machine.

Running the Game
Navigate to the directory containing the game's source code in your terminal. Run the game using the Go command: go run main.go


Game Rules - 
- At the start, the game prompts you to enter the number of players.
- Each player is then asked to enter their name and age.
- Next, you'll be prompted to enter the number of rounds you wish to play. Each game consists of 3 mini-rounds.
- In each mini-round, each player gets a chance to take a shot.
- A successful shot scores 3 points, and a miss scores 0.
- After all rounds are complete, the scores are tallied and displayed in a leaderboard.


Gameplay - 
- The game is turn-based, with each player taking a shot in their turn.
- The score is calculated based on the number of successful shots.
- At the end of all rounds, players are ranked according to their total scores.


Notes - 
- The game uses a simple command-line interface for interaction.
- Ensure that inputs (player names, number of players, and rounds) are entered as prompted to avoid errors.



