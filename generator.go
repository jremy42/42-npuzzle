package main

import ()

func gridGenerator(mapSize int) (board [][]uint8) {

	for {

		randomNumber := make(map[int]int)
		for i := 0; i < mapSize*mapSize; i++ {
			randomNumber[i] = i
		}
		board = make([][]uint8, mapSize)

		j := 0
		i := 0
		for _, number := range randomNumber {
			if i%mapSize == 0 {
				board[j] = make([]uint8, mapSize)
				j++
			}
			board[j-1][i%mapSize] = uint8(number)
			i++
		}
		if isSolvable(board) {
			break
		}
	}
	return board
}

func goal(mapSize int) (goal [][]uint8) {

	goal = make([][]uint8, mapSize)
	for i := range goal {
		goal[i] = make([]uint8, mapSize)
	}
	states := []Move2D{
		{'r', 1, 0},
		{'d', 0, 1},
		{'l', -1, 0},
		{'u', 0, -1},
	}
	goal[0][0] = 1
	for i, j, dir, count := 0, 0, 0, 1; count < (mapSize*mapSize)-1; {
		currMove := states[dir%4]
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
			goal[j][i] = uint8(count)
		}
	}
	return goal
}
