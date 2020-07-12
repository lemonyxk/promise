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

type Resolve func(Result)
type Reject func(Error)

type Result interface{}
type Error interface{}

type State func(resolve Resolve, reject Reject)
