package data_structure

import (
	"container/heap"
	"fmt"
	"math/rand"
	"sort"
)

type item struct {
	v int
}

func (receiver item) String() string {
	return fmt.Sprintf("%d", receiver.v)
}

func newItem(v int) *item {
	return &item{
		v: v,
	}
}

type items []item

func (h items) update(i int, v int) {
	h[i].v = v
	heap.Fix(&h, i)
}

// Push, Pop implements heap.Interface
func (h *items) Push(x any) {
	*h = append(*h, x.(item))
}
func (h *items) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// implements sort
func (a items) Len() int           { return len(a) }
func (a items) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a items) Less(i, j int) bool { return a[i].v < a[j].v }

func Play_Basic() {
	rand.Seed(0)

	data := make(items, 0)
	for i := 0; i < 100; i++ {
		data = append(data, *newItem(rand.Intn(100)))
	}

	// Sort
	fmt.Println("data: ", data)
	sort.Sort(data)
	fmt.Println("sorted: ", data)

	// Binary search
	binarySearch(data, 74)
	binarySearch(data, 75)

	v, found := binarySearchFind(data, 74)
	fmt.Printf("Find %v -> (%v,%v)\n", 74, v, found)

	v, found = binarySearchFind(data, 75)
	fmt.Printf("Find %v -> (%v,%v)\n", 75, v, found)
}

func Play_Heap() {
	rand.Seed(726251)

	data := make(items, 0)
	for i := 0; i < 10; i++ {
		data = append(data, *newItem(rand.Intn(100)))
	}

	heap.Init(&data)

	fmt.Printf("heap data: %v\n", data)

	doPushToHeap := func(i int) {
		v := *newItem(i)
		heap.Push(&data, v)
		fmt.Printf("push: %v\nheap data: %v\n", i, data)
	}

	doPushToHeap(-1)
	doPushToHeap(200)
	doPushToHeap(99)
	doPushToHeap(21)
	doPushToHeap(70)

	data.update(0, 999999)
	fmt.Printf("update: 0 -> 999999\nheap data: %v\n", data)

	total := data.Len()
	for i := 0; i < total; i++ {
		fmt.Printf("%d -> %v\n", i, heap.Pop(&data))
	}
}

func binarySearch(data items, target int) {
	index := sort.Search(len(data), func(i int) bool {
		return data[i].v >= target
	})
	if index < data.Len() && data[index].v == target {
		fmt.Printf("Binary search for %v -> index: %v, val: %v\n", target, index, data[index])
	} else {
		fmt.Printf("Binary search not found for value: %v \n", target)
	}
}

func binarySearchFind(data items, target int) (*item, bool) {
	index, found := sort.Find(len(data), func(i int) int {
		return target - data[i].v
	})

	if found {
		return &data[index], true
	} else {
		return nil, false
	}
}
