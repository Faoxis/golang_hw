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
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	switch l.len {
	case 0:
		l.front = item
		l.back = item
		l.len++
		return item
	default:
		item.Next = l.front
		l.front.Prev = item
		l.front = item
		l.len++
	}
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	switch l.len {
	case 0:
		l.front = item
		l.back = item
		l.len++
		return item
	default:
		item.Prev = l.back
		l.back.Next = item
		l.back = item
		l.len++
	}
	return item
}

func (l *list) Remove(item *ListItem) {
	switch l.len {
	case 0:
		return
	case 1:
		l.front = nil
		l.back = nil
	default:
		switch item {
		case nil:
			return
		case l.front:
			l.front = item.Next
			if item.Next != nil {
				item.Next.Prev = nil
			}
		case l.back:
			l.back = item.Prev
			if item.Prev != nil {
				item.Prev.Next = nil
			}
		default:
			item.Prev.Next = item.Next
			item.Next.Prev = item.Prev
		}
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || l.len == 0 || i == l.front {
		return
	}

	if i == l.back {
		l.back = i.Prev
		i.Prev.Next = nil
	} else {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	i.Next = l.front
	i.Prev = nil
	l.front.Prev = i
	l.front = i
}
