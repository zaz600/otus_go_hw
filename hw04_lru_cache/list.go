package hw04_lru_cache //nolint:golint,stylecheck

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
	size  int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	// nil <- (prev) front <-> ... <-> elem <-> ... <-> back (next) -> nil
	item := &ListItem{
		Value: v,
		Next:  l.front,
	}

	if l.size == 0 {
		l.front = item
		l.back = item
	} else {
		l.front.Prev = item
		l.front = item
	}
	l.size++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	// nil <- (prev) front <-> ... <-> elem <-> ... <-> back (next) -> nil
	item := &ListItem{
		Value: v,
		Prev:  l.back,
	}
	if l.size == 0 {
		l.front = item
		l.back = item
	} else {
		l.back.Next = item
		l.back = item
	}
	l.size++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if l.size == 1 {
		l.front = nil
		l.back = nil
		l.size--
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	l.size--
}

func (l *list) MoveToFront(i *ListItem) {
	// nil <- (prev) front <-> ... <-> elem <-> ... <-> back (next) -> nil
	if i == nil {
		return
	}
	l.PushFront(i.Value)
	l.Remove(i)
}

func NewList() List {
	return &list{}
}
