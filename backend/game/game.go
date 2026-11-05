package game

import "fmt"
import "math/rand/v2"

const META_GAME_GRID_SIZE = 2

type Coord struct {
	X, Y int
}

type Board struct {
	Width int
	Height int
	Grid [][]int
	Snakes []*Snake
	Updates []*Coord
	TileNumber int
}


type Snake struct {
	Head Coord 
	Tail Coord 
	UserDir int
	SnakeSize int
}


//0-2
const (
	OBJECT_TYPE_FLAG = 0b111
	NOTHING = 	0b000
	SNAKE = 	0b001
	APPLE = 	0b010
	HALF_ME = 	0b011
	DECAY_TAIL = 	0b100
	PORTAL_ORANGE = 	0b101
	PORTAL_BLUE = 	0b110
)

//3-5
const (
	DIRECTION_SHIFT = 3
	DIRECTION_TYPE_FLAG = 0b11 << DIRECTION_SHIFT
	UP = 	0b00 << DIRECTION_SHIFT
	DOWN = 	0b11 << DIRECTION_SHIFT
	LEFT = 	0b10 << DIRECTION_SHIFT
	RIGHT = 0b01 << DIRECTION_SHIFT
)

//8-24
const (
	COORD_X_SHIFT = 16
	COORD_Y_SHIFT = 8
	COORD_FLAG = 0b1111_1111
	RESETCOORD_FLAG = COORD_FLAG<<COORD_X_SHIFT | COORD_FLAG<<COORD_Y_SHIFT

)



func CreateBoard(width, height, tileNumber int) *Board {
	grid := make([][]int, height)
	for i := range(grid) {
		grid[i] = make([]int, width)
	}



	board := Board{
		Width: width,
		Height: height,
		Grid: grid,
		TileNumber: tileNumber,
	}

	for range width/5 {
		createObject(&board, APPLE)
	}

	createObject(&board, HALF_ME)

	createPortal(&board)


	for i := range height {
		for j := range width {
			board.Updates = append(board.Updates, &Coord{j, i})
		}
	}
	printBoarder(&board)

	return &board
}

func createPortal(board *Board) {
	bluePortal := createObject(board, PORTAL_ORANGE)
	orangePortal := createObject(board, PORTAL_BLUE)
	orangeX, orangeY := orangePortal.X, orangePortal.Y
	blueX, blueY := bluePortal.X, bluePortal.Y
	setSpot(board, *bluePortal, getSpot(board, *bluePortal)|orangeX<<COORD_X_SHIFT|orangeY<<COORD_Y_SHIFT)
	setSpot(board, *orangePortal, getSpot(board, *orangePortal)|blueX<<COORD_X_SHIFT|blueY<<COORD_Y_SHIFT)
}

func PrintBoard(board *Board) {
	for _, update := range board.Updates {
		printSpot(board, *update)
	}
	board.Updates = nil
}


func UpdateUserDir(board *Board, snake *Snake, dir int) {
	currentDir := snake.UserDir
	if currentDir ^ dir == DIRECTION_TYPE_FLAG {
		return
	}
	snake.UserDir = dir


	nextHeadXPos := snake.Head.X
	nextHeadYPos := snake.Head.Y

	switch(snake.UserDir) {
	case LEFT:
		nextHeadXPos -= 1
	case RIGHT:
		nextHeadXPos += 1
	case UP:
		nextHeadYPos -= 1
	case DOWN:
		nextHeadYPos += 1
	}

	nextHeadXPos = (nextHeadXPos + board.Width) % board.Width
	nextHeadYPos = (nextHeadYPos + board.Height) % board.Height

	setNextLocation(board, &snake.Head, nextHeadXPos, nextHeadYPos)

}


func setNextLocation(board *Board, tile *Coord, x, y int) {
	setSpot(board, *tile, getSpot(board, *tile) &^ RESETCOORD_FLAG)
	setSpot(board, *tile, getSpot(board, *tile) | x << COORD_X_SHIFT| y << COORD_Y_SHIFT)
}


func updateSnakeTile(board *Board, coord *Coord) {
	prevInfo := getSpot(board, *coord)
	newX := (prevInfo >> COORD_X_SHIFT) & COORD_FLAG
	newY := (prevInfo >> COORD_Y_SHIFT) & COORD_FLAG
	coord.X = newX
	coord.Y = newY
}


func cutHalf(board *Board, snake *Snake) int {
	halfSnake := snake.SnakeSize / 2
	for range halfSnake {
		setSpot(board, snake.Tail, getSpot(board, snake.Tail) &^ OBJECT_TYPE_FLAG | DECAY_TAIL)
		updateSnakeTile(board, &snake.Tail)
	}
	return snake.SnakeSize - halfSnake

}


func UpdateGame(board *Board) {
	for _, snake := range board.Snakes {
		previousSpot := Coord{snake.Head.X, snake.Head.Y}
		UpdateUserDir(board, snake, snake.UserDir)
		updateSnakeTile(board, &snake.Head)
		target := getSpot(board, snake.Head) & OBJECT_TYPE_FLAG

		switch (target) {
		case SNAKE:
			fmt.Print("\033[?25h")
			panic("you lose")
		case APPLE:
			snake.SnakeSize++
			createObject(board, APPLE)
		case DECAY_TAIL:
			shrinkSnake(board, snake)
			if snake.SnakeSize == 1  {
				break
			}
			snake.SnakeSize--
			shrinkSnake(board, snake)

		case PORTAL_BLUE, PORTAL_ORANGE:
			portalInfo := getSpot(board, snake.Head)
			otherPortalX := (portalInfo >> COORD_X_SHIFT) & COORD_FLAG 
			otherPortalY := (portalInfo >> COORD_Y_SHIFT) & COORD_FLAG 
			switch(snake.UserDir) {
			case LEFT:
				otherPortalX -= 1
			case RIGHT:
				otherPortalX += 1
			case UP:
				otherPortalY -= 1
			case DOWN:
				otherPortalY += 1
			}

			otherPortalX = (otherPortalX + board.Width) % board.Width
			otherPortalY = (otherPortalY + board.Height) % board.Height
			snake.Head.X = otherPortalX
			snake.Head.Y = otherPortalY
			setNextLocation(board, &previousSpot, otherPortalX, otherPortalY)
			setSpot(board, snake.Head, SNAKE)
			shrinkSnake(board, snake)

		case HALF_ME:
			snake.SnakeSize = cutHalf(board, snake)
			createObject(board, HALF_ME)
			shrinkSnake(board, snake)

		default:
			shrinkSnake(board, snake)
		}
		setSpot(board, snake.Head, SNAKE)
	}
}

func shrinkSnake(board *Board, snake *Snake) {
	tailPos := snake.Tail
	updateSnakeTile(board, &snake.Tail)
	setSpot(board, tailPos, NOTHING)
}


func printBoarder(board *Board) {
	boarderValue := "\u2588"

	offSetY := 1 + board.TileNumber/META_GAME_GRID_SIZE*(board.Height+2)
	offsetX := 1 + board.TileNumber%META_GAME_GRID_SIZE*(board.Width*2+2)

	for i := 0; i < board.Height+2; i++ {
		fmt.Printf("\033[%d;%dH%s", i+offSetY, offsetX, boarderValue)
		fmt.Printf("\033[%d;%dH%s", i+offSetY, offsetX+board.Width*2+1, boarderValue)
	}

	for i := 0; i < board.Width+1; i++ {
		fmt.Printf("\033[%d;%dH%s", offSetY, offsetX+i*2, boarderValue+boarderValue)
		fmt.Printf("\033[%d;%dH%s", offSetY+board.Height+1, offsetX+i*2, boarderValue+boarderValue)
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
	case PORTAL_ORANGE:
		toPrint = "\033[38;2;255;140;0m\u25C9\033[0m"   // orange ◉
	case PORTAL_BLUE:
		toPrint = "\033[38;2;0;120;255m\u25C9\033[0m"   // blue ◉

	}

	offSetY := 2 + board.TileNumber/META_GAME_GRID_SIZE*(board.Height+2)
	offsetX := 2 + board.TileNumber%META_GAME_GRID_SIZE*(board.Width*2+2)

	fmt.Printf("\033[%d;%dH%s", coord.Y+offSetY, coord.X*2+offsetX, toPrint)
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

func createObject(board *Board, givenType int) *Coord {
	potentialSpot := Coord{rand.IntN(board.Width), rand.IntN(board.Height)}
	for !availableSpot(board, potentialSpot) {
		potentialSpot = Coord{rand.IntN(board.Width), rand.IntN(board.Height)}
	}
	setSpot(board, potentialSpot, givenType)
	return &potentialSpot
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


