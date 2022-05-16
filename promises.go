/**
* @program: promise
*
* @description:
*
* @author: lemo
*
* @create: 2022-05-15 18:07
**/

package promise

type Promises[T any, P any] interface {
	Then(func([]T)) Promises[T, P]
	Catch(func(P)) Promises[T, P]
	Finally(func())
}

type promises[T any, P any] struct {
	fn func()

	ch chan []T
	eh chan P

	done chan bool
}

func (p *promises[T, P]) Then(then func([]T)) Promises[T, P] {
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

func (p *promises[T, P]) Catch(catch func(P)) Promises[T, P] {
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

func (p *promises[T, P]) Finally(finally func()) {
	finally()
}
