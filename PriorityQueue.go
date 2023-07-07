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

func (pq *PriorityQueue) update(item *Item, updatedNode Node) {
	item.node = updatedNode
	heap.Fix(pq, item.index)
}
