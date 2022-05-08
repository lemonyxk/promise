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

	var p fall[T, P]
	p.ch = make(chan []T, 1)
	p.eh = make(chan P, 1)

	p.fn = func() {
		var sucCounter int32 = 0
		var errCounter int32 = 0
		var results = make([]T, len(promises))

		var index = 0

		var fn func(int)

		fn = func(index int) {
			promises[index].Then(func(result T) {
				results[index] = result
				if atomic.AddInt32(&sucCounter, 1) == int32(len(promises)) {
					p.ch <- results
				} else {
					index++
					fn(index)
				}
			}).Catch(func(err P) {
				if atomic.AddInt32(&errCounter, 1) == 1 {
					p.eh <- err
				}
			})
		}

		fn(index)
	}

	return p
}
