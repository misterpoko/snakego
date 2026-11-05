package game

import "fmt"
import "math/rand/v2"


type Coord struct {
	X, Y int
}

type Board struct {
	Width int
	Height int
	Grid [][]int
	Snakes []*Snake
	Updates []*Coord
}


type Snake struct {
	Head Coord 
	Tail Coord 
	UserDir int
	SnakeSize int
}


const (
	OBJECT_TYPE_FLAG = 0b111
	NOTHING = 	0b000
	SNAKE = 	0b001
	APPLE = 	0b010
	HALF_ME = 	0b011
	DECAY_TAIL = 	0b100
)

const (
	DIRECTION_SHIFT = 3
	DIRECTION_TYPE_FLAG = 0b11 << DIRECTION_SHIFT
	UP = 	0b00 << DIRECTION_SHIFT
	DOWN = 	0b11 << DIRECTION_SHIFT
	LEFT = 	0b10 << DIRECTION_SHIFT
	RIGHT = 0b01 << DIRECTION_SHIFT
)



func CreateBoard(width, height int) *Board {
	grid := make([][]int, height)
	for i := range(grid) {
		grid[i] = make([]int, width)
	}



	board := Board{
		Width: width,
		Height: height,
		Grid: grid,
	}

	for range 20 {
		createObject(&board, APPLE)
	}

	createObject(&board, HALF_ME)


	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			board.Updates = append(board.Updates, &Coord{j, i})
		}
	}
	printBoarder(&board)

	return &board
}

func PrintBoard(board *Board) {
	for _, update := range board.Updates {
		printSpot(board, *update)
	}
	board.Updates = nil
}


func UpdateUserDir(snake *Snake, dir int) {
	currentDir := snake.UserDir
	if currentDir ^ dir == DIRECTION_TYPE_FLAG {
		return
	}
	snake.UserDir = dir

}


func updateSnakeDir(board *Board, coord *Coord) {
	prevInfo := getSpot(board, *coord) & DIRECTION_TYPE_FLAG
	switch(prevInfo) {
	case LEFT:
		coord.X -= 1
	case RIGHT:
		coord.X += 1
	case UP:
		coord.Y -= 1
	case DOWN:
		coord.Y += 1
	}

	coord.X = (coord.X + board.Width) % board.Width
	coord.Y = (coord.Y + board.Height) % board.Height
}


func cutHalf(board *Board, snake *Snake) int {
	halfSnake := snake.SnakeSize / 2
	for range halfSnake {
		setSpot(board, snake.Tail, getSpot(board, snake.Tail) &^ OBJECT_TYPE_FLAG | DECAY_TAIL)
		updateSnakeDir(board, &snake.Tail)
	}
	return snake.SnakeSize - halfSnake

}


func UpdateGame(board *Board) {
	for _, snake := range board.Snakes {
		setSpot(board, snake.Head, getSpot(board, snake.Head) &^ DIRECTION_TYPE_FLAG | snake.UserDir)
		updateSnakeDir(board, &snake.Head)

		target := getSpot(board, snake.Head) & OBJECT_TYPE_FLAG

		switch (target) {
		case SNAKE:
			fmt.Print("\033[?25h")
			panic("you lose")
		case APPLE:
			snake.SnakeSize++
			createObject(board, APPLE)
		case HALF_ME:
			snake.SnakeSize = cutHalf(board, snake)
			createObject(board, HALF_ME)
			fallthrough
		default:
			tailPos := snake.Tail
			updateSnakeDir(board, &snake.Tail)
			setSpot(board, tailPos, NOTHING)
		}

		setSpot(board, snake.Head, SNAKE | snake.UserDir)
	}
}


func printBoarder(board *Board) {
	boarderValue := "\u2588"

	for i := 0; i < board.Height + 2; i++ {
		fmt.Printf("\033[%d;%dH%s", i, 0, boarderValue)
		fmt.Printf("\033[%d;%dH%s", i, board.Width * 2 + 1, boarderValue)
	}

	for i := 0; i < board.Width + 1; i++ {
		fmt.Printf("\033[%d;%dH%s", 0, i * 2, boarderValue + boarderValue)
		fmt.Printf("\033[%d;%dH%s", board.Height + 2, i * 2, boarderValue + boarderValue)
	}

}


func printSpot(board *Board, coord Coord) {
	boardValue := getSpot(board, coord)
	toPrint := " "
	switch(boardValue & OBJECT_TYPE_FLAG) {
	case SNAKE:
		toPrint = "\033[38;2;0;255;0m\u25A0\033[0m" 
	case APPLE:
		toPrint = "\033[31m\u25CF\033[0m"
	case HALF_ME:
		toPrint = "\u25D0"
	case DECAY_TAIL:
		toPrint = "\033[38;2;139;69;19m\u25A0\033[0m"

	}
	fmt.Printf("\033[%d;%dH%s", coord.Y + 2, coord.X * 2 + 2, toPrint)
}


func setSpot(board *Board, coord Coord, value int) {
	board.Grid[coord.Y][coord.X] = value
	board.Updates = append(board.Updates, &coord);
}


func getSpot(board *Board, coord Coord) int {
	return board.Grid[coord.Y][coord.X]
}


func availableSpot(board *Board, coord Coord) bool {
	return board.Grid[coord.Y][coord.X] == 0
}

func createObject(board *Board, givenType int) {
	potentialSpot := Coord{rand.IntN(board.Width), rand.IntN(board.Height)}
	for !availableSpot(board, potentialSpot) {
		potentialSpot = Coord{rand.IntN(board.Width), rand.IntN(board.Height)}
	}
	setSpot(board, potentialSpot, givenType)
}


func CreateSnake(board *Board) *Snake {
	potentialSnake := Coord{rand.IntN(board.Width), rand.IntN(board.Height)}
	for !availableSpot(board, potentialSnake) {
		potentialSnake = Coord{rand.IntN(board.Width), rand.IntN(board.Height)}
	}
	setSpot(board, potentialSnake, SNAKE | RIGHT)
	newSnake := &Snake{
		Head: potentialSnake,
		Tail: potentialSnake,
		UserDir: RIGHT,
		SnakeSize: 1,
	}
	board.Snakes = append(board.Snakes, newSnake)
	return newSnake
}


