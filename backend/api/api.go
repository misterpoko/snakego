package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"snakego/backend/game"
	"time"
)

type dirCommand struct {
	Snake int
	Dir   int
}

var gameStateInfo *game.Board

var dirNames = map[string]int{
	"up":    game.UP,
	"down":  game.DOWN,
	"left":  game.LEFT,
	"right": game.RIGHT,
}

func createGame(w http.ResponseWriter, r *http.Request) {
	letGameSize := 25

	fmt.Print("\033[?25l")     // hide cursor
	fmt.Print("\033[H\033[2J") // clear screen
	board := game.CreateBoard(letGameSize, letGameSize, 0)

	go startGameLoop(board)
	gameStateInfo = board
	w.Write([]byte("game started"))
}

func startGameLoop(board *game.Board) {

	ticker := time.NewTicker(334 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		game.UpdateGame(board)
		game.PrintBoard(board)
	}
}

func addSnake(w http.ResponseWriter, r *http.Request) {
	snake := game.CreateSnake(gameStateInfo)
	idx := slices.Index(gameStateInfo.Snakes, snake)
	snakeNumber, err := json.Marshal(dirCommand{Snake: idx, Dir: -1})
	if err != nil {
		fmt.Println("trouble adding snake")
	}
	w.Write([]byte(snakeNumber))
}

func destroyCurrentGame(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\033[?25h")     // show cursor again
	fmt.Print("\033[H\033[2J") // clear screen
	gameStateInfo.Snakes = nil
	gameStateInfo.Grid = nil
	gameStateInfo.Updates = nil
}

func changeDir(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Snake int    `json:"snake"`
		Dir   string `json:"dir"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	dir, ok := dirNames[req.Dir]
	if !ok {
		http.Error(w, "dir must be up, down, left or right", http.StatusBadRequest)
		return
	}

	if req.Snake >= len(gameStateInfo.Snakes) {
		http.Error(w, "snake isn't a valid player", http.StatusBadRequest)
		return
	}

	game.UpdateUserDir(gameStateInfo, gameStateInfo.Snakes[req.Snake], dir)

	w.Write([]byte("ok"))
}

func gameInfo(w http.ResponseWriter, r *http.Request) {
	info, err := json.Marshal(gameStateInfo)
	if err != nil {
		fmt.Println("something bad happen")
	}

	w.Write([]byte(info))
}

func Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /quitAll", destroyCurrentGame)
	mux.HandleFunc("GET /createGame", createGame)
	mux.HandleFunc("GET /addSnake", addSnake)
	mux.HandleFunc("POST /changeDir", changeDir)
	mux.HandleFunc("GET /gameInfo", gameInfo)
	return mux
}

func Start() {
	http.ListenAndServe(":8080", Routes())
}
