// PriorityQueue implements the heap.Interface and sort.Interface, so it can
// be used as a priority queue when is passed to the "container/heap" methods.
// src: https://hackernoon.com/today-i-learned-using-priority-queue-in-golang-6f71868902b7
package main

type Item struct {
	container interface{}
	priority  int
	index     int
}

type PriorityQueue []*Item

func NewItem(value interface{}, prio int) *Item {
	return &Item{container: value, priority: prio}
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	itm := old[n-1]
	itm.index = -1
	*pq = old[0 : n-1]
	return itm.container
}
