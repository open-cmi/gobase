package wsem

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type Semaphore struct {
	w  *semaphore.Weighted
	wg sync.WaitGroup
}

func New(size int) *Semaphore {
	return &Semaphore{
		w:  semaphore.NewWeighted(int64(size)),
		wg: sync.WaitGroup{},
	}
}

func (s *Semaphore) Acquire(ctx context.Context) error {

	err := s.w.Acquire(ctx, 1)
	if err != nil {
		return err
	}
	s.wg.Add(1)
	return nil
}

func (s *Semaphore) Release() {
	s.wg.Done()
	s.w.Release(1)
}

func (s *Semaphore) Wait() {
	s.wg.Wait()
}
