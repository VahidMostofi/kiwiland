package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPQ(t *testing.T) {
	pq := new(PriorityQueue)
	heap.Init(pq)

	heap.Push(pq, NewItem(100, 100))
	heap.Push(pq, NewItem(50, 50))
	heap.Push(pq, NewItem(200, 200))
	assert.Equal(t, 50, heap.Pop(pq).(int))
	assert.Equal(t, 100, heap.Pop(pq).(int))
	heap.Push(pq, NewItem(10, 10))
	heap.Push(pq, NewItem(250, 250))
	assert.Equal(t, 10, heap.Pop(pq).(int))
	assert.Equal(t, 200, heap.Pop(pq).(int))
	heap.Push(pq, NewItem(1, 1))
	assert.Equal(t, 1, heap.Pop(pq).(int))
	assert.Equal(t, 250, heap.Pop(pq).(int))
}
