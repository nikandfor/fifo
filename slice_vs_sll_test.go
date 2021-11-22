package fifo

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

// Code from https://gist.github.com/Codebreaker101/cd42f5dbc0e84bf5c85c2ab963e4d551
// Which is from reddit post https://www.reddit.com/r/golang/comments/qz3jvj/memory_efficient_fifo_structure_could_it_be_that/
// Except Benchmarks

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

type QueueSlice struct {
	lock sync.Mutex // you don't have to do this if you don't want thread safety
	s    []int
}

func NewQueueSlice() *QueueSlice {
	return &QueueSlice{sync.Mutex{}, make([]int, 0)}
}

func (s *QueueSlice) Push(v int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *QueueSlice) Pop() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0
	}

	res := s.s[0]
	s.s = s.s[1:]
	return res
}

func (s *QueueSlice) Len() int {
	return len(s.s)
}

func (s *QueueSlice) String() string {
	return fmt.Sprintf("Capacity: %d Length: %d", cap(s.s), len(s.s))
}

func TestQueueSliceLocal(t *testing.T) {
	t.Skip()

	QueueSlice := NewQueueSlice()
	const SIZE = 10000000
	PrintMemUsage()
	fmt.Println(QueueSlice.String())
	now := time.Now()
	for i := 0; i < SIZE; i++ {
		//	val := int64(64384136813)
		QueueSlice.Push(i)
	}
	fmt.Printf("fill %s\n", time.Since(now))
	PrintMemUsage()
	fmt.Println(QueueSlice.String())
	now = time.Now()
	for QueueSlice.Len() > 0 {
		QueueSlice.Pop()
	}
	fmt.Printf("drain %s\n", time.Since(now))
	runtime.GC()
	PrintMemUsage()
	fmt.Println(QueueSlice.String())

}

/*
	Alloc = 0 MiB	TotalAlloc = 0 MiB	Sys = 10 MiB	NumGC = 0
	Capacity: 0 Length: 0
	Fill 39.882545762s
	Alloc = 222 MiB	TotalAlloc = 481 MiB	Sys = 397 MiB	NumGC = 11
	Capacity: 10597376 Length: 10000000
	Drain 350.078977ms
	Alloc = 157 MiB	TotalAlloc = 481 MiB	Sys = 397 MiB	NumGC = 12
	Capacity: 10597376 Length: 0
*/

func BenchmarkQueueSlicePush(b *testing.B) {
	q := NewQueueSlice()

	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
}

func BenchmarkQueueSlicePop(b *testing.B) {
	q := NewQueueSlice()

	for i := 0; i < b.N; i++ {
		q.Push(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}

func BenchmarkQueueSlice10Cycles(b *testing.B) {
	const M = 10

	q := NewQueueSlice()

	for j := 0; j < M; j++ {
		for i := 0; i < b.N; i++ {
			q.Push(i)
		}

		for i := 0; i < b.N; i++ {
			q.Pop()
		}
	}

	//	b.Logf("final cap: %v after %v cycles of (%v push + %v pop)", q.Len(), M, b.N, b.N)

	for j := 0; j < M; j++ {
		for i := 0; i < b.N; i++ {
			q.Push(i)
			q.Pop()
		}
	}

	//	b.Logf("final cap: %v after %v cycles of (%v push+pop)", q.Len(), M, b.N)
}

type node struct {
	data int
	next *node
}

type QueueSll struct {
	head  *node
	tail  *node
	count int
	lock  *sync.Mutex
}

func NewQueueSll() *QueueSll {
	q := &QueueSll{}
	q.lock = &sync.Mutex{}
	return q
}

func (q *QueueSll) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.count
}

func (q *QueueSll) Push(item int) {
	q.lock.Lock()
	defer q.lock.Unlock()

	n := &node{data: item}

	if q.tail == nil {
		q.tail = n
		q.head = n
	} else {
		q.tail.next = n
		q.tail = n
	}
	q.count++
}

func (q *QueueSll) Pop() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.head == nil {
		return 0
	}

	n := q.head
	q.head = n.next

	if q.head == nil {
		q.tail = nil
	}
	q.count--

	return n.data
}

func (q *QueueSll) Peek() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	n := q.head
	if n == nil {
		return 0
	}

	return n.data
}

func TestQueueSllLocal(t *testing.T) {
	t.Skip()

	queue := NewQueueSll()
	const SIZE = 10000000
	PrintMemUsage()
	fmt.Println(queue.Len())
	now := time.Now()
	for i := 0; i < SIZE; i++ {
		//	val := int64(64384136813)
		queue.Push(i)
	}
	fmt.Printf("fill %s\n", time.Since(now))
	PrintMemUsage()
	fmt.Println(queue.Len())
	now = time.Now()
	for queue.Len() > 0 {
		queue.Pop()
	}
	fmt.Printf("drain %s\n", time.Since(now))
	runtime.GC()
	PrintMemUsage()
	fmt.Println(queue.Len())
}

/*
	Alloc = 0 MiB	TotalAlloc = 0 MiB	Sys = 10 MiB	NumGC = 0
	Size 0
	Fill 2.322401648s
	Alloc = 229 MiB	TotalAlloc = 229 MiB	Sys = 260 MiB	NumGC = 6
	10000000
	Drain 694.36594ms
	Alloc = 0 MiB	TotalAlloc = 229 MiB	Sys = 261 MiB	NumGC = 7
	Size 0
*/

func BenchmarkQueueSllPush(b *testing.B) {
	q := NewQueueSll()

	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
}

func BenchmarkQueueSllPop(b *testing.B) {
	q := NewQueueSll()

	for i := 0; i < b.N; i++ {
		q.Push(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}

func BenchmarkQueueSll10Cycles(b *testing.B) {
	const M = 10

	q := NewQueueSll()

	for j := 0; j < M; j++ {
		for i := 0; i < b.N; i++ {
			q.Push(i)
		}

		for i := 0; i < b.N; i++ {
			q.Pop()
		}
	}

	//	b.Logf("final cap: %v after %v cycles of (%v push + %v pop)", q.Len(), M, b.N, b.N)

	for j := 0; j < M; j++ {
		for i := 0; i < b.N; i++ {
			q.Push(i)
			q.Pop()
		}
	}

	//	b.Logf("final cap: %v after %v cycles of (%v push+pop)", q.Len(), M, b.N)
}
