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

	"github.com/lemonyxk/promise"
)

func main() {

	log.Println("start")

	var r1 = promise.New(func(resolve func(int), reject func(error)) {
		go func() {
			time.Sleep(time.Millisecond * 300)
			resolve(1)
		}()
	})

	var r2 = promise.New(func(resolve func(int), reject func(error)) {
		go func() {
			time.Sleep(time.Millisecond * 200)
			resolve(2)
		}()
	})

	var r3 = promise.New(func(resolve func(int), reject func(error)) {
		go func() {
			time.Sleep(time.Millisecond * 100)
			resolve(3)
		}()
	})

	promise.Race(r1, r2, r3).Then(func(result int) {
		log.Println(result)
	}).Catch(func(err error) {
		log.Println(err)
	})

	promise.All(r1, r2, r3).Then(func(result []int) {
		log.Println(result)
	}).Catch(func(err error) {
		log.Println(err)
	})

	promise.Fall(r1, r2, r3).Then(func(result []int) {
		log.Println(result)
	}).Catch(func(err error) {
		log.Println(err)
	})

	log.Println("end")

	// 2020/07/13 02:00:15 start
	// 2020/07/13 02:00:16 3
	// 2020/07/13 02:00:19 [1 2 3]
	// 2020/07/13 02:00:25 [1 2 3]
	// 2020/07/13 02:00:25 end

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Kill)
	// <-signalChan
}
