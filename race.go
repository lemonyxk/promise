/**
* @program: promise
*
* @description:
*
* @author: lemo
*
* @create: 2020-07-12 01:50
**/

package promise

import (
	"sync/atomic"
)

type race[T any, P any] struct {
	promise[T, P]
}

func Race[T any, P any](promises ...Promise[T, P]) Promise[T, P] {

	var p race[T, P]
	p.ch = make(chan T, 1)
	p.eh = make(chan P, 1)

	p.fn = func() {
		var sucCounter int32 = 0
		var errCounter int32 = 0

		for i := 0; i < len(promises); i++ {
			var index = i
			go func() {
				promises[index].Then(func(result T) {
					if atomic.AddInt32(&sucCounter, 1) == 1 {
						p.ch <- result
					}
				}).Catch(func(err P) {
					if atomic.AddInt32(&errCounter, 1) == 1 {
						p.eh <- err
					}
				})
			}()
		}
	}

	return p
}
