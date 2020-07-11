/**
* @program: promise
*
* @description:
*
* @author: lemo
*
* @create: 2020-07-11 13:20
**/

package main

import (
	"log"
	"time"

	"github.com/Lemo-yxk/promise"
)

func main() {

	promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		go func() {
			time.Sleep(time.Second * 2)
			resolve("hello world!")
			reject("err")
		}()
	}).Then(func(result promise.Result) {
		log.Println(result)
	}).Catch(func(e promise.Error) {
		log.Println(e)
	})

	var p1 = promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		go func() {
			time.Sleep(time.Second * 1)
			resolve(1)
		}()
	})

	var p2 = promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		go resolve(2)
	})

	promise.All(p1, p2).Then(func(results []promise.Result) {
		log.Println(results)
	}).Catch(func(err promise.Error) {
		log.Println(err)
	})

	var e1 = promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		go func() {
			time.Sleep(time.Second * 1)
			resolve("e1 resolve")
		}()
	})

	var e2 = promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		go func() {
			time.Sleep(time.Second * 1)
			reject("e2 reject")
		}()
	})

	promise.All(e1, e2).Then(func(results []promise.Result) {
		log.Println(results)
	}).Catch(func(err promise.Error) {
		log.Println(err)
	})

	time.Sleep(time.Second * 3)
}
