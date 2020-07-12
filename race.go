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

type race struct {
	promise
}

func Race(promises ...Promise) Promise {

	var p race
	p.ch = make(chan data, 1)

	p.fn = func() {
		var sucCounter int32 = 0
		var errCounter int32 = 0

		for i := 0; i < len(promises); i++ {
			var index = i
			go func() {
				promises[index].Then(func(result Result) {
					if atomic.AddInt32(&sucCounter, 1) == 1 {
						p.ch <- data{res: result, err: nil}
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
