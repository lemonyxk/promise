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

import (
	"sync/atomic"
)

type Results[T any, P any] interface {
	Then(func(result []T)) Catch[P]
}

type results[T any, P any] struct {
	fn func(index int)

	index  int32
	chList []chan []T
	ehList []chan P
}

func (p *results[T, P]) Then(then func([]T)) Catch[P] {

	var index = atomic.AddInt32(&p.index, 1)
	// not first
	if index > 0 {
		p.chList = append(p.chList, make(chan []T, 1))
		p.ehList = append(p.ehList, make(chan P, 1))
		p.fn(int(index))
	}

	select {
	case res := <-p.chList[index]:
		var c = &catch[P]{err: empty[P](), run: false}
		then(res)
		return c
	case err := <-p.ehList[index]:
		var c = &catch[P]{err: err, run: true}
		return c
	}
}

type Result[T any, P any] interface {
	Then(func(result T)) Catch[P]
}

type result[T any, P any] struct {
	fn func(index int)

	index  int32
	chList []chan T
	ehList []chan P
}

func (p *result[T, P]) Then(then func(T)) Catch[P] {

	var index = atomic.AddInt32(&p.index, 1)
	// not first
	if index > 0 {
		p.chList = append(p.chList, make(chan T, 1))
		p.ehList = append(p.ehList, make(chan P, 1))
		p.fn(int(index))
	}

	select {
	case res := <-p.chList[index]:
		var c = &catch[P]{err: empty[P](), run: false}
		then(res)
		return c
	case err := <-p.ehList[index]:
		var c = &catch[P]{err: err, run: true}
		return c
	}
}
