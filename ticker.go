package main

import (
	"time"
)

var (
	pauseChannel     = make(chan bool)
	interruptChannel = make(chan bool)
)

const (
	TEN_MINUTES int32 = 15000 // 15 seconds for testing and debugging purposes
)

func broadcastPackets() {
	worldStateTicker := time.NewTicker(33 * time.Millisecond)
	for {
		<-worldStateTicker.C
		go sendPositions()
		go sendNetworkMessages()
	}
}

func countDown() {
	endTimer := time.NewTimer(time.Duration(state.TimeLeft) * time.Millisecond)
	gameTicker := time.NewTicker(1 * time.Millisecond)
	paused := false

	for {
		select {
		case <-gameTicker.C:
			state.TimeLeft--

		case shouldPause := <-pauseChannel:
			if shouldPause && !paused {
				endTimer.Stop()
				gameTicker.Stop()
				paused = true
			} else if !shouldPause && paused {
				endTimer.Reset(time.Duration(state.TimeLeft) * time.Millisecond)
				gameTicker = time.NewTicker(1 * time.Millisecond)
				paused = false
			}

		case <-interruptChannel:
			endTimer.Stop()
			gameTicker.Stop()

		case <-endTimer.C:
			endTimer.Stop()
			gameTicker.Stop()
			go intermission()
			return
		}
	}
}
