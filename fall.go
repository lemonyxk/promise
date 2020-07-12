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

type fall struct {
	promise
}

func Fall(promises ...Promise) Promise {

	var p fall
	p.ch = make(chan data, 1)

	p.fn = func() {
		var sucCounter int32 = 0
		var errCounter int32 = 0
		var results = make([]Result, len(promises))

		var index = 0

		var fn func(int)

		fn = func(index int) {
			promises[index].Then(func(result Result) {
				results[index] = result
				if atomic.AddInt32(&sucCounter, 1) == int32(len(promises)) {
					p.ch <- data{res: results, err: nil}
				} else {
					index++
					fn(index)
				}
			}).Catch(func(err Error) {
				if atomic.AddInt32(&errCounter, 1) == 1 {
					p.ch <- data{res: nil, err: err}
				}
			})
		}

		fn(index)
	}

	return p
}
