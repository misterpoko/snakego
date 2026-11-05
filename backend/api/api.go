package api



import (
	"net/http"
	"snakego/backend/game"
	"fmt"
	"time"
	"os"
)

var addSnakeCh = make(chan struct{})

func createGame(w http.ResponseWriter, r *http.Request) {
	letGameSize := 25

	fmt.Print("\033[?25l") // hide cursor
	fmt.Print("\033[H\033[2J") // clear screen
	board0 := game.CreateBoard(letGameSize, letGameSize, 0)

	go startGameLoop(board0)
	w.Write([]byte("game started"))
}

func startGameLoop(board *game.Board) {

	keys, restore, err := game.StartInput()
	if err != nil {
		panic(err)
	}
	defer restore()

	ticker := time.NewTicker(167 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case k := <-keys:
			for _, snake := range board.Snakes {
				game.UpdateUserDir(board, snake, k)
			}

		case <-addSnakeCh:
			board.Snakes = append(board.Snakes, game.CreateSnake(board))

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
	fmt.Print("\033[?25h") // show cursor again
	fmt.Print("\033[H\033[2J") // clear screen
	os.Exit(0)
}




func Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /quit", quit)
	mux.HandleFunc("GET /createGame", createGame)
	mux.HandleFunc("POST /addSnake", addSnake)
	return mux
}


func Start() {
	http.ListenAndServe(":8080", Routes())
}

