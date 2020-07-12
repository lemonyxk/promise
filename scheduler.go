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

type scheduler struct {
	queue []*promise
	mux   sync.Mutex
}

func (s *scheduler) add(p *promise) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.queue = append(s.queue, p)
}
