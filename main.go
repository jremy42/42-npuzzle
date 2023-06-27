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

	if file != "" {
		file := OpenFile(file)
		size, board := ParseInput(file)
		fmt.Println("size :", size, "board:", board)
	} else if mapSize > 0 {
		size, board := Generator(mapSize)
		fmt.Println("size :", size, "board:", board)
	}

	//input.GetInput(file)
}
