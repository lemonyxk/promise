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
	resultCh chan Result
	errCh    chan Error

	thenLink []Resolve

	then  bool
	catch bool
}

func (p *race) Then(fn Resolve) *race {
	if !p.then {
		p.then = true
		go func() {
			if result := <-p.resultCh; result != nil {
				fn(result)
				for i := 0; i < len(p.thenLink); i++ {
					p.thenLink[i](nil)
				}
			}
		}()
	} else {
		p.thenLink = append(p.thenLink, fn)
	}
	return p
}

func (p *race) Catch(fn Reject) {
	if !p.catch {
		p.catch = true
		go func() {
			if err := <-p.errCh; err != nil {
				fn(err)
			}
		}()
	}
}

func Race(promises ...*promise) *race {

	var p = &race{}

	p.resultCh = make(chan Result, 1)
	p.errCh = make(chan Error, 1)

	var sucCounter int32 = 0
	var errCounter int32 = 0

	go func() {
		for i := 0; i < len(promises); i++ {
			var index = i
			promises[index].Then(func(result Result) {
				if atomic.AddInt32(&sucCounter, 1) == 1 {
					p.resultCh <- result
					p.errCh <- nil
				}
			}).Catch(func(err Error) {
				if atomic.AddInt32(&errCounter, 1) == 1 {
					p.errCh <- err
					p.resultCh <- nil
				}
			})
		}
	}()

	return p
}
