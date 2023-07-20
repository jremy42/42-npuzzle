package main

func getValuePostion(board [][]uint8, toFind uint8) Pos2D {
	for i, line := range board {
		for j, value := range line {
			if value == toFind {
				return Pos2D{X: j, Y: i}
			}
		}
	}
	return Pos2D{-1, -1}
}

func swap[T any](a, b *T) {
	*a, *b = *b, *a
}

func moveUp(board [][] uint8) (ok bool, updatedBoard [][]uint8) {
	empty := getValuePostion(board, 0)
	updatedBoard = Deep2DSliceCopy(board)
	if empty.Y != 0 {
		swap(&updatedBoard[empty.Y][empty.X], &updatedBoard[empty.Y-1][empty.X])
		return true, updatedBoard
	}
	return false, nil
}

func moveDown(board [][]uint8) (ok bool, updatedBoard [][]uint8) {
	empty := getValuePostion(board, 0)
	updatedBoard = Deep2DSliceCopy(board)
	if empty.Y != len(board)-1 {
		swap(&updatedBoard[empty.Y][empty.X], &updatedBoard[empty.Y+1][empty.X])
		return true, updatedBoard
	}
	return false, nil
}

func moveLeft(board [][]uint8) (ok bool, updatedBoard [][]uint8) {
	empty := getValuePostion(board, 0)
	updatedBoard = Deep2DSliceCopy(board)
	if empty.X != 0 {
		swap(&updatedBoard[empty.Y][empty.X], &updatedBoard[empty.Y][empty.X-1])
		return true, updatedBoard
	}
	return false, nil
}

func moveRight(board [][]uint8) (ok bool, updatedBoard [][]uint8) {
	empty := getValuePostion(board, 0)
	updatedBoard = Deep2DSliceCopy(board)
	if empty.X != len(board)-1 {
		swap(&updatedBoard[empty.Y][empty.X], &updatedBoard[empty.Y][empty.X+1])
		return true, updatedBoard
	}
	return false, nil
}
