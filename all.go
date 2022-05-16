/**
* @program: promise
*
* @description:
*
* @author: lemo
*
* @create: 2020-07-12 01:47
**/

package promise

import (
	"sync/atomic"
)

type all[T any, P any] struct {
	promises[T, P]
}

func All[T any, P any](promises ...Promise[T, P]) Promises[T, P] {

	var p = new(all[T, P])
	p.ch = make(chan []T, 1)
	p.eh = make(chan P, 1)
	p.done = make(chan bool, 1)

	p.fn = func() {
		var sucCounter int32 = 0
		var errCounter int32 = 0
		var results = make([]T, len(promises))

		for i := 0; i < len(promises); i++ {
			var pi = i
			go func() {
				promises[pi].Then(func(result T) {
					results[pi] = result
					if atomic.AddInt32(&sucCounter, 1) == int32(len(promises)) {
						p.ch <- results
						p.done <- true
					}
				}).Catch(func(err P) {
					if atomic.AddInt32(&errCounter, 1) == 1 {
						p.eh <- err
						p.done <- false
					}
				})
			}()
		}
	}

	p.fn()

	return p
}
