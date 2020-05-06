package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int                          // длина списка
	Front() *listItem                  // первый listItem
	Back() *listItem                   // последний listItem
	PushFront(v interface{}) *listItem // добавить значение в начало
	PushBack(v interface{}) *listItem  // добавить значение в конец
	Remove(i *listItem)                // удалить элемент
	MoveToFront(i *listItem)           // переместить элемент в начало
}

type listItem struct {
	Value interface{} // значение
	Next  *listItem   // следующий элемент
	Prev  *listItem   // предыдущий элемент
}

type list struct {
	front  *listItem // первый элемент
	back   *listItem // последний элемент
	length int       // длина
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *listItem {
	return l.front
}

func (l *list) Back() *listItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *listItem {
	newItem := &listItem{
		Value: v,
		Next:  nil,
	}

	if l.front == nil {
		newItem.Prev = nil
		l.back = newItem
	} else {
		newItem.Prev = l.front
		l.front.Next = newItem
	}

	l.front = newItem
	l.length++

	return newItem
}

func (l *list) PushBack(v interface{}) *listItem {
	newItem := &listItem{
		Value: v,
		Next:  l.back,
		Prev:  nil,
	}

	if l.back == nil {
		l.front = newItem
	} else {
		l.back.Prev = newItem
	}

	l.back = newItem
	l.length++

	return newItem
}

func (l *list) Remove(i *listItem) {
	if i != nil {
		if i.Next != nil {
			i.Next.Prev = i.Prev
		} else {
			l.front = i.Prev
		}

		if i.Prev != nil {
			i.Prev.Next = i.Next
		} else {
			l.back = i.Next
		}

		l.length--
	}
}

func (l *list) MoveToFront(i *listItem) {
	if i != nil {
		l.Remove(i)
		l.PushFront(i.Value)
	}
}

func NewList() List {
	return &list{}
}
