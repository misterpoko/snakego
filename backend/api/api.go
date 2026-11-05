package api



import (
	"net/http"
	"snakego/backend/game"
	"fmt"
	"time"
)

func createGame(w http.ResponseWriter, r *http.Request) {
	letGameSize := 25

	fmt.Print("\033[?25l") // hide cursor
	fmt.Print("\033[H\033[2J") // clear screen
	board0 := game.CreateBoard(letGameSize, letGameSize, 0)
	snake0 := game.CreateSnake(board0)



	go startGameLoop(board0, snake0)
}

func startGameLoop(board *game.Board, snake *game.Snake) {

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
			game.UpdateUserDir(board, snake, k)

		case <-ticker.C:
			game.UpdateGame(board)
			game.PrintBoard(board)
		}
	}
}


func Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /createGame", createGame)
	return mux
}


func Start() {
	http.ListenAndServe(":8080", Routes())
}

