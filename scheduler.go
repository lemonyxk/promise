/**
* @program: promise
*
* @description:
*
* @author: lemo
*
* @create: 2020-07-12 16:07
**/

package promise

import "sync"

type scheduler[T any, P any] struct {
	queue []*promise[T, P]
	mux   sync.Mutex
}

func (s *scheduler[T, P]) add(p *promise[T, P]) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.queue = append(s.queue, p)
}
