/**
* @program: promise
*
* @description:
*
* @author: lemo
*
* @create: 2022-05-17 16:25
**/

package promise

import (
	"errors"
	"testing"
	"time"
)

func assert(t *testing.T, condition bool, msg string, v ...interface{}) {
	if !condition {
		t.Errorf(msg, v...)
	}
}

func TestPromise(t *testing.T) {
	var p = New(func(resolve func(int), reject func(error)) {
		resolve(1)
		reject(errors.New("error"))
	})

	p.Then(func(v int) {
		assert(t, v == 1, "then value is %d", v)
	})

	p.Catch(func(err error) {
		assert(t, err.Error() == "error", "catch error is %s", err.Error())
	})

	var p1 = New(func(resolve func(int), reject func(error)) {
		reject(errors.New("error"))
		resolve(1)
	})

	p1.Then(func(v int) {
		assert(t, v == 1, "then value is %d", v)
	})

	p1.Catch(func(err error) {
		assert(t, err.Error() == "error", "catch error is %s", err.Error())
	})

	var now = time.Now()

	var p2 = New(func(resolve func(int), reject func(error)) {
		time.AfterFunc(time.Second*1, func() {
			resolve(1)
		})
	})

	p2.Then(func(v int) {
		var sub = time.Now().Sub(now).Seconds()
		assert(t, int(sub) == 1, "then time sub is %d", int(sub))
		assert(t, v == 1, "then value is %d", v)
	})
}

func TestResolve(t *testing.T) {
	var p = Resolve(1)

	p.Then(func(v int) {
		assert(t, v == 1, "then value is %d", v)
	})
}

func TestReject(t *testing.T) {
	var p = Reject[int](errors.New("error"))

	p.Catch(func(err error) {
		assert(t, err.Error() == "error", "catch error is %s", err.Error())
	})
}

func TestAll(t *testing.T) {
	var now = time.Now()

	var p1 = New(func(resolve func(int), reject func(error)) {
		go func() {
			time.Sleep(time.Second * 1)
			resolve(1)
		}()
	})
	var p2 = New(func(resolve func(int), reject func(error)) {
		go func() {
			time.Sleep(time.Second * 1)
			resolve(2)
		}()
	})
	var p3 = New(func(resolve func(int), reject func(error)) {
		go func() {
			time.Sleep(time.Second * 1)
			resolve(3)
		}()
	})

	var p = All(p1, p2, p3)

	p.Then(func(v []int) {
		var sub = time.Now().Sub(now).Seconds()
		assert(t, int(sub) == 1, "then time sub is %d", int(sub))
		assert(t, v[0] == 1, "then value is %d", v[0])
		assert(t, v[1] == 2, "then value is %d", v[1])
		assert(t, v[2] == 3, "then value is %d", v[2])
	})
}

func TestFall(t *testing.T) {
	var now = time.Now()

	var p1 = New(func(resolve func(int), reject func(error)) {
		time.Sleep(time.Second * 1)
		resolve(1)
	})
	var p2 = New(func(resolve func(int), reject func(error)) {
		time.Sleep(time.Second * 1)
		resolve(2)
	})
	var p3 = New(func(resolve func(int), reject func(error)) {
		time.Sleep(time.Second * 1)
		resolve(3)
	})

	var p = Fall(p1, p2, p3)

	p.Then(func(v []int) {
		var sub = time.Now().Sub(now).Seconds()
		assert(t, int(sub) == 3, "then time sub is %d", int(sub))
		assert(t, v[0] == 1, "then value is %d", v[0])
		assert(t, v[1] == 2, "then value is %d", v[1])
		assert(t, v[2] == 3, "then value is %d", v[2])
	})
}

func TestRace(t *testing.T) {
	var p1 = New(func(resolve func(int), reject func(error)) {
		go func() {
			time.Sleep(time.Millisecond * 3)
			resolve(1)
		}()
	})
	var p2 = New(func(resolve func(int), reject func(error)) {
		go func() {
			time.Sleep(time.Millisecond * 1)
			resolve(2)
		}()
	})
	var p3 = New(func(resolve func(int), reject func(error)) {
		go func() {
			time.Sleep(time.Millisecond * 2)
			resolve(3)
		}()
	})

	var p = Race(p1, p2, p3)

	p.Then(func(v int) {
		assert(t, v == 2, "then value is %d", v)
	})
}

func TestPromise_Finally(t *testing.T) {
	var run = false
	Reject[int](errors.New("run")).Finally(func() {
		run = true
	})

	assert(t, run, "finally run")
}
