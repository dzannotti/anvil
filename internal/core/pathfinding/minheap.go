package pathfinding

import "container/heap"

type nodeHeap []*node

func (h nodeHeap) Len() int           { return len(h) }
func (h nodeHeap) Less(i, j int) bool { return h[i].fScore < h[j].fScore }
func (h nodeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *nodeHeap) Push(x any) {
	*h = append(*h, x.(*node))
}

func (h *nodeHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

type minHeap struct {
	heap *nodeHeap
}

func newMinHeap() *minHeap {
	h := make(nodeHeap, 0, 64)
	heap.Init(&h)
	return &minHeap{heap: &h}
}

func (h *minHeap) Push(n *node) {
	heap.Push(h.heap, n)
}

func (h *minHeap) Pop() *node {
	if h.Empty() {
		return nil
	}
	return heap.Pop(h.heap).(*node)
}

func (h *minHeap) Empty() bool {
	return h.heap.Len() == 0
}