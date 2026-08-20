package scheduler

import (
	"context"
	"sync"
	"time"
)

type Job struct {
	Name string
	At   time.Time
	Run  func(context.Context) error
}
type Scheduler struct {
	mu   sync.Mutex
	jobs []Job
	stop chan struct{}
}

func New() *Scheduler             { return &Scheduler{stop: make(chan struct{})} }
func (s *Scheduler) Add(j Job)    { s.mu.Lock(); defer s.mu.Unlock(); s.jobs = append(s.jobs, j) }
func (s *Scheduler) Pending() int { s.mu.Lock(); defer s.mu.Unlock(); return len(s.jobs) }
func (s *Scheduler) RunOnce(ctx context.Context, now time.Time) []error {
	s.mu.Lock()
	ready := []Job{}
	keep := []Job{}
	for _, j := range s.jobs {
		if !j.At.After(now) {
			ready = append(ready, j)
		} else {
			keep = append(keep, j)
		}
	}
	s.jobs = keep
	s.mu.Unlock()
	errs := []error{}
	for _, j := range ready {
		if j.Run != nil {
			if e := j.Run(ctx); e != nil {
				errs = append(errs, e)
			}
		}
	}
	return errs
}
