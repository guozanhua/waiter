package main

import (
	"time"
)

var (
	pauseChannel     = make(chan bool)
	interruptChannel = make(chan bool)
)

const (
	TEN_MINUTES int32 = 600000
)

func countDown() {
	endTimer := time.NewTimer(time.Duration(state.TimeLeft) * time.Millisecond)
	gameTicker := time.NewTicker(1 * time.Millisecond)
	worldStateTicker := time.NewTicker(33 * time.Millisecond)
	paused := false

	for {
		select {
		case <-gameTicker.C:
			state.TimeLeft--

		case <-worldStateTicker.C:
			go sendPositions()
			go sendNetworkMessages()

		case shouldPause := <-pauseChannel:
			if shouldPause && !paused {
				endTimer.Stop()
				gameTicker.Stop()
				worldStateTicker.Stop()
				paused = true
			} else if !shouldPause && paused {
				endTimer.Reset(time.Duration(state.TimeLeft) * time.Millisecond)
				gameTicker = time.NewTicker(1 * time.Millisecond)
				worldStateTicker = time.NewTicker(33 * time.Millisecond)
				paused = false
			}

		case <-interruptChannel:
			endTimer.Stop()
			gameTicker.Stop()
			worldStateTicker.Stop()

		case <-endTimer.C:
			endTimer.Stop()
			gameTicker.Stop()
			worldStateTicker.Stop()
			go intermission()
			return
		}
	}
}
