package input

import (
	"fmt"
	"go/scanner"
	"os"
	"strconv"
	"bufio"
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



func parseInput(file *os.File) (int, int[]) {
	//open file and read input
	defer file.Close()
    scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		fmt.Println(word)
	}
	return size, parseInput
}

func GetInput(file *os.File) Board {
	board := Board{make([][]int, size)}
	size, parseInput := parseInput(file)
	for i := 0; i < size; i++ {
		board.Board[i] = make([]int, size)
		for j := 0; j < size; j++ {
			fmt.Fscanf(file, "%d", &board.Board[i][j])
		}
	}
	return board
}
