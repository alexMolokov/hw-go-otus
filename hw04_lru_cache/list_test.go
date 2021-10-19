package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("cap one list", func(t *testing.T) {
		l := NewList()
		l.PushBack("10")

		first := l.Front()
		last := l.Back()

		require.Equal(t, 1, l.Len())
		require.Equal(t, first, last)

		l.MoveToFront(last)
		first = l.Front()
		last = l.Back()

		require.Equal(t, first, last)
		require.Equal(t, "10", first.Value.(string))
		require.Nil(t, first.Next)
		require.Nil(t, first.Prev)

		l.Remove(last)
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("can put struct", func(t *testing.T) {
		type T struct {
			str string
			i   int
		}
		l := NewList()

		v1 := T{str: "10", i: 10}
		v2 := T{str: "20", i: 20}

		first := l.PushBack(v1)
		last := l.PushBack(v2)

		require.Equal(t, 2, l.Len())
		l.MoveToFront(first)
		require.Equal(t, first, l.Front())
		l.MoveToFront(last)
		require.Equal(t, last, l.Front())

		front := l.Front()
		data := front.Value.(T)

		require.Equal(t, "20", data.str)
		require.Equal(t, 20, data.i)
	})

	t.Run("complex", func(t *testing.T) {
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
	})
}
