package main

import (
	"time"
)

// global channels used to notify the server of intermission, or the ticker of an interrupting event (like a map change mid-game)
var (
	interruptTickerChannel = make(chan bool)
	intermissionChannel    = make(chan bool)
)

const (
	TEN_MINUTES time.Duration = 10 * time.Minute
)

// interrupt is for map changes in the middle of the game, etc.
// intermission is the channel to notify the main loop of intermission
func countDown(d time.Duration, interrupt chan bool, intermission chan bool) {
	state.TimeLeft = int32(d.Nanoseconds() / 1000000)
	end := time.After(d)
	t := time.NewTicker(1 * time.Millisecond)

	for {
		select {
		case <-t.C:
			state.TimeLeft--
		case <-interrupt:
			t.Stop()
			return
		case <-end:
			intermission <- true
		}
	}
}
