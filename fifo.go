package fifo

import (
	"fmt"
	//"log"
)

type (
	FIFO[T any] struct {
		*arr[T]
	}

	arr[T any] struct {
		b    []T
		r, w int
		p    *arr[T]

		mask int
	}
)

func New[T any](n int) *FIFO[T] {
	return &FIFO[T]{
		arr: newArrN[T](n, (*arr[T])(nil)),
	}
}

func (f *FIFO[T]) Push(el T) {
	if f.arr == nil {
		f.arr = newArr((*arr[T])(nil))
	}

	if f.w > 2 * f.localCap() && f.localLen() < f.localCap() / 2 && f.localCap() > 16 { // shrink
		f.arr = newArrN[T](f.localCap() / 2, f.arr)
	}

	if f.push(el) {
		return
	}

	f.arr = newArr[T](f.arr)

	_ = f.push(el)

	return
}

func (f *FIFO[T]) Pop() (el T) {
	if f.arr == nil {
		return
	}

	el, _ = f.pop()

	return
}

func newArr[T any](p *arr[T]) *arr[T] {
	n := 16

	if p != nil {
		n = len(p.b) * 2
	}

	return newArrN[T](n, p)
}

func newArrN[T any](n int, p *arr[T]) *arr[T] {
	if n&(n-1) != 0 {
		panic(n)
	}

	return &arr[T]{
		b:    make([]T, n),
		p:    p,
		mask: n - 1,
	}
}

func (b *arr[T]) push(el T) bool {
	if (b.w + 1) == b.r+len(b.b) {
		return false
	}

	b.b[b.w&b.mask] = el
	b.w++

	return true
}

func (b *arr[T]) pop() (el T, ok bool) {
	if b.p != nil {
		el, ok = b.p.pop()

		if b.p.Len() == 0 {
			b.p = nil
		}

		return
	}

	if b.r == b.w {
		return el, false
	}

	el = b.b[b.r&b.mask]

	var zero T
	b.b[b.r&b.mask] = zero

	b.r++

	return el, true
}

func (b *arr[T]) Len() int {
	if b == nil {
		return 0
	}

	return b.w - b.r + b.p.Len()
}

func (b *arr[T]) Cap() int {
	if b == nil {
		return 0
	}

	return len(b.b) + b.p.Cap()
}

func (b *arr[T]) localLen() int {
	if b == nil {
		return 0
	}

	return b.w - b.r
}

func (b *arr[T]) localCap() int {
	return cap(b.b)
}

func (b *arr[T]) String() string {
	if b == nil {
		return "nil"
	}

	return fmt.Sprintf("arr{r:%d, w:%d, b:%v, p:%v}", b.r, b.w, b.b, b.p)
}
