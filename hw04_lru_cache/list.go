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
	if l.len == 0 {
		l.front = item
		l.back = item
		l.len++
		return item
	} else {
		item.Next = l.front
		l.front.Prev = item
		l.front = item
		l.len++
	}
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.len == 0 {
		l.front = item
		l.back = item
		l.len++
		return item
	} else {
		item.Prev = l.back
		l.back.Next = item
		l.back = item
		l.len++
	}
	return item
}

func (l *list) Remove(item *ListItem) {
	if item == nil || l.len == 0 {
		return
	}
	foundItem := l.find(item)
	if foundItem == nil {
		return
	}

	if item == l.front {
		l.front = foundItem.Next
		if foundItem.Next != nil {
			foundItem.Next.Prev = nil
		}
	} else if item == l.back {
		l.back = foundItem.Prev
		if foundItem.Prev != nil {
			foundItem.Prev.Next = nil
		}
	} else {
		foundItem.Prev.Next = foundItem.Next
		foundItem.Next.Prev = foundItem.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || l.len == 0 {
		return
	}
	foundItem := l.find(i)
	if foundItem == nil || foundItem == l.front {
		return
	}

	if foundItem == l.back {
		l.back = foundItem.Prev
		foundItem.Prev.Next = nil
	} else {
		foundItem.Prev.Next = foundItem.Next
		foundItem.Next.Prev = foundItem.Prev
	}

	foundItem.Next = l.front
	foundItem.Prev = nil
	l.front.Prev = foundItem
	l.front = foundItem
}

func (l *list) find(item *ListItem) *ListItem {
	current := l.front
	for current != nil {
		if current == item {
			return current
		}
		current = current.Next
	}
	return nil
}
