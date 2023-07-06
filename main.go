package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var evals = []eval{
	//{"dijkstra", dijkstra},
	//	{"greedy_hamming", greedy_hamming},
	//	{"greedy_inv", greedy_inv},
	//{"greedy_manhattan", greedy_manhattan},
	//	{"astar_hamming", astar_hamming},
	{"astar_manhattan", astar_manhattan_generator(1)},
	//{"astar_inversion", astar_inv},
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

func getNextMoves(startPos, goalPos [][]int, scoreFx func(pos, startPos, goalPos [][]int, path []byte) int, path []byte, currentNode Node, seenNodes *map[string]int, elapsed *[4]time.Duration) (nextPaths [][]byte, nextNodes []Node) {
	for key, fx := range directions {
		start := time.Now()
		ok, nextPos := fx(currentNode.world)
		end := time.Now()
		elapsed[1] += end.Sub(start)
		if !ok {
			continue
		}
		score := scoreFx(nextPos, startPos, goalPos, path)
		nextNode := Node{nextPos, score}
		start = time.Now()
		keyNode := matrixToString(nextPos)
		seenNodesScore, alreadyExplored := (*seenNodes)[keyNode]
		//alreadyExplored := posAlreadySeen(seenNodes, nextPos)
		end = time.Now()
		elapsed[2] += end.Sub(start)
		if alreadyExplored == false ||
			score < seenNodesScore {
			start = time.Now()
			toAdd := DeepSliceCopyAndAdd(path, key)
			nextPaths = append(nextPaths, toAdd)
			nextNodes = append(nextNodes, nextNode)
			//nextSeen = append(nextSeen, nextNode)
			(*seenNodes)[keyNode] = score
			end = time.Now()
			elapsed[3] += end.Sub(start)
		}
	}
	return
}

func algo(world [][]int, scoreFx func(pos, startPos, goalPos [][]int, path []byte) int) (currentPath []byte, seenNodes map[string]int, tries int, maxSizeQueue int) {
	goalPos := goal(len(world))
	startPos := Deep2DSliceCopy(world)
	//seenNodes = []Node{{startPos, 0}}
	seenNodes = make(map[string]int, 10000)
	seenNodes[matrixToString(startPos)] = 0
	posQueue := []Node{{startPos, 0}}
	pathQueue := [][]byte{{}}

	var elapsed [4]time.Duration
	for ; len(posQueue) > 0; tries++ {
		maxSizeQueue = Max(maxSizeQueue, len(posQueue))

		nextIndex := getNextNodeIndex(posQueue)

		currentNode := posQueue[nextIndex]
		currentPath = pathQueue[nextIndex]

		if tries > 0 && tries%1000 == 0 {
			fmt.Printf("%d k tries. Len of try : %d. Score : %d\n", tries/1000, len(currentPath), currentNode.score)
		}

		posQueue = append(posQueue[:nextIndex], posQueue[nextIndex+1:]...)
		pathQueue = append(pathQueue[:nextIndex], pathQueue[nextIndex+1:]...)

		//fmt.Printf("new try %d with score %d\n", tries, currentNode.score)
		//fmt.Printf("new try %d\n. QueueSize %d\n", tries, len(posQueue))
		//fmt.Printf("score %d => current try %v\n", currentNode.score, currentNode.world)
		//fmt.Printf("current path %v\n", currentPath)
		//time.Sleep(0 * time.Millisecond)
		if isEqual(goalPos, currentNode.world) {
			fmt.Println("getNexMoves total time :", elapsed[0].String())
			fmt.Println("time applying moves:", elapsed[1].String())
			fmt.Println("time finding if node already exists in nodes list:", elapsed[2].String())
			fmt.Println("time adding node to queues:", elapsed[3].String())
			return
		}

		start := time.Now()
		nextPaths, nextPoses := getNextMoves(startPos, goalPos, scoreFx, currentPath, currentNode, &seenNodes, &elapsed)
		end := time.Now()
		elapsed[0] += end.Sub(start)

		//fmt.Println("next Paths", nextPaths)
		posQueue = append(posQueue, nextPoses...)
		pathQueue = append(pathQueue, nextPaths...)
		//seenNodes = append(seenNodes, nextSeen...)
	}
	return nil, seenNodes, tries, maxSizeQueue
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

		/* board = [][]int{
				{2, 1, 7},
				{4, 0, 8},
				{3, 5, 6},
			} */

	board = [][]int{
		{10, 3, 14, 4},
		{2, 12, 5, 6},
		{11, 0, 15, 7},
		{1, 9, 8, 13},
	}
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

	/*
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
			//displayBoard(board, path, seenPos, eval.name+" in "+elapsed.String(), tries, sizeMax)
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
