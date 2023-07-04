package main

//func getEmptySpot(board [][]int)

func moveUp(board [][]int) [][]int {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board); j++ {
			if board[i][j] == 0 && i != 0 {
				board[i][j], board[i-1][j] = board[i-1][j], board[i][j]
				return board
			}
		}
	}
	return board
}

func moveDown(board [][]int) [][]int {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board); j++ {
			if board[i][j] == 0 && i != len(board)-1 {
				board[i][j], board[i+1][j] = board[i+1][j], board[i][j]
				return board
			}
		}
	}
	return board
}

func moveLeft(board [][]int) [][]int {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board); j++ {
			if board[i][j] == 0 && j != 0 {
				board[i][j], board[i][j-1] = board[i][j-1], board[i][j]
				return board
			}
		}
	}
	return board
}

func moveRight(board [][]int) [][]int {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board); j++ {
			if board[i][j] == 0 && j != len(board)-1 {
				board[i][j], board[i][j+1] = board[i][j+1], board[i][j]
				return board
			}
		}
	}
	return board
}
