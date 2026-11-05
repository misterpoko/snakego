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
			if k == game.Quit {
				fmt.Print("\033[?25h")   // show cursor again
				return
			}
			game.UpdateUserDir(board0, snake0, k)

		case <-ticker.C:
			game.UpdateGame(board0)
			game.PrintBoard(board0)
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

