package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"
)

type moveFx func([][]int) (bool, [][]int)

type evalFx func(pos, startPos, goalPos [][]int, path []byte) int

type eval struct {
	name string
	fx   evalFx
}

var evals = []eval{
	//{"dijkstra", dijkstra},
	//	{"greedy_hamming", greedy_hamming},
	//	{"greedy_inv", greedy_inv},
	{"greedy_manhattan", greedy_manhattan},
	{"astar_custom", astar_custom},
	//{"astar_manhattan_weighted_4", astar_weighted_4},
	//{"astar_manhattan_weighted_2", astar_weighted_2},
	//{"astar_manhattan_weighted_1.5", astar_weighted_1dot5},
	//	{"astar_hamming", astar_hamming},
		{"astar_manhattan", astar_manhattan},
	//	{"astar_inversion", astar_inv},
}

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

func getNextMoves(startPos, goalPos [][]int, scoreFx func(pos, startPos, goalPos [][]int, path []byte) int, path []byte, currentNode Node, seenNodes []Node) (nextPaths [][]byte, nextNodes []Node, nextSeen []Node) {
	for key, fx := range directions {
		ok, nextPos := fx(currentNode.world)
		if !ok {
			continue
		}
		score := scoreFx(nextPos, startPos, goalPos, path)
		nextNode := Node{nextPos, score}
		alreadyExplored := posAlreadySeen(seenNodes, nextPos)
		if alreadyExplored.world == nil ||
			score < alreadyExplored.score {
			toAdd := DeepSliceCopyAndAdd(path, key)
			nextPaths = append(nextPaths, toAdd)
			nextNodes = append(nextNodes, nextNode)
			nextSeen = append(nextSeen, nextNode)
		}
	}
	return
}

func algo(world [][]int, scoreFx func(pos, startPos, goalPos [][]int, path []byte) int) (currentPath []byte, seenNodes []Node, tries int, maxSizeQueue int) {
	goalPos := goal(len(world))
	startPos := Deep2DSliceCopy(world)
	seenNodes = []Node{{startPos, 0}}
	posQueue := DeepSliceCopyAndAdd(seenNodes)
	pathQueue := [][]byte{{}}

	var elapsed time.Duration
	for ; len(posQueue) > 0; tries++ {
		maxSizeQueue = Max(maxSizeQueue, len(posQueue))

		nextIndex := getNextNodeIndex(posQueue)

		currentNode := posQueue[nextIndex]
		currentPath = pathQueue[nextIndex]

		if tries > 0 && tries%1000 == 0 {
			fmt.Println("1k. Len of try :", len(currentPath), "score :", currentNode.score)
		}

		posQueue = append(posQueue[:nextIndex], posQueue[nextIndex+1:]...)
		pathQueue = append(pathQueue[:nextIndex], pathQueue[nextIndex+1:]...)

		//fmt.Printf("new try %d with score %d\n", tries, currentNode.score)
		//fmt.Printf("new try %d\n. QueueSize %d\n", tries, len(posQueue))
		//fmt.Printf("score %d => current try %v\n", currentNode.score, currentNode.world)
		//fmt.Printf("current path %v\n", currentPath)
		//time.Sleep(0 * time.Millisecond)
		if isEqual(goalPos, currentNode.world) {
			return
		}
		start := time.Now()

		nextPaths, nextPoses, nextSeen := getNextMoves(startPos, goalPos, scoreFx, currentPath, currentNode, seenNodes)
		end := time.Now()
		elapsed += end.Sub(start)
		//fmt.Println("found in :", elapsed.String())

		//fmt.Println("next Paths", nextPaths)
		posQueue = append(posQueue, nextPoses...)
		pathQueue = append(pathQueue, nextPaths...)
		seenNodes = append(seenNodes, nextSeen...)
	}
	return nil, seenNodes, tries, maxSizeQueue
}

func dijkstra(pos, startPos, goalPos [][]int, path []byte) int {
	score := len(path) + 1
	return score
}

func astar_hamming(pos, startPos, goalPos [][]int, path []byte) int {
	score := len(path) + 1
	for j, row := range goalPos {
		for i, value := range row {
			if pos[j][i] != value {
				score++
			}
		}
	}
	return score
}

func astar_manhattan(pos, startPos, goalPos [][]int, path []byte) int {
	score := len(path) + 1
	//j is Y
	for j, row := range goalPos {
		//i is X
		for i, value := range row {
			if pos[j][i] != value {
				wrongPositon := getValuePostion(pos, value)
				score += int(math.Abs(float64(wrongPositon.X-i)) + math.Abs(float64(wrongPositon.Y-j)))
			}
		}
	}
	return score
}

func astar_weighted_4(pos, startPos, goalPos [][]int, path []byte) int {
	initDist := len(path) + 1
	//j is Y

	heuristic := 0
	for j, row := range goalPos {
		//i is X
		for i, value := range row {
			if pos[j][i] != value {
				wrongPositon := getValuePostion(pos, value)
				heuristic += int(math.Abs(float64(wrongPositon.X-i)) + math.Abs(float64(wrongPositon.Y-j)))
			}
		}
	}
	return initDist + 4*heuristic
}

func astar_weighted_2(pos, startPos, goalPos [][]int, path []byte) int {
	initDist := len(path) + 1
	//j is Y

	heuristic := 0
	for j, row := range goalPos {
		//i is X
		for i, value := range row {
			if pos[j][i] != value {
				wrongPositon := getValuePostion(pos, value)
				heuristic += int(math.Abs(float64(wrongPositon.X-i)) + math.Abs(float64(wrongPositon.Y-j)))
			}
		}
	}
	return initDist + 2*heuristic
}

func astar_weighted_1dot5(pos, startPos, goalPos [][]int, path []byte) int {
	initDist := len(path) + 1
	//j is Y

	heuristic := 0
	for j, row := range goalPos {
		//i is X
		for i, value := range row {
			if pos[j][i] != value {
				wrongPositon := getValuePostion(pos, value)
				heuristic += int(math.Abs(float64(wrongPositon.X-i)) + math.Abs(float64(wrongPositon.Y-j)))
			}
		}
	}
	return initDist + (heuristic*3)/2
}

func astar_custom(pos, startPos, goalPos [][]int, path []byte) int {
	initDist := len(path) + 1
	//j is Y
	heuristic := 0
	for j, row := range goalPos {
		//i is X
		for i, value := range row {
			if pos[j][i] != value {
				wrongPositon := getValuePostion(pos, value)
				heuristic += int(math.Abs(float64(wrongPositon.X-i)) + math.Abs(float64(wrongPositon.Y-j)))
			}
		}
	}
	return initDist + 10 * heuristic //+ 0 * greedy_inv(pos, startPos, goalPos, path)
}

func greedy_manhattan(pos, startPos, goalPos [][]int, path []byte) int {
	score := 0
	for j, row := range goalPos {
		for i, value := range row {
			if pos[j][i] != value {
				wrongPositon := getValuePostion(pos, value)
				score += int(math.Abs(float64(wrongPositon.X-i)) + math.Abs(float64(wrongPositon.Y-j)))
			}
		}
	}
	return score
}

func greedy_hamming(pos, startPos, goalPos [][]int, path []byte) int {
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

func greedy_inv(pos, startPos, goalPos [][]int, path []byte) int {
	score := 0
	flattenedPos := make([]int, 0, len(pos)*len(pos))
	for _, row := range pos {
		for _, value := range row {
			flattenedPos = append(flattenedPos, value)
		}
	}
	inversion := 0
	for i := range flattenedPos {
		for j := i + 1; j < len(flattenedPos); j++ {
			if flattenedPos[i] > 0 && flattenedPos[j] > 0 && flattenedPos[i] > flattenedPos[j] {
				inversion++
			}
		}
	}
	return score + inversion
}

func astar_inv(pos, startPos, goalPos [][]int, path []byte) int {
	score := len(path) + 1
	flattenedPos := make([]int, 0, len(pos)*len(pos))
	for _, row := range pos {
		for _, value := range row {
			flattenedPos = append(flattenedPos, value)
		}
	}
	inversion := 0
	for i := range flattenedPos {
		for j := i + 1; j < len(flattenedPos); j++ {
			if flattenedPos[i] > 0 && flattenedPos[j] > 0 && flattenedPos[i] > flattenedPos[j] {
				inversion++
			}
		}
	}
	return score + inversion
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
			{2, 4, 5},
			{6, 7, 8},
			{1, 3, 0},
		}
	*/

	/*
		board = [][]int{
			{2, 1, 7},
			{4, 0, 8},
			{3, 5, 6},
		}
	*/

	/*
		board = [][]int{
			{10, 3, 14, 4},
			{2, 12, 5, 6},
			{11, 0, 15, 7},
			{1, 9, 8, 13},
		}
			ok1, falseBoard1 := moveRight(goal(len(board)))
			ok2, falseBoard2 := moveUp(falseBoard1)
			ok3, falseBoard3 := moveLeft(falseBoard2)
			ok4, falseBoard4 := moveDown(falseBoard3)
			if !ok1 || !ok2 || !ok3  || !ok4 {
				fmt.Println("Init failure")
				os.Exit(1)
			}
			board = falseBoard4
	*/

		}

		board = [][]int{
			{2, 3, 4, 6},
			{1, 12, 14, 15},
			{11, 9, 13, 5},
			{10, 8, 0, 7},
		}
	*/
	/*
		board = [][]int{
			{14, 4, 0, 12, 1},
			{24, 11, 15, 10, 5},
			{17, 2, 19, 23, 21},
			{9, 3, 20, 8, 6},
			{7, 18, 16, 22, 13},
		}
	*/

	/*
		ok1, falseBoard1 := moveRight(goal(len(board)))
		ok2, falseBoard2 := moveUp(falseBoard1)
		ok3, falseBoard3 := moveLeft(falseBoard2)
		ok4, falseBoard4 := moveDown(falseBoard3)
		if !ok1 || !ok2 || !ok3  || !ok4 {
			fmt.Println("Init failure")
			os.Exit(1)
		}
		board = falseBoard4
	*/
	if !isSolvable(board) {
		fmt.Println("Board is not solvable")
		os.Exit(1)
	}
	fmt.Println("Board is :", board)
	for _, eval := range evals {
		fmt.Println("Now starting with :", eval.name)
		start := time.Now()
		path, seenPos, tries, sizeMax := algo(board, eval.fx)
		end := time.Now()
		elapsed := end.Sub(start)
		if path != nil {
			displayBoard(board, path, seenPos, tries, sizeMax)
			fmt.Println("Succes with :", eval.name, "in ", elapsed.String(), "!")
			fmt.Printf("len of solution %v, %d pos seen, %d tries, %d space complexity\n", len(path), len(seenPos), tries, sizeMax)
			fmt.Println(string(path))
		} else {
			fmt.Println("No solution !")
		}
	}
	/*
		for playBoard(board) {
			mapSize = 3
			board = gridGenerator(mapSize)
		}
	*/
}
