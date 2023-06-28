package main

import (
	"flag"
	"fmt"
)

var (
	file      string
	mapSize   int
	heuristic string
)

func main() {

	flag.StringVar(&file, "f", "", "usage : -f [filename]")
	flag.IntVar(&mapSize, "s", 0, "usage : -s [size]")
	flag.StringVar(&heuristic, "h", "m", "usage : -h m for manhattan or e for euclidean")
	flag.Parse()
	var board [][]int
	var size int

	if file != "" {
		file := OpenFile(file)
		size, board = ParseInput(file)
		fmt.Println("size :", size, "board:", board)
	} else if mapSize > 0 {
		size, board = Generator(mapSize)
		fmt.Println("size :", size, "board:", board)
	}
	for PrintBoard(board) {
		mapSize = 3
		size, board = Generator(mapSize)
	}
	//input.GetInput(file)
}
