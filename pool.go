package gosocket

import (
	"sync"
)

type Pool[T any] struct {
	p sync.Pool
}

func NewPool[T any](fn func() T) *Pool[T] {
	return &Pool[T]{
		p: sync.Pool{
			New: func() any {
				return fn()
			},
		},
	}
}

func (p *Pool[T]) Get() T {
	return p.p.Get().(T)
}

func (p *Pool[T]) Put(v T) {
	p.p.Put(v)
}
