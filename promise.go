/**
* @program: promise
*
* @description:
*
* @author: lemo
*
* @create: 2020-07-12 01:46
**/

package promise

import (
	"sync/atomic"
)

type promise struct {
	resolve Resolve
	reject  Reject

	resultCh chan Result
	errCh    chan Error

	thenLink []Resolve

	then  bool
	catch bool
}

func (p *promise) Then(resolve Resolve) *promise {
	if !p.then {
		p.then = true
		go func() {
			if result := <-p.resultCh; result != nil {
				resolve(result)
				for i := 0; i < len(p.thenLink); i++ {
					p.thenLink[i](nil)
				}
			}
		}()
	} else {
		p.thenLink = append(p.thenLink, resolve)
	}
	return p
}

func (p *promise) Catch(reject Reject) {
	if !p.catch {
		p.catch = true
		go func() {
			if err := <-p.errCh; err != nil {
				reject(err)
			}
		}()
	}
}

func New(fn State) *promise {
	var p = &promise{}

	p.resultCh = make(chan Result, 1)
	p.errCh = make(chan Error, 1)

	var counter int32 = 0

	p.resolve = func(result Result) {
		if atomic.AddInt32(&counter, 1) == 1 {
			p.resultCh <- result
			p.errCh <- nil
		}
	}
	p.reject = func(err Error) {
		if atomic.AddInt32(&counter, 1) == 1 {
			p.errCh <- err
			p.resultCh <- nil
		}
	}
	fn(p.resolve, p.reject)
	return p
}
