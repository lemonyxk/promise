/**
* @program: promise
*
* @description:
*
* @author: lemo
*
* @create: 2020-07-12 01:46
**/

package promise

import (
	"sync/atomic"
)

type Promise[T any, P any] interface {
	Then(func(result T)) Catch[P]
}

type Promises[T any, P any] interface {
	Then(func(result []T)) Catch[P]
}

type Catch[T any] interface {
	Catch(func(err T))
}

type promise[T any, P any] struct {
	ch chan T
	eh chan P
	fn func()
}

type promises[T any, P any] struct {
	ch chan []T
	eh chan P
	fn func()
}

func (p promise[T, P]) Then(then func(T)) Catch[P] {
	p.fn()

	select {
	case res := <-p.ch:
		var c = catch[P]{err: empty[P](), run: false}
		then(res)
		return c
	case err := <-p.eh:
		var c = catch[P]{err: err, run: true}
		return c
	}
}

func (p promises[T, P]) Then(then func([]T)) Catch[P] {
	p.fn()

	select {
	case res := <-p.ch:
		var c = catch[P]{err: empty[P](), run: false}
		then(res)
		return c
	case err := <-p.eh:
		var c = catch[P]{err: err, run: true}
		return c
	}
}

type catch[T any] struct {
	err T
	run bool
}

func (c catch[T]) Catch(catch func(T)) {
	if c.run {
		catch(c.err)
	}
}

func New[T any, P any](state State[T, P]) Promise[T, P] {

	var p promise[T, P]
	p.ch = make(chan T, 1)
	p.eh = make(chan P, 1)

	p.fn = func() {
		var counter int32 = 0
		// just one can be exec
		var resolve = func(result T) {
			if atomic.AddInt32(&counter, 1) == 1 {
				p.ch <- result
			}
		}
		var reject = func(err P) {
			if atomic.AddInt32(&counter, 1) == 1 {
				p.eh <- err
			}
		}
		state(resolve, reject)
	}

	return p
}
