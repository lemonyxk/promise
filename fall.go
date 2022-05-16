/**
* @program: promise
*
* @description:
*
* @author: lemo
*
* @create: 2020-07-12 13:10
**/

package promise

import (
	"sync/atomic"
)

type fall[T any, P any] struct {
	promises[T, P]
}

func Fall[T any, P any](promises ...Promise[T, P]) Promises[T, P] {

	var p = new(fall[T, P])
	p.ch = make(chan []T, 1)
	p.eh = make(chan P, 1)
	p.done = make(chan bool, 1)

	p.fn = func() {
		var sucCounter int32 = 0
		var errCounter int32 = 0
		var results = make([]T, len(promises))

		var pi = 0

		var fn func(int)

		fn = func(pi int) {
			promises[pi].Then(func(result T) {
				results[pi] = result
				if atomic.AddInt32(&sucCounter, 1) == int32(len(promises)) {
					p.ch <- results
					p.done <- true
				} else {
					pi++
					fn(pi)
				}
			}).Catch(func(err P) {
				if atomic.AddInt32(&errCounter, 1) == 1 {
					p.eh <- err
					p.done <- false
				}
			})
		}

		fn(pi)
	}

	p.fn()

	return p
}
