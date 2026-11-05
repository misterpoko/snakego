package main

import "snakego/backend/game"
import "time"
import "fmt"


func main() {

	fmt.Print("\033[?25l") // hide cursor
	fmt.Print("\033[H\033[2J") // clear screen
	board := game.CreateBoard(25, 25)
	snake := game.CreateSnake(board)

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
			game.UpdateUserDir(snake, k)

		case <-ticker.C:
			game.UpdateGame(board)
			game.PrintBoard(board)
		}
	}

}
