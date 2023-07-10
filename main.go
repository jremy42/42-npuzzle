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

func printTimeInfo(elapsed []time.Duration) {

	fmt.Println("time to get next node to explore from Queue :", elapsed[0].String())
	fmt.Println("time applying moves:", elapsed[1].String())
	fmt.Println("time calculating costs and creating node:", elapsed[2].String())
	fmt.Println("time finding if node already exists in nodes list:", elapsed[3].String())
	fmt.Println("time adding node to queues:", elapsed[4].String())
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

func getNextMoves(startPos, goalPos [][]int, scoreFx evalFx, path []byte, currentNode *Item, elapsed []time.Duration, data *safeData, index int, workers int) {
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
		keyNode, queueIndex := matrixToString(nextPos, workers)
		data.muSeen[queueIndex].Lock()
		seenNodesScore, alreadyExplored := data.seenNodes[queueIndex][keyNode]
		data.muSeen[queueIndex].Unlock()
		end = time.Now()
		elapsed[3] += end.Sub(start)
		if !alreadyExplored ||
			score < seenNodesScore {
			start = time.Now()
			item := &Item{node: nextNode}
			data.muQueue[queueIndex].Lock()
			//fmt.Println("Push [1] :index is : ", index, "len of queue :", len(*data.posQueue[index]))
			heap.Push(data.posQueue[queueIndex], item)
			//fmt.Println("Push [1] :index is : ", index, "len of queue :", len(*data.posQueue[index]))
			data.muQueue[queueIndex].Unlock()
			data.muSeen[queueIndex].Lock()
			data.seenNodes[queueIndex][keyNode] = score
			data.muSeen[queueIndex].Unlock()
			end = time.Now()
			elapsed[4] += end.Sub(start)
		}
	}
}


//Trouver un moyen de wait si on est pas sur d'avoir la solution optimale :
// au momen de poper, si on a deja une sol, et que le min des queues est sup a notre sol, on <- end et on quite
func algo(world [][]int, scoreFx evalFx, data *safeData, index int, workers int) {
	goalPos := goal(len(world))
	startPos := Deep2DSliceCopy(world)
	var elapsed [8]time.Duration
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
			//data.muQueue[index].Unlock()
			//data.mu.Unlock()
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
			//data.muQueue[index].Unlock()
			//data.mu.Unlock()
			fmt.Println(index, "End of sim")
			return
		}
		start := time.Now()
		data.muQueue[index].Lock()
		//fmt.Println("Pop : [1] :index is : ", index, "len of queue :", len(*data.posQueue[index]))
		currentNode := (heap.Pop(data.posQueue[index])).(*Item) // Parfois erreur ????
		//fmt.Println("Pop [2] : index is : ", index, "len of queue :", len(*data.posQueue[index]))
		data.muQueue[index].Unlock()
		if foundSol != nil && currentNode.node.score > foundSol.node.score {
				data.mu.Lock()
				data.path = foundSol.node.path
				data.over = true
				data.end <- true
				data.mu.Unlock()
				return
		}
		end := time.Now()
		elapsed[0] += end.Sub(start)

		currentPath := currentNode.node.path
		if tries > 0 && tries%100000 == 0 {
			fmt.Printf("[%d] Time so far : %s | %d * 100k tries. Len of try : %d. Score : %d Len of Queue : %d\n", index, time.Since(startAlgo), tries/100000, len(currentNode.node.path), currentNode.node.score, lenqueue)
		}

		if isEqual(goalPos, currentNode.node.world) {
			printTimeInfo(elapsed[:])
			data.mu.Lock()
			if checkOptimalSolution(currentNode, data) {
				data.path = currentPath
				data.over = true
				data.end <- true
				data.mu.Unlock()
				return
			} else {
				fmt.Println("Found non optimal solution")
				foundSol = currentNode
				data.mu.Unlock()
			}
		}
		getNextMoves(startPos, goalPos, scoreFx, currentPath, currentNode, elapsed[:], data, index, workers)
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

func initData(board [][]int, workers int) (data *safeData) {
	data = &safeData{}
	startPos := Deep2DSliceCopy(board)
	data.seenNodes = make([]map[string]int, workers)
	keyNode, _ := matrixToString(startPos, workers)
	for i := 0; i < workers; i++ {
		data.seenNodes[i] = make(map[string]int, 1)
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
	data.muSeen = make([]sync.Mutex, workers)
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
		data := initData(board, workers)
		for i := 0; i < workers; i++ {
			go algo(board, eval.fx, data, i, workers)
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
