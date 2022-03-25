package main

import "sync"

type RunSummary struct {
	mx      sync.Mutex
	Scanned int
	Updated int
	Errors  int
}

func (me *RunSummary) Do(f func(sum *RunSummary)) {
	me.mx.Lock()
	defer me.mx.Unlock()
	f(me)
}

func (me *RunSummary) IncrScanned() {
	me.Do(func(sum *RunSummary) {
		me.Scanned++
	})
}

func (me *RunSummary) IncrUpdated() {
	me.Do(func(sum *RunSummary) {
		me.Updated++
	})
}

func (me *RunSummary) IncrErrors() {
	me.Do(func(sum *RunSummary) {
		me.Errors++
	})
}
