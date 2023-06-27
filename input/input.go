package input

import (
	"fmt"
	//"go/scanner"
	"os"
	//"strconv"
	"bufio"
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
		if !isSpace(r) {
			break
		}
		// Skip comments.
		if r == '#' {
			for i := start; i < len(data); i++ {
				if data[i] == '\n' {
					start = i
					break
				}
			}
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

func ParseInput(file *os.File) (size int) {
	fmt.Println("ParseInput")
	scanner := bufio.NewScanner(file)
	fmt.Print("befeore split")
	scanner.Split(scanWords)
	fmt.Print("after split")
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		word := scanner.Text()
		fmt.Println(word)
	}
	fmt.Println("after scan")
	defer file.Close()
	size = 0
	return
}

/* func GetInput(file *os.File) int {
	size := ParseInput(file)
	for i := 0; i < size; i++ {
		board.Board[i] = make([]int, size)
		for j := 0; j < size; j++ {
			fmt.Fscanf(file, "%d", &board.Board[i][j])
		}
	}
	return size
}
*/
