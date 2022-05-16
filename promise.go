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
	Then(func(T)) Promise[T, P]
	Catch(func(P)) Promise[T, P]
	Finally(func())
}

type promise[T any, P any] struct {
	fn func()

	ch chan T
	eh chan P

	done chan bool
}

func (p *promise[T, P]) Then(then func(T)) Promise[T, P] {
	var done = <-p.done
	if done {
		var res = <-p.ch
		then(res)
		p.ch <- res
		p.done <- done
		return p
	} else {
		p.done <- done
		return p
	}
}

func (p *promise[T, P]) Catch(catch func(P)) Promise[T, P] {

	var done = <-p.done
	if done {
		p.done <- done
		return p
	} else {
		var err = <-p.eh
		catch(err)
		p.eh <- err
		p.done <- done
		return p
	}
}

func (p *promise[T, P]) Finally(finally func()) {
	finally()
}

func New[T any, P any](state func(resolve func(T), reject func(P))) Promise[T, P] {

	var p = new(promise[T, P])
	p.ch = make(chan T, 1)
	p.eh = make(chan P, 1)
	p.done = make(chan bool, 1)

	p.fn = func() {
		var counter int32 = 0
		// just one can be exec
		var resolve = func(result T) {
			if atomic.AddInt32(&counter, 1) == 1 {
				p.ch <- result
				p.done <- true
			}
		}
		var reject = func(err P) {
			if atomic.AddInt32(&counter, 1) == 1 {
				p.eh <- err
				p.done <- false
			}
		}
		state(resolve, reject)
	}

	p.fn()

	return p
}
