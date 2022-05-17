```go
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
"errors"
"log"
"time"

	"github.com/lemonyxk/promise"
)

func main() {

	log.Println("start")

	var r1 = promise.New(func(resolve func(int), reject func(error)) {
		go func() {
			log.Println("r1 start")
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

	promise.New(func(resolve func(string), reject func(error)) {
		go func() {
			time.Sleep(time.Millisecond * 100)
			reject(errors.New("error"))
			resolve("test1") // not execute
		}()
	}).Catch(func(err error) {
		log.Println("reject:", err)
	}).Then(func(result string) {
		log.Println("resolve:", result)
	}).Finally(func() {
		log.Println("finally")
	})

	log.Println("end")

	promise.Resolve("test").Then(func(result string) {
		log.Println("resolve:", result)
	}).Catch(func(err error) {
		log.Println("reject:", err)
	}).Finally(func() {
		log.Println("finally")
	})

	promise.Reject[string](errors.New("error")).Then(func(result string) {
		log.Println("resolve:", result)
	}).Catch(func(err error) {
		log.Println("reject:", err)
	}).Finally(func() {
		log.Println("finally")
	})

	// 2022/05/16 18:45:24 start
	// 2022/05/16 18:46:00 r1 start
	// 2022/05/16 18:45:24 3
	// 2022/05/16 18:45:24 [1 2 3]
	// 2022/05/16 18:45:24 [1 2 3]
	// 2022/05/16 18:45:24 reject: test
	// 2022/05/16 18:45:24 finally
	// 2022/05/16 18:45:24 end

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Kill)
	// <-signalChan
}

```