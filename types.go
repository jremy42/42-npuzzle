package main

type Move2D struct {
	dir byte
	X   int
	Y   int
}

type Pos2D struct {
	X int
	Y int
}

type Node struct {
	//world [][]uint8
	path  []byte
	score int
}

type moveFx func([][]uint8) (bool, [][]uint8)

type evalFx func(pos, startPos, goalPos [][]uint8, path []byte) int

type eval struct {
	name string
	fx   evalFx
}
