package main

import (
	"fmt"
)

func IsSolvable(board [][]int) bool {
	fmt.Println("board :", board)
	return true
}

func isSolved(board [][]int) bool {
	if board[len(board)/2][len(board)/2] != 0 {
		return false
	}
	for i := 0; i < len(board)-1; i++ {
		for j := 0; j < len(board)-1; j++ {
			if board[i][j] != i*len(board)+j+1 {
				return false
			}
		}
	}
	return true
}
