/**
* @program: lemo
*
* @description:
*
* @author: lemo
*
* @create: 2020-07-11 01:00
**/

package promise

import (
	"sync/atomic"
	"time"
)

type Resolve func(Result)
type Reject func(Error)
type Result interface{}
type Error interface{}

type State func(resolve Resolve, reject Reject)

type promise struct {
	resolve Resolve
	reject  Reject
	result  Result
	err     Error

	resultCh chan bool
	errCh    chan bool
	done     chan time.Time
	next     []Resolve
	then     bool
	catch    bool
	count    int32
}

func (p *promise) Then(resolve Resolve) *promise {
	if !p.then {
		p.then = true
		go func() {
			if <-p.resultCh {
				resolve(p.result)
				p.result = nil
				for i := 0; i < len(p.next); i++ {
					p.next[i](p.result)
				}
				p.done <- time.Now()
			}
		}()
	} else {
		p.next = append(p.next, resolve)
	}
	return p
}

func (p *promise) Catch(reject Reject) {
	if !p.catch {
		p.catch = true
		go func() {
			if <-p.errCh {
				reject(p.err)
				p.done <- time.Now()
			}
		}()
	}
}

func New(fn State) *promise {
	var p = &promise{}

	p.resultCh = make(chan bool, 1)
	p.errCh = make(chan bool, 1)
	p.done = make(chan time.Time, 1)
	p.resolve = func(result Result) {
		if atomic.AddInt32(&p.count, 1) == 1 {
			p.result = result
			p.resultCh <- true
			p.errCh <- false
		}
	}
	p.reject = func(err Error) {
		if atomic.AddInt32(&p.count, 1) == 1 {
			p.err = err
			p.errCh <- true
			p.resultCh <- false
		}
	}
	fn(p.resolve, p.reject)
	return p
}

type promiseAll struct {
	results []Result
	err     Error

	resultCh chan bool
	errCh    chan bool
	next     []func(results []Result)

	then  bool
	catch bool
}

func (p *promiseAll) Then(fn func(results []Result)) *promiseAll {
	if !p.then {
		p.then = true
		go func() {
			if <-p.resultCh {
				fn(p.results)
				p.results = nil
				for i := 0; i < len(p.next); i++ {
					p.next[i](p.results)
				}
			}
		}()
	} else {
		p.next = append(p.next, fn)
	}
	return p
}

func (p *promiseAll) Catch(fn func(err Error)) {
	if !p.catch {
		p.catch = true
		go func() {
			if <-p.errCh {
				fn(p.err)
			}
		}()
	}
}

func All(promises ...*promise) *promiseAll {

	var p = &promiseAll{}

	p.resultCh = make(chan bool, 1)
	p.errCh = make(chan bool, 1)

	p.results = make([]Result, len(promises))
	p.err = nil

	var sucCh = make(chan struct{}, 1)
	var errCh = make(chan struct{}, 1)

	var counter = 0

	go func() {
		for i := 0; i < len(promises); i++ {
			var index = i
			promises[index].Then(func(result Result) {
				p.results[index] = result
				sucCh <- struct{}{}
			}).Catch(func(e Error) {
				p.err = e
				errCh <- struct{}{}
			})
		}
	}()

	go func() {
		for {
			select {
			case <-sucCh:
				counter++
				if counter == len(promises) {
					p.resultCh <- true
					p.errCh <- false
				}
			case <-errCh:
				p.errCh <- true
				p.resultCh <- false
			}
		}
	}()

	return p
}
