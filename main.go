package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Matrix struct {
	Width  int
	Height int
}

func NewMatrix(width int, height int) *Matrix {
	matrix := new(Matrix)
	matrix.Width = width
	matrix.Height = height
	return matrix
}

type Mover struct {
	X int
	Y int
}

func NewMover() *Mover {
	mover := new(Mover)
	mover.X = rand.Intn(matrix.Width)
	mover.Y = rand.Intn(matrix.Height)
	return mover
}

func moveRandom(mover *Mover) {
	p := rand.Intn(2) // move on X or Y?
	x := rand.Intn(2)
	y := rand.Intn(2)
	if p == 1 {
		if (x == 1 || mover.X == 0) && mover.X < matrix.Width {
			mover.X++
		} else {
			mover.X--
		}
	} else {
		if (y == 1 || mover.Y == 0) && mover.Y < matrix.Height {
			mover.Y++
		} else {
			mover.Y--
		}
	}
}

func overlapping(mover1 *Mover, mover2 *Mover) bool {
	return mover1.X == mover2.X && mover1.Y == mover2.Y
}

func pursue(chaser *Mover, fugitive *Mover) {
	if chaser.X == fugitive.X {
		pursueInLine(&chaser.Y, &fugitive.Y)
	} else {
		pursueInLine(&chaser.X, &fugitive.X)
	}
}

// moves chaser, in X or Y, to reduce its distance from fugitive
func pursueInLine(chaserCoord *int, fugitiveCoord *int) {
	if *chaserCoord > *fugitiveCoord {
		*chaserCoord--
	} else {
		*chaserCoord++
	}
}

// moves chaser, in X or Y, to increase its distance from fugitive
func escapeInLine(chaserCoord *int, fugitiveCoord *int) {
	if *chaserCoord > *fugitiveCoord {
		*chaserCoord++
	} else {
		*chaserCoord--
	}
}

// moves fugitive away from chaser, in X or Y
func escape(fugitive *Mover, chaser *Mover) {
	if chaser.X == fugitive.X {
		escapeInLine(&chaser.Y, &fugitive.Y)
	} else {
		escapeInLine(&chaser.X, &fugitive.X)
	}
}

func moverSeesMover(mover1 *Mover, mover2 *Mover) bool {
	return mover1.X == mover2.X || mover1.Y == mover2.Y
}

func fugitiveMove(fugitive *Mover, chasers []*Mover) {
	moved := false
	for _, chaser := range chasers {
		if moverSeesMover(fugitive, chaser) {
			escape(fugitive, chaser)
			moved = true
			break
		}
	}
	if !moved {
		moveRandom(fugitive)
	}
}

func chaserMoves(chasers []*Mover, fugitive *Mover, moves int) int {
	for _, chaser := range chasers {
		if overlapping(chaser, fugitive) {
			fmt.Println("Caught at", chaser.X, chaser.Y, "after", moves, "moves")
			return 1
		} else if moverSeesMover(chaser, fugitive) {
			pursue(chaser, fugitive)
		} else {
			moveRandom(chaser)
		}
	}

	return 0
}

func pursuit(chasers []*Mover, fugitive *Mover) int {
	for moves := 1; moves < 1000000; moves++ {
		fugitiveMove(fugitive, chasers)
		if chaserMoves(chasers, fugitive, moves) == 1 {
			// fugitive was caught
			return 1
		}
	}

	return 0
}

var matrix *Matrix

func main() {
	// use a different seed evety time the program runs
	rand.Seed(time.Now().UTC().UnixNano())

	// create the matrix
	matrix = NewMatrix(5, 5)

	// create the chasers
	nChasers := 4
	chasers := make([]*Mover, nChasers)
	for i, _ := range chasers {
		chasers[i] = NewMover()
	}

	// create the fugitive
	fugitive := NewMover()

	// start
	pursuit(chasers, fugitive)
}
