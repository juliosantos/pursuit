package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

type Matrix struct {
	width  int
	height int
}

type Position struct {
	x int
	y int
}

type Mover struct {
	position Position
	role     string
}

func newMover(role string) *Mover {
	mover := new(Mover)
	mover.role = role
	mover.position.x = rand.Intn(matrix.width)
	mover.position.y = rand.Intn(matrix.height)
	return mover
}

func (mover *Mover) moveRandom() {
	p := rand.Intn(2) // move on x or y?
	x := rand.Intn(2)
	y := rand.Intn(2)
	if p == 1 {
		if (x == 1 || mover.position.x == 0) && mover.position.x < matrix.width {
			mover.position.x++
		} else {
			mover.position.x--
		}
	} else {
		if (y == 1 || mover.position.y == 0) && mover.position.y < matrix.height {
			mover.position.y++
		} else {
			mover.position.y--
		}
	}
}

func (mover1 *Mover) overlapping(mover2 *Mover) bool {
	return mover1.position == mover2.position
}

func (mover1 *Mover) canSee(mover2 *Mover) bool {
	return mover1.position.x == mover2.position.x || mover1.position.y == mover2.position.y
}

// moves chaser closer to fugitive, in x or y
func (chaser *Mover) pursue(fugitive *Mover) {
	if chaser.position.x == fugitive.position.x {
		pursueInLine(&chaser.position.y, fugitive.position.y)
	} else {
		pursueInLine(&chaser.position.x, fugitive.position.x)
	}
}

// moves chaser, in x or y, to reduce its distance from fugitive
func pursueInLine(chaserCoord *int, fugitiveCoord int) {
	if *chaserCoord > fugitiveCoord {
		*chaserCoord--
	} else {
		*chaserCoord++
	}
}

// moves fugitive away from chaser, in x or y
func (fugitive *Mover) escape(chaser *Mover) {
	if chaser.position.x == fugitive.position.x {
		escapeInLine(&fugitive.position.x, chaser.position.x, matrix.height)
	} else {
		escapeInLine(&fugitive.position.y, chaser.position.y, matrix.width)
	}
}

// moves chaser, in x or y, to increase its distance from fugitive
func escapeInLine(fugitiveCoord *int, chaserCoord int, maxCoord int) {
	if *fugitiveCoord > chaserCoord && *fugitiveCoord < maxCoord {
		*fugitiveCoord++
	} else if *fugitiveCoord > 0 {
		*fugitiveCoord--
	}
}

func fugitiveMove(fugitive *Mover, chasers []*Mover) {
	moved := false
	for _, chaser := range chasers {
		if fugitive.canSee(chaser) {
			fugitive.escape(chaser)
			moved = true
			break
		}
	}
	if !moved {
		fugitive.moveRandom()
	}
}

func chaserMoves(chasers []*Mover, fugitive *Mover) int {
	for _, chaser := range chasers {
		if chaser.canSee(fugitive) {
			chaser.pursue(fugitive)
			if chaser.overlapping(fugitive) {
				return 1
			}
		} else {
			chaser.moveRandom()
		}
	}

	return 0
}

func pursuit(chasers []*Mover, fugitive *Mover) int {
	state(chasers, fugitive)
	for moves := 1; ; moves++ {
		fugitiveMove(fugitive, chasers)
		if chaserMoves(chasers, fugitive) == 1 {
			// fugitive was caught
			state(chasers, fugitive)
			fmt.Println("Caught at", fugitive.position, "after", moves, "moves")
			return 1
		}
		state(chasers, fugitive)
	}

	return 0
}

func state(chasers []*Mover, fugitive *Mover) {
	board := make([][][]*Mover, matrix.height)

	for y, _ := range board {
		board[y] = make([][]*Mover, matrix.width)

		for x, _ := range board[y] {
			board[y][x] = make([]*Mover, 0)
			if fugitive.position.x == x && fugitive.position.y == y {
				board[y][x] = append(board[y][x], fugitive)
			}

			for _, chaser := range chasers {
				if chaser.position.x == x && chaser.position.y == y {
					board[y][x] = append(board[y][x], chaser)
				}
			}
		}
	}

	for _, row := range board {
		for _, col := range row {
			if len(col) == 1 && col[0] == fugitive {
				fmt.Print(" F ")
			} else if len(col) > 1 && col[0] == fugitive {
				fmt.Print(" X ")
			} else {
				fmt.Print(" ", len(col), " ")
			}
		}
		fmt.Print("\n")
	}

	fmt.Print("\n")
}

var matrix Matrix

func parseFlags() (*int, *int, *int) {
	cols := flag.Int("cols", 5, "number of columns")
	rows := flag.Int("rows", 5, "number of rows")
	nChasers := flag.Int("nChasers", 3, "number of chasers")
	flag.Parse()
	return cols, rows, nChasers
}

func main() {
	// use a different seed evety time the program runs
	rand.Seed(time.Now().UTC().UnixNano())

	// get command line arguments
	cols, rows, nChasers := parseFlags()

	// create the matrix
	matrix = Matrix{*cols, *rows}

	// create the chasers
	chasers := make([]*Mover, *nChasers)
	for i, _ := range chasers {
		chasers[i] = newMover("chaser")
	}

	// create the fugitive
	fugitive := newMover("fugitive")

	// start
	pursuit(chasers, fugitive)
}
