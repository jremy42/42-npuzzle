package input

import (
	"fmt"
	//"go/scanner"
	"os"
	//"strconv"
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



func ParseInput(file *os.File) (size int) {
	// ouvrir le fichier et lire l'entrée
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		fmt.Println(word)
	}

	// fermer le fichier
	defer file.Close()
	size = 0
	return
}

/* func GetInput(file *os.File) int {
	size := parseInput(file)
	for i := 0; i < size; i++ {
		board.Board[i] = make([]int, size)
		for j := 0; j < size; j++ {
			fmt.Fscanf(file, "%d", &board.Board[i][j])
		}
	}
	return size
} */
