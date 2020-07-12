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

type Promise interface {
	Then(func(result Result)) Catch
}

type Catch interface {
	Catch(func(err Error))
}

type data struct {
	res Result
	err Error
}

type promise struct {
	ch chan data
	fn func()
}

func (p promise) Then(then func(Result)) Catch {
	p.fn()
	r := <-p.ch
	var c = catch{err: r.err}
	if r.err != nil {
		return c
	}
	then(r.res)
	return c
}

type catch struct {
	err Error
}

func (c catch) Catch(catch func(Error)) {
	if c.err != nil {
		catch(c.err)
	}
}

func New(state State) Promise {

	var p promise
	p.ch = make(chan data, 1)

	p.fn = func() {
		var counter int32 = 0
		// just one can be exec
		var resolve = func(result Result) {
			if atomic.AddInt32(&counter, 1) == 1 {
				p.ch <- data{res: result, err: nil}
			}
		}
		var reject = func(err Error) {
			if atomic.AddInt32(&counter, 1) == 1 {
				p.ch <- data{res: nil, err: err}
			}
		}
		state(resolve, reject)
	}

	return p
}
