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
	"sync"
	"sync/atomic"
)

type Promise[T any, P any] interface {
	Then(func(result T)) Catch[P]
}

type promise[T any, P any] struct {
	fn func(index int)

	index    int32
	chList   []chan T
	ehList   []chan P
	thenList []chan T

	result    T
	err       P
	hasResult bool
	hasErr    bool

	mux sync.Mutex
}

func (p *promise[T, P]) Then(then func(T)) Catch[P] {

	var index = atomic.AddInt32(&p.index, 1)
	// not first
	if index > 0 {
		p.mux.Lock()

		if p.hasResult {
			then(p.result)
			p.mux.Unlock()
			var c = &catch[P]{err: empty[P](), run: false}
			return c
		}

		if p.hasErr {
			p.mux.Unlock()
			var c = &catch[P]{err: p.err, run: true}
			return c
		}

		var t = make(chan T, 1)

		p.thenList = append(p.thenList, t)
		p.mux.Unlock()

		var res = <-t
		then(res)

		var c = &catch[P]{err: empty[P](), run: false}
		return c
	}

	select {
	case res := <-p.chList[index]:
		p.mux.Lock()
		then(res)
		p.result = res
		p.hasResult = true
		for i := 0; i < len(p.thenList); i++ {
			p.thenList[i] <- res
		}
		p.mux.Unlock()
		var c = &catch[P]{err: empty[P](), run: false}
		return c
	case err := <-p.ehList[index]:
		p.mux.Lock()
		p.err = err
		p.hasErr = true
		p.mux.Unlock()
		var c = &catch[P]{err: err, run: true}
		return c
	}
}

type Catch[T any] interface {
	Catch(func(err T))
}

type catch[T any] struct {
	err T
	run bool
}

func (c *catch[T]) Catch(catch func(T)) {
	if c.run {
		catch(c.err)
	}
}

func New[T any, P any](state func(resolve func(T), reject func(P))) Promise[T, P] {

	var p = new(promise[T, P])
	p.index = -1
	p.chList = append(p.chList, make(chan T, 1))
	p.ehList = append(p.ehList, make(chan P, 1))

	p.fn = func(index int) {
		var counter int32 = 0
		// just one can be exec
		var resolve = func(result T) {
			if atomic.AddInt32(&counter, 1) == 1 {
				p.chList[index] <- result
			}
		}
		var reject = func(err P) {
			if atomic.AddInt32(&counter, 1) == 1 {
				p.ehList[index] <- err
			}
		}
		state(resolve, reject)
	}

	p.fn(0)

	return p
}
