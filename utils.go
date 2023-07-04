package main

import ()

func isEqual[T comparable](a, b [][]T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func Index[T comparable](slice []T, toFind T) int {
	for i, v := range slice {
		if v == toFind {
			return i
		}
	}
	return -1
}

func posAlreadySeen(slice []Node, toFind [][]int) int {
	for index, v := range slice {
		if isEqual(v.world, toFind) {
			return index
		}
	}
	return -1
}

func DeepSliceCopyAndAdd[T any](slice []T, elems ...T) []T {
	newSlice := make([]T, len(slice), len(slice)+len(elems))
	if len(slice) > 0 {
		copy(newSlice, slice)
	}
	newSlice = append(newSlice, elems...)
	return newSlice
}

func Deep2DSliceCopy[T any](slice [][]T) [][]T {
	newSlice := make([][]T, len(slice))
	for i, row := range slice {
		newSlice[i] = make([]T, len(row))
		for j, value := range row {
			newSlice[i][j] = value
		}
	}
	return newSlice
}

func Max(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}
