package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	//"time"
)

type moveFx func([][]int) (bool, [][]int)

var directions = map[byte]moveFx{
	'U': moveUp,
	'D': moveDown,
	'L': moveLeft,
	'R': moveRight,
}

func getNextNodeIndex(queue []Node) int {
	retScore := queue[0].score
	ret := 0
	for index, value := range queue {
		if value.score < retScore {
			retScore = value.score
			ret = index
		}
	}
	return ret
}

func getNextMoves(startPos, goalPos [][]int, scoreFx func(pos, startPos, goalPos [][]int) int, path []byte, currentNode Node, seenNodes []Node) (nextPaths [][]byte, nextNodes []Node, nextSeen []Node) {
	for key, fx := range directions {
		ok, nextPos := fx(currentNode.world)
		if !ok {
			continue
		}
		score := scoreFx(nextPos, startPos, goalPos)
		nextNode := Node{nextPos, score}
		if posAlreadySeen(seenNodes, nextPos) == -1 ||
			score < seenNodes[posAlreadySeen(seenNodes, nextPos)].score {
			toAdd := DeepSliceCopyAndAdd(path, key)
			//fmt.Println("to Add", toAdd, "path", path)
			nextPaths = append(nextPaths, toAdd)
			nextNodes = append(nextNodes, nextNode)
			nextSeen = append(nextSeen, nextNode)
		}
	}
	return
}

func algo(world [][]int, scoreFx func(pos, startPos, goalPos [][]int) int) (currentPath []byte, seenNodes []Node, tries int, maxSizeQueue int) {
	goalPos := goal(len(world))
	startPos := Deep2DSliceCopy(world)
	seenNodes = []Node{{startPos, 0}}
	posQueue := DeepSliceCopyAndAdd(seenNodes)
	pathQueue := [][]byte{{}}

	for ; len(posQueue) > 0; tries++ {
		maxSizeQueue = Max(maxSizeQueue, len(posQueue))

		nextIndex := getNextNodeIndex(posQueue)
		currentNode := posQueue[nextIndex]
		posQueue = append(posQueue[:nextIndex], posQueue[nextIndex+1:]...)

		currentPath = pathQueue[nextIndex]
		pathQueue = append(pathQueue[:nextIndex], pathQueue[nextIndex+1:]...)

		//fmt.Printf("new try %d\n. QueueSize %d\n", tries, len(posQueue))
		//fmt.Printf("score %d => current try %v\n", currentNode.score, currentNode.world)
		//fmt.Printf("current path %v\n", currentPath)
		time.Sleep(0 * time.Millisecond)
		if isEqual(goalPos, currentNode.world) {
			return
		}
		nextPaths, nextPoses, nextSeen := getNextMoves(startPos, goalPos, scoreFx, currentPath, currentNode, seenNodes)
		//fmt.Println("next Paths", nextPaths)
		posQueue = append(posQueue, nextPoses...)
		pathQueue = append(pathQueue, nextPaths...)
		seenNodes = append(seenNodes, nextSeen...)
	}
	return nil, seenNodes, tries, maxSizeQueue
}

func FIFO() func(pos, startPos, goalPos [][]int) int {
	count := 0
	return func(pos, startPos, goalPos [][]int) int {
		count++
		return count
	}
}

func hamming_distance(pos, startPos, goalPos [][]int) int {
	score := 0
	for i, row := range goalPos {
		for j, value := range row {
			if pos[i][j] != value {
				score++
			}
		}
	}
	return score
}

func manhattan_distance(pos, startPos, goalPos [][]int) int {
	score := 0
	for i, row := range goalPos {
		for j, value := range row {
			if pos[i][j] != value {
				score++
			}
		}
	}
	return score
}

func main() {
	var (
		file      string
		mapSize   int
		heuristic string
	)

	flag.StringVar(&file, "f", "", "usage : -f [filename]")
	flag.IntVar(&mapSize, "s", 3, "usage : -s [size]")
	flag.StringVar(&heuristic, "h", "m", "usage : -h m for manhattan or e for euclidean")
	flag.Parse()
	var board [][]int

	if file != "" {
		file := OpenFile(file)
		_, board = ParseInput(file)
	} else if mapSize > 0 {
		board = gridGenerator(mapSize)
	} else {
		fmt.Println("Invalid Map size")
		os.Exit(1)
	}

	/*
	board = [][]int{
		{0, 2, 4},
		{5, 8, 1},
		{3, 6, 7},
	}
	*/

	/*
	ok1, falseBoard1 := moveRight(goal(len(board)))
	ok2, falseBoard2 := moveUp(falseBoard1)
	ok3, falseBoard3 := moveLeft(falseBoard2)
	if !ok1 || !ok2 || !ok3 {
		fmt.Println("Init failure")
		os.Exit(1)
	}
	path, seenPos, tries, sizeMax := algo(falseBoard3, FIFO())
	*/

	path, seenPos, tries, sizeMax := algo(board, FIFO())
	//path, seenPos, tries, sizeMax := algo(board, hamming_distance)

	if path != nil {
		fmt.Println("Succes !")
		fmt.Printf("len of solution %v, %d pos seen, %d tries, %d space complexity\n", len(path), len(seenPos), tries, sizeMax)
		fmt.Println(string(path))
		fmt.Println(board)
	} else {
		fmt.Println("No solution !")
	}
	displayBoard(board)
	/*
		for playBoard(board) {
			mapSize = 3
			board = gridGenerator(mapSize)
		}
	*/
}
