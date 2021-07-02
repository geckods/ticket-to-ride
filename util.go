package main

import (
	"container/heap"
	"reflect"
)

const UintSize = 32 << (^uint(0) >> 32 & 1) // 32 or 64


const (
	MaxInt  = 1<<(UintSize-2) - 1 // 1<<31 - 1 or 1<<63 - 1
	MaxUint = 1<<UintSize - 1     // 1<<32 - 1 or 1<<64 - 1
)


func itemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func normalizeFloatSlice(sl *[]float64){
	sum := float64(0)
	for _,elem := range *sl {
		sum+=elem
	}
	if sum == 0.0 {
		return
	}
	for i,_ := range *sl {
		(*sl)[i]/=sum
	}
}

func scaleFloat(x,a,b float64) float64{
	return a+((b-a)*x)
}

// An Item is something we manage in a priority queue. In our case, we're using it for djikstra, so it's distance and index
type DjikstraItem struct {
	vertexIndex int // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}


// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*DjikstraItem

func (pq PriorityQueue) Len() int { return len(pq) }

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
	item := x.(*DjikstraItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	//old[n-1] = nil  // avoid memory leak
	//item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *DjikstraItem, value int, priority int) {
	item.vertexIndex = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
