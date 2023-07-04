package main

func getEmptySpot(board [][]int) pos2D {
	for i, line := range board {
		for j, value := range line {
			if value == 0 {
				return pos2D{X: j, Y: i}
			}
		}
	}
	return pos2D{-1, -1}
}

func swap[T any](a, b *T) {
	*a, *b = *b, *a
}

func moveUp(board [][]int, empty pos2D) {
	if empty.Y != 0 {
		swap(&board[empty.Y][empty.X], &board[empty.Y-1][empty.X])
	}
}

func moveDown(board [][]int, empty pos2D) {
	if empty.Y != len(board) -1 {
		swap(&board[empty.Y][empty.X], &board[empty.Y+1][empty.X])
	}
}

func moveLeft(board [][]int, empty pos2D){
	if empty.X != 0 {
		swap(&board[empty.Y][empty.X], &board[empty.Y][empty.X - 1])
	}
}

func moveRight(board [][]int, empty pos2D){
	if empty.X != len(board) -1 {
		swap(&board[empty.Y][empty.X], &board[empty.Y][empty.X + 1])
	}
}
