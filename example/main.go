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
	"os"
	"os/signal"
	"time"

	"github.com/lemoyxk/promise"
)

func main() {

	promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		log.Println("start1")    // sync
		resolve("hello world!1") // async
		log.Println("start3")    // sync
	}).Then(func(result promise.Result) {
		log.Println(result)
	}).Catch(func(e promise.Error) {
		log.Println(e)
	})

	promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		log.Println("start2")
		resolve("hello world!2")
	}).Then(func(result promise.Result) {
		log.Println(result)
	}).Catch(func(e promise.Error) {
		log.Println(e)
	})

	log.Println("end")

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
			time.Sleep(time.Millisecond * 50)
			reject("e1 reject")
			resolve("e1 resolve")
		}()
	})

	var e2 = promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		go func() {
			time.Sleep(time.Millisecond * 100)
			reject("e2 reject")
		}()
	})

	promise.All(e1, e2).Then(func(results []promise.Result) {
		log.Println(results)
	}).Catch(func(err promise.Error) {
		log.Println(err)
	})

	var r1 = promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		go func() {
			time.Sleep(time.Millisecond)
			resolve(1)
		}()
	})

	var r2 = promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		go func() {
			time.Sleep(time.Millisecond * 2)
			resolve(2)
		}()
	})

	var r3 = promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		go func() {
			time.Sleep(time.Millisecond * 3)
			resolve(3)
		}()
	})

	promise.Race(r1, r2, r3).Then(func(result promise.Result) {
		log.Println(result)
	})

	signalChan := make(chan os.Signal, 1)
	// 通知
	signal.Notify(signalChan, os.Kill)
	<-signalChan
}
