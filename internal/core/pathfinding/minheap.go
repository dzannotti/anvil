package pathfinding

type minHeap struct {
	nodes []*node
}

func newMinHeap() *minHeap {
	return &minHeap{
		nodes: make([]*node, 0, 64), // Pre-allocate some capacity
	}
}

func (h *minHeap) Len() int {
	return len(h.nodes)
}

func (h *minHeap) Push(n *node) {
	h.nodes = append(h.nodes, n)
	h.upHeapify(len(h.nodes) - 1)
}

func (h *minHeap) Pop() *node {
	if len(h.nodes) == 0 {
		return nil
	}

	root := h.nodes[0]
	lastIdx := len(h.nodes) - 1
	h.nodes[0] = h.nodes[lastIdx]
	h.nodes = h.nodes[:lastIdx]

	if lastIdx > 0 {
		h.downHeapify(0)
	}

	return root
}

func (h *minHeap) upHeapify(idx int) {
	for idx > 0 {
		parentIdx := (idx - 1) / 2
		if h.nodes[parentIdx].fScore <= h.nodes[idx].fScore {
			break
		}
		h.nodes[parentIdx], h.nodes[idx] = h.nodes[idx], h.nodes[parentIdx]
		idx = parentIdx
	}
}

func (h *minHeap) downHeapify(idx int) {
	smallest := idx
	left := 2*idx + 1
	right := 2*idx + 2

	if left < len(h.nodes) && h.nodes[left].fScore < h.nodes[smallest].fScore {
		smallest = left
	}

	if right < len(h.nodes) && h.nodes[right].fScore < h.nodes[smallest].fScore {
		smallest = right
	}

	if smallest != idx {
		h.nodes[idx], h.nodes[smallest] = h.nodes[smallest], h.nodes[idx]
		h.downHeapify(smallest)
	}
}

func (h *minHeap) Empty() bool {
	return len(h.nodes) == 0
}
