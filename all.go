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

type all struct {
	promise
}

func All(promises ...Promise) Promise {

	var p all
	p.ch = make(chan data, 1)

	p.fn = func() {
		var sucCounter int32 = 0
		var errCounter int32 = 0
		var results = make([]Result, len(promises))

		for i := 0; i < len(promises); i++ {
			var index = i
			go func() {
				promises[index].Then(func(result Result) {
					results[index] = result
					if atomic.AddInt32(&sucCounter, 1) == int32(len(promises)) {
						p.ch <- data{res: results, err: nil}
					}
				}).Catch(func(err Error) {
					if atomic.AddInt32(&errCounter, 1) == 1 {
						p.ch <- data{res: nil, err: err}
					}
				})
			}()
		}
	}

	return p
}
