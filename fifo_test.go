package fifo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFIFO(t *testing.T) {
	f := New[int](2)

	t.Logf("f: %+v", f.arr)

	f.Push(1)
	assert.Equal(t, 1, f.Pop())

	f.Push(2)
	assert.Equal(t, 2, f.Pop())

	f.Push(3)

	assert.Equal(t, 1, f.Len())
	assert.Equal(t, 2, f.Cap())

	assert.Equal(t, 3, f.Pop())

	assert.Equal(t, 0, f.Len())
	assert.Equal(t, 2, f.Cap())

	f.Push(4)
	f.Push(5)
	f.Push(6)

	t.Logf("f: %+v", f.arr)

	assert.Equal(t, 3, f.Len())
	assert.Equal(t, 6, f.Cap())

	assert.Equal(t, 4, f.Pop())

	t.Logf("f: %+v", f.arr)

	assert.Equal(t, 2, f.Len())
	assert.Equal(t, 4, f.Cap())
}

func BenchmarkFIFOPush(b *testing.B) {
	q := New[int](1 << 10)

	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
}

func BenchmarkFIFOPop(b *testing.B) {
	q := New[int](1 << 10)

	for i := 0; i < b.N; i++ {
		q.Push(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}

func BenchmarkFIFO10Cycles(b *testing.B) {
	const M = 10

	q := New[int](1 << 10)

	for j := 0; j < M; j++ {
		for i := 0; i < b.N; i++ {
			q.Push(i)
		}

		for i := 0; i < b.N; i++ {
			q.Pop()
		}
	}

	//	b.Logf("final cap: %v after %v cycles of (%v push + %v pop)", q.Cap(), M, b.N, b.N)

	for j := 0; j < M; j++ {
		for i := 0; i < b.N; i++ {
			q.Push(i)
			q.Pop()
		}
	}

	//	b.Logf("final cap: %v after %v cycles of (%v push+pop)", q.Cap(), M, b.N)
}
