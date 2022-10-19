package main

import (
	"sync"
	"time"
)

type RunSummary struct {
	mx          sync.Mutex
	Scanned     int
	Errors      int
	DurationSec int

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
