package main

import (
	"container/heap"
)


type Item struct {
	node  Node
	index int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].node.score < pq[j].node.score
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[j].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	newItem := x.(*Item)
	newItem.index = n
	*pq = append(*pq, newItem)
}

func (pq *PriorityQueue) Pop() any {
	n := len(*pq)
	item := (*pq)[n-1]
	(*pq) = (*pq)[:n-1]
	return item
}

//not mandatory but improves speed. Need index
func (pq *PriorityQueue) update(item *Item, updatedNode Node) {
	item.node = updatedNode
	heap.Fix(pq, item.index)
}

/*
import (
	"container/heap"
)
type Node struct {
	world [][]int
	score int
}

func main() {
	items := []Node{
		Node{world: [][]int{[]int{1, 2, 3}}, score: 3},
		Node{world: [][]int{[]int{7, 9, 1}}, score: 8},
		Node{world: [][]int{[]int{0, 1, 7}}, score: 4},
	}
	pq := make(PriorityQueue, len(items))

	for i, node := range items {
		pq[i] = &Item{
			node:  node,
			//index: i,
		}
	}
	heap.Init(&pq)

	newItem := &Item{
		node :Node{world: [][]int{[]int{7, 7, 7}}, score  : 5},
	}
	heap.Push(&pq, newItem)
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%v\n", item.node)
	}
}
*/
