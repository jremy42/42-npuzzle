package main

import (
	"container/heap"
	"fmt"
	"os"
	"sync"
	"time"
)

func initData(board [][]int, workers int, seenNodesSplit int) (data safeData) {
	startPos := board
	data.seenNodes = make([]map[string]int, seenNodesSplit)
	keyNode, _, _ := matrixToStringSelector(startPos, workers, seenNodesSplit)
	for i := 0; i < seenNodesSplit; i++ {
		data.seenNodes[i] = make(map[string]int, 1000)
		data.seenNodes[i][keyNode] = 0
	}
	data.posQueue = make([]*PriorityQueue, workers)
	for i := 0; i < workers; i++ {

		queue := make(PriorityQueue, 1, 1000)
		queue[0] = &Item{node: Node{world: startPos, score: 0, path: []byte{}}}
		data.posQueue[i] = &queue
		heap.Init(data.posQueue[i])
	}
	data.over = false
	data.win = false
	data.muQueue = make([]sync.Mutex, workers)
	data.muSeen = make([]sync.Mutex, seenNodesSplit)
	data.maxSizeQueue = make([]int, workers)
	data.idle = 0
	return
}

func algo(world [][]int, scoreFx evalFx, data *safeData, workerIndex int, workers int, seenNodesSplit int) {
	goalPos := goal(len(world))
	startPos := world
	var foundSol *Item
	startAlgo := time.Now()
	isIdle := false
	for {
		over, tries, lenqueue, idle := refreshData(data, workerIndex)
		if over {
			fmt.Fprintf(os.Stderr, "[%2d] - Someone ended sim. Leaving now\n", workerIndex)
			return
		}
		if idle >= workers {
			fmt.Fprintf(os.Stderr, "[%2d] - Everyone is idle\n", workerIndex)
			return
		}
		if isIdle && lenqueue > 0 {
			isIdle = false
			data.mu.Lock()
			data.idle--
			data.mu.Unlock()
		}
		if lenqueue == 0 {
			if !isIdle {
				data.mu.Lock()
				data.idle++
				data.mu.Unlock()
				isIdle = true
			}
			continue
		}
		currentNode := getNextNode(data, workerIndex)
		if currentNode == nil {
			continue
		}
		if foundSol != nil && currentNode.node.score > foundSol.node.score {
			data.mu.Lock()
			terminateSearch(data, foundSol.node.path, foundSol.node.score)
			data.mu.Unlock()
			return
		}
		printInfo(workerIndex, tries, currentNode, startAlgo, lenqueue)
		if isEqual(goalPos, currentNode.node.world) {
			data.mu.Lock()
			if checkOptimalSolution(currentNode, data) {
				fmt.Fprintf(os.Stderr, "\x1b[32m[%2d] - Found an OPTIMAL solution\n\x1b[0m", workerIndex)
				terminateSearch(data, currentNode.node.path, currentNode.node.score)
				data.mu.Unlock()
				return
			} else {
				fmt.Fprintf(os.Stderr, "\x1b[33m[%2d] - Found a NON optimal solution\n\x1b[0m", workerIndex)
				foundSol = currentNode
				data.mu.Unlock()
			}
		}
		getNextMoves(startPos, goalPos, scoreFx, currentNode.node.path, currentNode, data, workerIndex, workers, seenNodesSplit)
	}
}

func checkOptimalSolution(currentNode *Item, data *safeData) bool {
	bestNodes := make([]*Item, len(data.posQueue))
	for i := range data.posQueue {
		data.muQueue[i].Lock()
		if data.posQueue[i].Len() > 0 {
			bestNodes[i] = heap.Pop(data.posQueue[i]).(*Item)
		} else {
			bestNodes[i] = nil
		}
		data.muQueue[i].Unlock()
	}
	for i := range bestNodes {
		if bestNodes[i] != nil && bestNodes[i].node.score <= currentNode.node.score {
			for j := range bestNodes {
				data.muQueue[j].Lock()
				if bestNodes[j] != nil {
					heap.Push(data.posQueue[j], bestNodes[j])
				}
				data.muQueue[j].Unlock()
			}
			return false
		}
	}

	return true
}

func terminateSearch(data *safeData, solutionPath []byte, score int) {
	data.path = solutionPath
	data.over = true
	data.win = true
	data.winScore = score
}

func getNextMoves(startPos, goalPos [][]int, scoreFx evalFx, path []byte, currentNode *Item, data *safeData, index int, workers int, seenNodesSplit int) {
	if data.tries%1000 == 0 {
		availableRAM, err := getAvailableRAM()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if availableRAM>>20 < minRAMAvailableMB {
			fmt.Println("Not enough RAM to continue, try with another heuristic")
			os.Exit(0)
		}
	}
	for _, dir := range directions {
		ok, nextPos := dir.fx(currentNode.node.world)
		if !ok {
			continue
		}
		score := scoreFx(nextPos, startPos, goalPos, path)
		nextPath := DeepSliceCopyAndAdd(path, dir.name)
		nextNode := Node{world: nextPos, path: nextPath, score: score}
		keyNode, queueIndex, seenNodeIndex := matrixToStringSelector(nextPos, workers, seenNodesSplit)
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

func noMoreNodesToExplore(data *safeData) bool {
	data.mu.Lock()
	totalLen := 0
	for i := range data.posQueue {
		data.muQueue[i].Lock()
		length := data.posQueue[i].Len()
		totalLen += length
		data.muQueue[i].Unlock()
	}
	data.mu.Unlock()
	if totalLen == 0 {
		fmt.Fprintln(os.Stderr, "all queues are empty. Leaving")
		return true
	} else {
		return false
	}
}

func refreshData(data *safeData, workerIndex int) (over bool, tries, lenqueue int, idle int) {
	data.mu.Lock()

	data.tries++
	tries = data.tries
	over = data.over
	idle = data.idle
	data.mu.Unlock()
	data.muQueue[workerIndex].Lock()
	lenqueue = len(*data.posQueue[workerIndex])
	data.maxSizeQueue[workerIndex] = Max(data.maxSizeQueue[workerIndex], lenqueue)
	data.muQueue[workerIndex].Unlock()
	return
}

func printInfo(workerIndex int, tries int, currentNode *Item, startAlgo time.Time, lenqueue int) {
	if tries > 0 && tries%100000 == 0 {
		fmt.Fprintf(os.Stderr, "[%2d] Time so far : %s | %d * 100k tries. Len of try : %d. Score : %d Len of Queue : %d\n", workerIndex, time.Since(startAlgo), tries/100000, len(currentNode.node.path), currentNode.node.score, lenqueue)
	}
}

func getNextNode(data *safeData, workerIndex int) (currentNode *Item) {
	data.muQueue[workerIndex].Lock()
	if data.posQueue[workerIndex].Len() != 0 {
		currentNode = (heap.Pop(data.posQueue[workerIndex])).(*Item)
	}
	data.muQueue[workerIndex].Unlock()
	return
}
