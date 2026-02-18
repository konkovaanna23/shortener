package pool

import (
	"sync"
)

type Resettable interface {
	Reset()
}

type Pool[T Resettable] struct {
	p sync.Pool
}

func NewPool[T Resettable](factory func() T) *Pool[T] {
	return &Pool[T]{
		p: sync.Pool{
			New: func() any {
				return factory()
			},
		},
	}

}

func (pl *Pool[T]) Get() T {
	return pl.p.Get().(T)
}

func (pl *Pool[T]) Put(v T) {
	if any(v) == nil {
		return
	}
	v.Reset()
	pl.p.Put(v)
}
