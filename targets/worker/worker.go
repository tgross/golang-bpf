package main

import (
	"fmt"
	"sync"
)

type State struct {
	Counts int
	lock   sync.Mutex
}

func (s *State) Inc() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Counts++
	return s.Counts
}

type count struct {
	i int
}

type Worker struct {
	state  *State
	counts []*count
}

func newWorker(s *State) *Worker {
	return &Worker{
		state:  s,
		counts: []*count{},
	}
}

func (w *Worker) work(wg *sync.WaitGroup, times int) {
	defer wg.Done()
	for j := 0; j < times; j++ {
		result := w.state.Inc()
		w.counts = append(w.counts, &count{i: result})
	}
}

func (w *Worker) results() []*count {
	return w.counts
}

func main() {
	s := &State{}
	workers := []*Worker{}
	total := 0
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		worker := newWorker(s)
		workers = append(workers, worker)
		go worker.work(&wg, i)
	}
	wg.Wait()

	for _, worker := range workers {
		results := worker.results()
		for _, result := range results {
			total = total + result.i
		}
	}
	fmt.Println(total)
}
