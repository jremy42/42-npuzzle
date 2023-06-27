package main

import (
	"flag"
	"fmt"
	"github.com/jremy42/42-npuzzle/input"
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

	if file == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	file := input.OpenFile(file)
	input.ParseInput(file)
	fmt.Println("print file", file)
	//input.GetInput(file)
}
