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

import "sync/atomic"

type all struct {
	resultCh chan []Result
	errCh    chan Error

	thenLink []Resolves

	then  bool
	catch bool
}

func (p *all) Then(fn Resolves) *all {
	if !p.then {
		p.then = true
		go func() {
			if results := <-p.resultCh; results != nil {
				fn(results)
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

func (p *all) Catch(fn Reject) {
	if !p.catch {
		p.catch = true
		go func() {
			if err := <-p.errCh; err != nil {
				fn(err)
			}
		}()
	}
}

func All(promises ...*promise) *all {

	var p = &all{}

	p.resultCh = make(chan []Result, 1)
	p.errCh = make(chan Error, 1)

	var results = make([]Result, len(promises))

	var sucCounter int32 = 0
	var errCounter int32 = 0

	go func() {
		for i := 0; i < len(promises); i++ {
			var index = i
			promises[index].Then(func(result Result) {
				results[index] = result
				if atomic.AddInt32(&sucCounter, 1) == int32(len(promises)) {
					p.resultCh <- results
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
