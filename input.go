package main

import (
	"fmt"
	//"go/scanner"
	"bufio"
	"os"
	"strconv"
	"unicode/utf8"
)

type Board struct {
	Board [][]int
}

// OpenFile opens a file and returns a pointer to it

func OpenFile(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file")
		os.Exit(1)
	}
	return file
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func scanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])

		// Skip comments.
		if r == '#' {
			for i := start; i < len(data); i++ {
				if data[i] == '\n' {
					start = i
					break
				}
			}
		}
		if !isSpace(r) {
			break
		}
	}
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if r == '#' {
			return i, data[start:i], nil
		}
		if isSpace(r) {
			return i + width, data[start:i], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}

func ParseInput(file *os.File) (size int, board [][]int) {
	scanner := bufio.NewScanner(file)
	scanner.Split(scanWords)
	inputArray := make([]int, 0, 100)

	for scanner.Scan() {
		word := scanner.Text()
		if word == "" {
			continue
		}
		num, err := strconv.Atoi(word)
		if err != nil {
			fmt.Println("Error parsing input")
			os.Exit(1)
		}
		inputArray = append(inputArray, num)
	}
	size = inputArray[0]
	if size*size != len(inputArray)-1 {
		fmt.Println("Error parsing input")
		os.Exit(1)
	}

	board = make([][]int, size)
	for i := 0; i < size; i++ {
		board[i] = make([]int, size)
		for j := 0; j < size; j++ {
			if inputArray[i*size+j+1] < 0 || inputArray[i*size+j+1] > size*size-1 {
				fmt.Println("Error parsing input")
				os.Exit(1)
			}
			board[i][j] = inputArray[i*size+j+1]
		}
	}
	file.Close()
	return
}
