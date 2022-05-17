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

type Promise[T any] interface {
	Then(func(T)) Promise[T]
	Catch(func(error)) Promise[T]
	Finally(func())
}

type promise[T any] struct {
	fn func()

	ch chan T
	eh chan error

	done chan bool
}

func (p *promise[T]) Then(then func(T)) Promise[T] {
	var done = <-p.done
	if done {
		var res = <-p.ch
		then(res)
		p.ch <- res
		p.done <- done
		return p
	} else {
		p.done <- done
		return p
	}
}

func (p *promise[T]) Catch(catch func(error)) Promise[T] {

	var done = <-p.done
	if done {
		p.done <- done
		return p
	} else {
		var err = <-p.eh
		catch(err)
		p.eh <- err
		p.done <- done
		return p
	}
}

func (p *promise[T]) Finally(finally func()) {
	finally()
}

func New[T any](state func(resolve func(T), reject func(error))) Promise[T] {

	var p = new(promise[T])
	p.ch = make(chan T, 1)
	p.eh = make(chan error, 1)
	p.done = make(chan bool, 1)

	p.fn = func() {
		var counter int32 = 0
		// just one can be exec
		var resolve = func(result T) {
			if atomic.AddInt32(&counter, 1) == 1 {
				p.ch <- result
				p.done <- true
			}
		}
		var reject = func(err error) {
			if atomic.AddInt32(&counter, 1) == 1 {
				p.eh <- err
				p.done <- false
			}
		}
		state(resolve, reject)
	}

	p.fn()

	return p
}

func Resolve[T any](result T) Promise[T] {
	return New[T](func(resolve func(T), reject func(error)) {
		resolve(result)
	})
}

func Reject[T any](err error) Promise[T] {
	return New[T](func(resolve func(T), reject func(error)) {
		reject(err)
	})
}

func All[T any](promises ...Promise[T]) Promise[[]T] {
	return New[[]T](func(resolve func([]T), reject func(error)) {
		var result = make([]T, len(promises))
		var counter int32 = 0

		for i := 0; i < len(promises); i++ {
			var pi = i
			go func() {
				promises[pi].Then(func(res T) {
					result[pi] = res
					if atomic.AddInt32(&counter, 1) == int32(len(promises)) {
						resolve(result)
					}
				}).Catch(func(err error) {
					reject(err)
				})
			}()
		}
	})
}

func Race[T any](promises ...Promise[T]) Promise[T] {
	return New(func(resolve func(T), reject func(error)) {
		var sucCounter int32 = 0
		var errCounter int32 = 0

		for i := 0; i < len(promises); i++ {
			var pi = i
			go func() {
				promises[pi].Then(func(result T) {
					if atomic.AddInt32(&sucCounter, 1) == 1 {
						resolve(result)
					}
				}).Catch(func(err error) {
					if atomic.AddInt32(&errCounter, 1) == 1 {
						reject(err)
					}
				})
			}()
		}
	})
}

func Fall[T any](promises ...Promise[T]) Promise[[]T] {
	return New[[]T](func(resolve func([]T), reject func(error)) {
		var sucCounter int32 = 0
		var errCounter int32 = 0
		var results = make([]T, len(promises))

		var pi = 0

		var fn func(int)

		fn = func(pi int) {
			promises[pi].Then(func(result T) {
				results[pi] = result
				if atomic.AddInt32(&sucCounter, 1) == int32(len(promises)) {
					resolve(results)
				} else {
					pi++
					fn(pi)
				}
			}).Catch(func(err error) {
				if atomic.AddInt32(&errCounter, 1) == 1 {
					reject(err)
				}
			})
		}

		fn(pi)
	})
}
