package main

import (
	"sync"
	"time"
)

func NewRunSummary() *RunSummary {
	return &RunSummary{}
}
func FinishSummary(s *RunSummary) {
	s.Stop()
	log.Infof("scanned:%d errors:%d duration_sec:%d changed:%d",
		s.Scanned, s.Errors, s.DurationSec, s.Changed)
}

type RunSummary struct {
	mx sync.Mutex

	Scanned     int // how many repos were scanned
	Errors      int // how many repos had errors
	DurationSec int // total time spent
	Changed     int // how many repositories were changed on git pull

	started time.Time
	stopped time.Time
}

func (me *RunSummary) Start() {
	me.started = time.Now()
}

func (me *RunSummary) Stop() {
	me.stopped = time.Now()
	me.DurationSec = int(me.stopped.Sub(me.started).Seconds())
}

func (me *RunSummary) Do(f func(sum *RunSummary)) {
	me.mx.Lock()
	defer me.mx.Unlock()
	f(me)
}

func (me *RunSummary) IncrScanned() {
	me.Do(func(sum *RunSummary) {
		sum.Scanned++
	})
}

func (me *RunSummary) IncrErrors() {
	me.Do(func(sum *RunSummary) {
		sum.Errors++
	})
}
