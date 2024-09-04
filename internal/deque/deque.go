package deque

import "container/list"

// Deque 是一个使用泛型的双端队列
type Deque[T any] struct {
	items *list.List
}

// NewDeque 创建一个新的 Deque
func NewDeque[T any]() *Deque[T] {
	return &Deque[T]{items: list.New()}
}

// PushFront 向队列前端添加元素
func (d *Deque[T]) PushFront(item T) {
	d.items.PushFront(item)
}

// PushBack 向队列后端添加元素
func (d *Deque[T]) PushBack(item T) {
	d.items.PushBack(item)
}

// PopFront 从队列前端移除元素
func (d *Deque[T]) PopFront() (T, bool) {
	if d.items.Len() == 0 {
		var zero T
		return zero, false
	}
	element := d.items.Front()
	d.items.Remove(element)
	return element.Value.(T), true
}

// PopBack 从队列后端移除元素
func (d *Deque[T]) PopBack() (T, bool) {
	if d.items.Len() == 0 {
		var zero T
		return zero, false
	}
	element := d.items.Back()
	d.items.Remove(element)
	return element.Value.(T), true
}

// IsEmpty 检查队列是否为空
func (d *Deque[T]) IsEmpty() bool {
	return d.items.Len() == 0
}

// Size 获取队列的大小
func (d *Deque[T]) Size() int {
	return d.items.Len()
}
