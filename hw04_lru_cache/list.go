package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	first *ListItem
	last  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if i == l.first {
		l.first = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i == l.last {
		l.last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := &ListItem{Value: v, Next: l.first}

	if l.len == 0 {
		l.last = i
	} else {
		l.first.Prev = i
	}

	l.first = i
	l.len++
	return i
}

func (l *list) PushBack(v interface{}) *ListItem {
	i := &ListItem{Value: v, Prev: l.last}

	if l.len == 0 {
		l.first = i
	} else {
		l.last.Next = i
	}

	l.last = i
	l.len++
	return i
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.first {
		return
	}

	if i == l.last {
		l.last = i.Prev
		i.Prev.Next = nil
	} else {
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}

	i.Prev = nil
	i.Next = l.first
	l.first.Prev = i
	l.first = i
}

func NewList() List {
	return new(list)
}
