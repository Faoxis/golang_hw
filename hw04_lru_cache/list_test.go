package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyList(t *testing.T) {
	l := NewList()

	require.Equal(t, 0, l.Len())
	require.Nil(t, l.Front())
	require.Nil(t, l.Back())
}

func TestComplexOperations(t *testing.T) {
	l := NewList()

	l.PushFront(10) // [10]
	l.PushBack(20)  // [10, 20]
	l.PushBack(30)  // [10, 20, 30]
	require.Equal(t, 3, l.Len())

	middle := l.Front().Next // 20
	l.Remove(middle)         // [10, 30]
	require.Equal(t, 2, l.Len())

	for i, v := range [...]int{40, 50, 60, 70, 80} {
		if i%2 == 0 {
			l.PushFront(v)
		} else {
			l.PushBack(v)
		}
	} // [80, 60, 40, 10, 30, 50, 70]

	require.Equal(t, 7, l.Len())
	require.Equal(t, 80, l.Front().Value)
	require.Equal(t, 70, l.Back().Value)

	l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
	l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

	elems := make([]int, 0, l.Len())
	for i := l.Front(); i != nil; i = i.Next {
		elems = append(elems, i.Value.(int))
	}
	require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
}

func TestOneElement(t *testing.T) {
	l := NewList()

	l.PushFront(10)
	require.Equal(t, 1, l.Len())
	require.Equal(t, 10, l.Front().Value)
	require.Equal(t, 10, l.Back().Value)
	require.Nil(t, l.Front().Next)
	require.Nil(t, l.Front().Prev)
}

func TestPushBack(t *testing.T) {
	l := NewList()

	l.PushBack(10)
	l.PushBack(20)
	l.PushBack(30)

	require.Equal(t, 3, l.Len())
	require.Equal(t, 10, l.Front().Value)
	require.Equal(t, 30, l.Back().Value)
}

func TestMoveToFront(t *testing.T) {
	l := NewList()

	l.PushBack(10) // [10]
	l.PushBack(20) // [10,20]
	l.PushBack(30) // [10,20,30]

	l.MoveToFront(l.Back()) // [30,10,20]
	require.Equal(t, 30, l.Front().Value)
	require.Equal(t, 20, l.Back().Value)

	l.MoveToFront(l.Back()) // [20,30,10]
	require.Equal(t, 20, l.Front().Value)
	require.Equal(t, 10, l.Back().Value)
}

func TestRemoveMiddleElement(t *testing.T) {
	l := NewList()

	l.PushBack(10)
	l.PushBack(20)
	l.PushBack(30)
	l.Remove(l.Front().Next) // [10,30]

	require.Equal(t, 2, l.Len())
	require.Equal(t, 10, l.Front().Value)
	require.Equal(t, 30, l.Back().Value)
}

func TestMoveMiddleToFront(t *testing.T) {
	l := NewList()

	l.PushBack(10)
	l.PushBack(20)
	l.PushBack(30)
	l.MoveToFront(l.Front().Next) // [20,10,30]

	require.Equal(t, 3, l.Len())
	require.Equal(t, 20, l.Front().Value)
	require.Equal(t, 30, l.Back().Value)
}

func TestMultipleOperations(t *testing.T) {
	l := NewList()

	l.PushFront(10)              // [10]
	l.PushBack(20)               // [10,20]
	l.PushBack(30)               // [10,20,30]
	l.PushFront(40)              // [40,10,20,30]
	l.Remove(l.Front().Next)     // [40,20,30]
	l.PushBack(50)               // [40,20,30,50]
	l.MoveToFront(l.Back().Prev) // [30,40,20,50]

	require.Equal(t, 4, l.Len())
	require.Equal(t, 30, l.Front().Value)
	require.Equal(t, 50, l.Back().Value)
}

func TestEmptyListOperations(t *testing.T) {
	l := NewList()

	require.Equal(t, 0, l.Len())
	require.Nil(t, l.Front())
	require.Nil(t, l.Back())

	l.Remove(nil)      // не должно вызывать панику
	l.MoveToFront(nil) // не должно вызывать панику
}

func TestRemoveAllElementsSequentially(t *testing.T) {
	l := NewList()

	l.PushBack(10) // [10]
	l.PushBack(20) // [10,20]
	l.PushBack(30) // [10,20,30]

	l.Remove(l.Front()) // [20,30]
	require.Equal(t, 2, l.Len())
	require.Equal(t, 20, l.Front().Value)

	l.Remove(l.Front()) // [30]
	require.Equal(t, 1, l.Len())
	require.Equal(t, 30, l.Front().Value)

	l.Remove(l.Front()) // []
	require.Equal(t, 0, l.Len())
	require.Nil(t, l.Front())
}

func TestComplexMoveToFront(t *testing.T) {
	l := NewList()

	l.PushBack(10)          // [10]
	l.PushBack(20)          // [10,20]
	l.PushBack(30)          // [10,20,30]
	l.PushBack(40)          // [10,20,30,40]
	l.MoveToFront(l.Back()) // [40,10,20,30]
	l.MoveToFront(l.Back()) // [30,40,10,20]
	l.MoveToFront(l.Back()) // [20,30,40,10]

	require.Equal(t, 4, l.Len())
	require.Equal(t, 20, l.Front().Value)
	require.Equal(t, 10, l.Back().Value)
}

func TestAlternatingPushFrontAndBack(t *testing.T) {
	l := NewList()

	l.PushFront(10) // [10]
	l.PushBack(20)  // [10,20]
	l.PushFront(30) // [30,10,20]
	l.PushBack(40)  // [30,10,20,40]
	l.PushFront(50) // [50,30,10,20,40]

	require.Equal(t, 5, l.Len())
	require.Equal(t, 50, l.Front().Value)
	require.Equal(t, 40, l.Back().Value)
	require.Equal(t, 10, l.Front().Next.Next.Value)
}

func TestRemovingOneElement(t *testing.T) {
	l := NewList()

	el := l.PushFront(10)
	l.Remove(el)

	require.Equal(t, 0, l.Len())

	el = l.Front()
	require.Nil(t, el)

	el = l.Back()
	require.Nil(t, el)
}
