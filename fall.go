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
	results[T, P]
}

func Fall[T any, P any](promises ...Promise[T, P]) Results[T, P] {

	var p = new(fall[T, P])
	p.index = -1
	p.chList = append(p.chList, make(chan []T, 1))
	p.ehList = append(p.ehList, make(chan P, 1))

	p.fn = func(index int) {
		var sucCounter int32 = 0
		var errCounter int32 = 0
		var results = make([]T, len(promises))

		var pi = 0

		var fn func(int)

		fn = func(pi int) {
			promises[pi].Then(func(result T) {
				results[pi] = result
				if atomic.AddInt32(&sucCounter, 1) == int32(len(promises)) {
					p.chList[index] <- results
				} else {
					pi++
					fn(pi)
				}
			}).Catch(func(err P) {
				if atomic.AddInt32(&errCounter, 1) == 1 {
					p.ehList[index] <- err
				}
			})
		}

		fn(pi)
	}

	p.fn(0)

	return p
}
