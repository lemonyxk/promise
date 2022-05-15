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
	results[T, P]
}

func All[T any, P any](promises ...Promise[T, P]) Results[T, P] {

	var p = new(all[T, P])
	p.index = -1
	p.chList = append(p.chList, make(chan []T, 1))
	p.ehList = append(p.ehList, make(chan P, 1))

	p.fn = func(index int) {
		var sucCounter int32 = 0
		var errCounter int32 = 0
		var results = make([]T, len(promises))

		for i := 0; i < len(promises); i++ {
			var pi = i
			go func() {
				promises[pi].Then(func(result T) {
					results[pi] = result
					if atomic.AddInt32(&sucCounter, 1) == int32(len(promises)) {
						p.chList[index] <- results
					}
				}).Catch(func(err P) {
					if atomic.AddInt32(&errCounter, 1) == 1 {
						p.ehList[index] <- err
					}
				})
			}()
		}
	}

	p.fn(0)

	return p
}
