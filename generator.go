package main

import (
	"fmt"
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
