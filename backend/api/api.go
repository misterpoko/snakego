package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"snakego/backend/game"
	"time"
)

type dirCommand struct {
	Snake int
	Dir   int
}

var dirCh = make(chan dirCommand)

var dirNames = map[string]int{
	"up":    game.UP,
	"down":  game.DOWN,
	"left":  game.LEFT,
	"right": game.RIGHT,
}

var addSnakeCh = make(chan struct{})

func createGame(w http.ResponseWriter, r *http.Request) {
	letGameSize := 25

	fmt.Print("\033[?25l")     // hide cursor
	fmt.Print("\033[H\033[2J") // clear screen
	board0 := game.CreateBoard(letGameSize, letGameSize, 0)

	go startGameLoop(board0)
	w.Write([]byte("game started"))
}

func startGameLoop(board *game.Board) {

	ticker := time.NewTicker(167 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {

		case <-addSnakeCh:
			board.Snakes = append(board.Snakes, game.CreateSnake(board))

		case info := <-dirCh:
			game.UpdateUserDir(board, board.Snakes[info.Snake], info.Dir)

		case <-ticker.C:
			game.UpdateGame(board)
			game.PrintBoard(board)
		}
	}
}

func addSnake(w http.ResponseWriter, r *http.Request) {
	addSnakeCh <- struct{}{}
	w.Write([]byte("snake added"))
}

func quit(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\033[?25h")     // show cursor again
	fmt.Print("\033[H\033[2J") // clear screen
	os.Exit(0)
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

	dirCh <- dirCommand{Snake: req.Snake, Dir: dir}
	w.Write([]byte("ok"))
}

func Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /quit", quit)
	mux.HandleFunc("GET /createGame", createGame)
	mux.HandleFunc("POST /addSnake", addSnake)
	mux.HandleFunc("POST /changeDir", changeDir)
	return mux
}

func Start() {
	http.ListenAndServe(":8080", Routes())
}
