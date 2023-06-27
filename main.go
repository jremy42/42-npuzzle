package main

import (
	"flag"
	"github.com/jremy/42-npuzzle/input"
	"os"
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

	if file == "" || mapSize == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	file := input.OpenFile(file)
	board := input.GetInput(file, mapSize)
	board.PrintBoard()
}
