package gutils

import "time"

//TimeFunction counts the duration of a function, and return the elapsed time in ms
func TimeFunction(toRun func()) int64 {
	start := time.Nanosecond.Milliseconds()
	toRun()
	return time.Nanosecond.Milliseconds() - start
}


