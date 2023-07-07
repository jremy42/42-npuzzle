package main

import (
	"container/heap"
	"flag"
	"fmt"
	//"log"
	"os"
	"time"
)

var evals = []eval{
	//{"dijkstra", dijkstra},
	//	{"greedy_hamming", greedy_hamming},
	//	{"greedy_inv", greedy_inv},
	{"greedy_manhattan", greedy_manhattan},
	{"astar_manhattan", astar_manhattan_generator(1)},
	//{"astar_manhattan", astar_manhattan_generator(2)},
	//	{"astar_hamming", astar_hamming},
	//{"astar_inversion", astar_inv},
}

var directions = []struct {
	name byte
	fx   moveFx
}{
	{'U', moveUp},
	{'D', moveDown},
	{'L', moveLeft},
	{'R', moveRight},
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

func getNextMoves(startPos, goalPos [][]int, scoreFx func(pos, startPos, goalPos [][]int, path []byte) int, path []byte, currentNode *Item, seenNodes *map[string]int, posQueue *PriorityQueue, elapsed *[8]time.Duration) {
	for _, dir := range directions {
		start := time.Now()
		ok, nextPos := dir.fx(currentNode.node.world)
		end := time.Now()
		elapsed[1] += end.Sub(start)
		if !ok {
			continue
		}
		start = time.Now()
		score := scoreFx(nextPos, startPos, goalPos, path)
		nextPath := DeepSliceCopyAndAdd(path, dir.name)
		nextNode := Node{nextPos, nextPath, score}
		end = time.Now()
		elapsed[2] += end.Sub(start)
		start = time.Now()
		keyNode := matrixToString(nextPos)
		seenNodesScore, alreadyExplored := (*seenNodes)[keyNode]
		end = time.Now()
		elapsed[3] += end.Sub(start)
		if !alreadyExplored ||
			score < seenNodesScore {
			start = time.Now()
			item := &Item{node: nextNode}
			heap.Push(posQueue, item)
			(*seenNodes)[keyNode] = score
			end = time.Now()
			elapsed[4] += end.Sub(start)
		}
	}
}

func algo(world [][]int, scoreFx func(pos, startPos, goalPos [][]int, path []byte) int) (currentPath []byte, seenNodes map[string]int, tries int, maxSizeQueue int) {
	goalPos := goal(len(world))
	startPos := Deep2DSliceCopy(world)
	seenNodes = make(map[string]int, 1000000)
	seenNodes[matrixToString(startPos)] = 0
	posQueue := make(PriorityQueue, 1, 1000000)
	posQueue[0] = &Item{node: Node{world: startPos, score: 0, path: []byte{}}}
	heap.Init(&posQueue)

	var elapsed [8]time.Duration
	startAlgo := time.Now()
	for lenqueue := 1; len(posQueue) > 0; lenqueue, tries = len(posQueue), tries+1 {
		maxSizeQueue = Max(maxSizeQueue, lenqueue)

		start := time.Now()
		currentNode := heap.Pop(&posQueue).(*Item)
		end := time.Now()
		elapsed[0] += end.Sub(start)

		currentPath = currentNode.node.path
		if tries > 0 && tries%100000 == 0 {
			fmt.Printf("Time so far : %s | %d * 100k tries. Len of try : %d. Score : %d Len of Queue : %d\n", time.Since(startAlgo), tries/100000, len(currentNode.node.path), currentNode.node.score, lenqueue)
		}

		if isEqual(goalPos, currentNode.node.world) {
			fmt.Println("getNexMoves total time :", elapsed[0].String())
			fmt.Println("time applying moves:", elapsed[1].String())
			fmt.Println("time calculating costs and creating node:", elapsed[2].String())
			fmt.Println("time finding if node already exists in nodes list:", elapsed[3].String())
			fmt.Println("time adding node to queues:", elapsed[4].String())
			return
		}
		getNextMoves(startPos, goalPos, scoreFx, currentPath, currentNode, &seenNodes, &posQueue, &elapsed)
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
