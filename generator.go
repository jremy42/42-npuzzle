package main

import (
	"fmt"
	"time"
)

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

func findGoal(mapSize int) (goal [][]int) {

	goal = make([][]int, mapSize)
	states := [4]byte{'r', 'd', 'l', 'u'}

	for i := 0; i < mapSize; i++ {
		goal[i] = make([]int, mapSize)
	}
	countLetter := 1
	iState := 0
	for i, j := 0, 0; countLetter < mapSize*mapSize-1; {
		state := states[iState%4]

		switch state {
		case 'r':
			for j < mapSize-1 && goal[i][j] == 0 {
				goal[i][j] = countLetter
				j++
				countLetter++
				fmt.Println("i :", i, "j :", j, "countLetter :", countLetter, "State :", string(state))
				fmt.Println("goal :", goal)
				time.Sleep(1000 * time.Millisecond)

			}
			if goal[i][j] != 0 {
				j--
			}
		case 'd':
			for i < mapSize-1 && goal[i][j] == 0 {
				goal[i][j] = countLetter
				i++
				countLetter++
				fmt.Println("i :", i, "j :", j, "countLetter :", countLetter, "State :", string(state))
				fmt.Println("goal :", goal)
				time.Sleep(1000 * time.Millisecond)
			}
			if goal[i][j] != 0 {
				i--
			}
		case 'l':
			for j > 0 && goal[i][j] == 0 {
				goal[i][j] = countLetter
				j--
				countLetter++
				fmt.Println("i :", i, "j :", j, "countLetter :", countLetter, "State :", string(state))
				fmt.Println("goal :", goal)
				time.Sleep(1000 * time.Millisecond)
			}
			if goal[i][j] != 0 {
				j++
			}

		case 'u':
			for i > 4/iState && goal[i][j] == 0 {
				goal[i][j] = countLetter
				i--
				countLetter++
				fmt.Println("i :", i, "j :", j, "countLetter :", countLetter, "State :", string(state))
				fmt.Println("goal :", goal)
				time.Sleep(1000 * time.Millisecond)
			}
		}
		iState++
	}
	return goal // TODO
}
