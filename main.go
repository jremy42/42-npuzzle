package main

import (
	"container/heap"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

var evals = []eval{
	//{"dijkstra", dijkstra},
	//	{"greedy_hamming", greedy_hamming},
	//	{"greedy_inv", greedy_inv},
	//{"greedy_manhattan", greedy_manhattan},
	{"astar_manhattan", astar_manhattan_generator(1)},
	//{"astar_manhattan", astar_manhattan_generator(2)},
	//	{"astar_hamming", astar_hamming},
	//{"astar_inversion", astar_inv},
}

type safeData struct {
	mu sync.Mutex

	muQueue  []sync.Mutex
	posQueue []*PriorityQueue

	muSeen    []sync.Mutex
	seenNodes []map[string]int

	tries        int
	maxSizeQueue int

	path []byte
	over bool
	end  chan bool
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

func terminateSearch(data *safeData, solutionPath []byte) {
	data.path = solutionPath
	data.over = true
	data.end <- true
}

func getNextMoves(startPos, goalPos [][]int, scoreFx evalFx, path []byte, currentNode *Item, data *safeData, index int, workers int, seenNodeMap int) {
	for _, dir := range directions {
		ok, nextPos := dir.fx(currentNode.node.world)
		if !ok {
			continue
		}
		score := scoreFx(nextPos, startPos, goalPos, path)
		nextPath := DeepSliceCopyAndAdd(path, dir.name)
		nextNode := Node{nextPos, nextPath, score}
		keyNode, queueIndex, seenNodeIndex := matrixToString(nextPos, workers, seenNodeMap)
		data.muSeen[seenNodeIndex].Lock()
		seenNodesScore, alreadyExplored := data.seenNodes[seenNodeIndex][keyNode]
		data.muSeen[seenNodeIndex].Unlock()
		if !alreadyExplored ||
			score < seenNodesScore {
			item := &Item{node: nextNode}
			data.muQueue[queueIndex].Lock()
			heap.Push(data.posQueue[queueIndex], item)
			data.muQueue[queueIndex].Unlock()
			data.muSeen[seenNodeIndex].Lock()
			data.seenNodes[seenNodeIndex][keyNode] = score
			data.muSeen[seenNodeIndex].Unlock()
		}
	}
}

func algo(world [][]int, scoreFx evalFx, data *safeData, index int, workers int, seenNodesMap int) {
	goalPos := goal(len(world))
	startPos := Deep2DSliceCopy(world)
	var foundSol *Item
	startAlgo := time.Now()
	for {
		data.mu.Lock()
		data.tries++
		tries := data.tries
		over := data.over
		data.muQueue[index].Lock()
		lenqueue := len(*data.posQueue[index])
		data.muQueue[index].Unlock()
		data.maxSizeQueue = Max(data.maxSizeQueue, lenqueue)
		data.mu.Unlock()
		if lenqueue == 0 {
			fmt.Println(index, "Empty queue. Waiting")
			time.Sleep(1 * time.Millisecond)
			//Check if all is empty, and exit if so
			data.mu.Lock()
			totalLen := 0
			for _, value := range data.posQueue {
				totalLen += value.Len()
			}
			data.mu.Unlock()
			if totalLen == 0 {
				fmt.Println("all queues are empty. Leaving")
				return
			} else {
				continue
			}
		} else if over {
			fmt.Println(index, "End of sim")
			return
		}
		data.muQueue[index].Lock()
		currentNode := (heap.Pop(data.posQueue[index])).(*Item) // Parfois erreur ????
		data.muQueue[index].Unlock()
		if foundSol != nil && currentNode.node.score > foundSol.node.score {
			data.mu.Lock()
			terminateSearch(data, foundSol.node.path)
			data.mu.Unlock()
			return
		}

		currentPath := currentNode.node.path
		if tries > 0 && tries%100000 == 0 {
			fmt.Printf("[%d] Time so far : %s | %d * 100k tries. Len of try : %d. Score : %d Len of Queue : %d\n", index, time.Since(startAlgo), tries/100000, len(currentNode.node.path), currentNode.node.score, lenqueue)
		}

		if isEqual(goalPos, currentNode.node.world) {
			data.mu.Lock()
			if checkOptimalSolution(currentNode, data) {
				terminateSearch(data, currentNode.node.path)
				data.mu.Unlock()
				return
			} else {
				fmt.Println("Found non optimal solution")
				foundSol = currentNode
				data.mu.Unlock()
			}
		}
		getNextMoves(startPos, goalPos, scoreFx, currentPath, currentNode, data, index, workers, seenNodesMap)
	}
}

func checkOptimalSolution(currentNode *Item, data *safeData) bool {
	bestNodes := make([]*Item, 0, len(data.posQueue))
	for i := range data.posQueue {
		if data.posQueue[i].Len() > 0 {
			bestNodes = append(bestNodes, heap.Pop(data.posQueue[i]).(*Item))
			fmt.Println("best nodes", bestNodes[len(bestNodes)-1])
		} else {
			bestNodes = append(bestNodes, nil)
		}
	}
	fmt.Println("current score :", currentNode.node.score)
	for i := range bestNodes {
		if bestNodes[i] != nil && bestNodes[i].node.score <= currentNode.node.score {
			fmt.Println("current score :", currentNode.node.score, "next score :", bestNodes[i].node.score)
			for j := range bestNodes {
				heap.Push(data.posQueue[j], bestNodes[j])
			}
			return false
		}
	}
	return true
}

func initData(board [][]int, workers int, seenNodesMap int) (data *safeData) {
	data = &safeData{}
	startPos := Deep2DSliceCopy(board)
	data.seenNodes = make([]map[string]int, seenNodesMap)
	keyNode, _, _ := matrixToString(startPos, workers, seenNodesMap)
	for i := 0; i < seenNodesMap; i++ {
		data.seenNodes[i] = make(map[string]int, 1000000)
		data.seenNodes[i][keyNode] = 0
	}
	data.posQueue = make([]*PriorityQueue, workers)
	for i := 0; i < workers; i++ {

		queue := make(PriorityQueue, 1, 1000000)
		queue[0] = &Item{node: Node{world: startPos, score: 0, path: []byte{}}}
		data.posQueue[i] = &queue
		heap.Init(data.posQueue[i])
	}
	data.end = make(chan bool)
	data.over = false
	data.muQueue = make([]sync.Mutex, workers)
	data.muSeen = make([]sync.Mutex, seenNodesMap)
	return
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
		workers := 8
		seenNodeMap := 32
		data := initData(board, workers, seenNodeMap)
		for i := 0; i < workers; i++ {
			go algo(board, eval.fx, data, i, workers, seenNodeMap)
		}
		<-data.end
		end := time.Now()
		elapsed := end.Sub(start)
		if data.path != nil {
			//displayBoard(board, path, seenPos, eval.name+" in "+elapsed.String(), tries, sizeMax)
			fmt.Println("Succes with :", eval.name, "in ", elapsed.String(), "!")
			fmt.Printf("len of solution %v, %d pos seen, %d tries, %d space complexity\n", len(data.path), len(data.seenNodes), data.tries, data.maxSizeQueue)
			fmt.Println(string(data.path))
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
