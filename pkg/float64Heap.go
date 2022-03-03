package pkg

import (
	"container/heap"
	"math"
)

type Float64Heap []float64

func (h Float64Heap) Len() int           { return len(h) }
func (h Float64Heap) Less(i, j int) bool { return h[i] > h[j] } // 大根堆
func (h Float64Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *Float64Heap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(float64))
}

func (h *Float64Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// GetTopRankAverage 从5%～30%部分
func (h *Float64Heap) GetTopRankAverage() (average float64) {
	startRank := int(math.Floor(float64(h.Len()) * 0.05))
	endRank := int(math.Ceil(float64(h.Len()) * 0.3))
	if endRank > h.Len()-1 {
		endRank = h.Len() - 1
	}
	if startRank < 0 {
		startRank = 0
	}
	if startRank > endRank {
		return 0
	}
	sum := 0.0
	for i := 0; i < endRank; i++ {
		if i >= startRank {
			sum += heap.Pop(h).(float64)
		} else {
			h.Pop()
		}
	}
	return sum / float64(endRank-startRank+1)
}
