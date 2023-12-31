package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"strconv"
)

func isEqual[T comparable](a, b [][]T) bool {
	for i := range a {
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func isEqualTable[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
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

func DeepSliceCopyAndAdd[T any](slice []T, elems ...T) []T {
	newSlice := make([]T, len(slice), len(slice)+len(elems))
	copy(newSlice, slice)
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

func Min(a, b int) int {
	if a >= b {
		return b
	} else {
		return a
	}
}

func Abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

func openDir(dir string) []fs.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func matrixToString2(matrix [][]int) string {

	results := ""

	for i := 0; i < len(matrix); i++ {

		for j := 0; j < len(matrix[i]); j++ {
			results += strconv.Itoa(matrix[i][j]) + "."
		}
	}
	return results
}

func matrixToStringSelector(matrix [][]int, worker int, seenNodeMap int) (key string, queueIndex int, seenNodeIndex int) {
	if len(matrix) < 10 {
		return matrixToStringOptimal(matrix, worker, seenNodeMap)
	} else {
		return matrixToStringNoOpti(matrix, worker, seenNodeMap)
	}
}

func matrixToStringOptimal(matrix [][]int, worker int, seenNodeMap int) (key string, queueIndex int, seenNodeIndex int) {

	results := make([]byte, len(matrix)*len(matrix)*4)
	size := len(matrix)

	spot := 0
	for i := 0; i < size; i++ {

		for j := 0; j < size; j++ {
			queueIndex += matrix[i][j] * i * j
			seenNodeIndex += matrix[i][j] * i * j
			results[spot] = byte(matrix[i][j] / 10)
			results[spot+1] = byte(matrix[i][j] % 10)
			results[spot+2] = '.'
			spot += 3
		}
	}
	queueIndex %= worker
	seenNodeIndex %= seenNodeMap
	return string(results), queueIndex, seenNodeIndex
}

func matrixToStringNoOpti(matrix [][]int, worker int, seenNodeMap int) (key string, queueIndex int, seenNodeIndex int) {

	results := ""
	size := len(matrix)

	for i := 0; i < size; i++ {

		for j := 0; j < size; j++ {
			queueIndex += matrix[i][j] * i * j
			seenNodeIndex += matrix[i][j] * i * j
			results += strconv.Itoa(matrix[i][j]) + "."

		}
	}
	queueIndex %= worker
	seenNodeIndex %= seenNodeMap
	return results, queueIndex, seenNodeIndex
}
