package gutils

import "time"

//DoEvery Execute function at a set interval from its own thread
func DoEvery(d time.Duration, f func()) {
	go func() {
		for range time.Tick(d) {
			f()
		}
	}()
}
