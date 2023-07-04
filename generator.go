package main

import (
	"fmt"
	//"time"
)

type move2D struct {
	dir byte
	X   int
	Y   int
}

func Generator(mapSize int) (size int, board [][]int) {

	randomNumber := make(map[int]int)
	for i := 0; i < mapSize*mapSize; i++ {
		randomNumber[i] = i
	}
	fmt.Println("randomNumber :", randomNumber)

	board = make([][]int, mapSize)

	j := 0
	i := 0
	for _, number := range randomNumber {
		fmt.Println("number :", number)
		if i%mapSize == 0 {
			board[j] = make([]int, mapSize)
			j++
		}
		board[j-1][i%mapSize] = number
		i++
	}

	return mapSize, board
}

//j <=> y
//x <=> i
func findGoal(mapSize int) (goal [][]int) {

	goal = make([][]int, mapSize)
	for i := range goal {
		goal[i] = make([]int, mapSize)
	}
	states := []move2D{
		{'r', 1, 0},
		{'d', 0, 1},
		{'l', -1, 0},
		{'u', 0, -1},
	}
	goal[0][0] = 1
	for i, j, dir, count := 0, 0, 0, 1; count < (mapSize*mapSize)-1; {
		currMove := states[dir % 4]
		nextJ := j + currMove.Y
		nextI := i + currMove.X
		if nextI > mapSize-1 ||
			nextI < 0 ||
			nextJ > mapSize-1 ||
			nextJ < 0 ||
			goal[nextJ][nextI] != 0 {
			dir++
		} else {
			i, j = nextI, nextJ
			count++
			goal[j][i] = count
		}
	}
	return goal
}
