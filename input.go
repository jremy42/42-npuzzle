package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Board struct {
	Board [][]int
}

func OpenFile(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file")
		os.Exit(1)
	}
	return file
}

func ParseInput(file *os.File) (size int, board [][]int) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	inputArray := readInputArray(scanner)
	size = extractSize(inputArray)
	board = createBoard(size, inputArray)
	file.Close()
	return size, board
}

func alreadyInArray(array []int, num int) bool {
	if len(array) == 0 {
		return false
	}
	for i := 1; i < len(array); i++ {
		if array[i] == num {
			return true
		}
	}
	return false
}

func readInputArray(scanner *bufio.Scanner) []int {
	inputArray := make([]int, 0, 100)

	for scanner.Scan() {
		line := scanner.Text()
		comment := strings.Index(line, "#")
		if comment != -1 {
			line = line[:comment]
		}
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		words := strings.Fields(line)
		for _, word := range words {
			num, err := strconv.Atoi(word)
			if err != nil || num < 0 || alreadyInArray(inputArray, num) {
				fmt.Println("Error parsing input 1")
				os.Exit(1)
			}
			inputArray = append(inputArray, num)
		}
	}

	return inputArray
}

func extractSize(inputArray []int) int {
	if len(inputArray) == 0 {
		fmt.Println("Error parsing input 2")
		os.Exit(1)
	}
	size := inputArray[0]
	if size*size != len(inputArray)-1 {
		fmt.Println("Error parsing input 3")
		os.Exit(1)
	}
	return size
}

func createBoard(size int, inputArray []int) [][]int {
	board := make([][]int, size)
	for i := 0; i < size; i++ {
		board[i] = make([]int, size)
		for j := 0; j < size; j++ {
			if inputArray[i*size+j+1] < 0 || inputArray[i*size+j+1] > size*size-1 {
				fmt.Println("Error parsing input 4", size, inputArray[i*size+j+1], i, j)
				os.Exit(1)
			}
			board[i][j] = inputArray[i*size+j+1]
		}
	}
	return board
}
