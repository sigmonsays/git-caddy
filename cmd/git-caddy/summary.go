package main

import "sync"

type RunSummary struct {
	mx      sync.Mutex
	Scanned int
	Errors  int
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
