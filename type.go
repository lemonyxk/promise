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

type Resolve[T any] func(T)
type Reject[T any] func(T)

// type Result[T any] any
// type Error[T any] any

type State[T any, P any] func(resolve Resolve[T], reject Reject[P])
