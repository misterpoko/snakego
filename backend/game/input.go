package game

import (
	"golang.org/x/term"
	"os"
)

const Quit = -1

func StartInput() (keys <-chan int, restore func(), err error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, nil, err
	}

	ch := make(chan int, 8)
	go func() {
		buf := make([]byte, 1)
		for {
			if _, err := os.Stdin.Read(buf); err != nil {
				return
			}
			switch buf[0] {
			case 'w', 'W':
				ch <- UP
			case 's', 'S':
				ch <- DOWN
			case 'a', 'A':
				ch <- LEFT
			case 'd', 'D':
				ch <- RIGHT
			case 'q', 3: // q or Ctrl+C
				ch <- Quit
			}
		}
	}()

	return ch, func() { term.Restore(fd, oldState) }, nil
}
