package main

import "sync"

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
	world [][]int
	path  []byte
	score int
}

type moveFx func([][]int) (bool, [][]int)

type evalFx func(pos, startPos, goalPos [][]int, path []byte) int

type eval struct {
	name string
	fx   evalFx
}

type safeData struct {
	mu sync.Mutex

	muQueue  []sync.Mutex
	posQueue []*PriorityQueue

	muSeen    []sync.Mutex
	seenNodes []map[string]int

	tries        int
	maxSizeQueue []int

	path                []byte
	over                bool
	win                 bool
	winScore            int
	idle                int
	closedSetComplexity int
}

type idaData struct {
	fx                  evalFx
	maxScore            int
	path                []byte
	states              [][][]int
	hashes              []string
	goal                [][]int
	closedSetComplexity int
	tries               int
	ramFailure          bool
}

type Result struct {
	path                []byte
	closedSetComplexity int
	tries               int
	ramFailure          bool
}
