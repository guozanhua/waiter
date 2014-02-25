package main

import (
	"./enet"
	"time"
)

var (
	// channel to pause the game
	pauseChannel = make(chan bool)

	// channel to interrupt the game (for example when a master changes mode or map mid-game)
	interruptChannel = make(chan bool)
)

const (
	TEN_MINUTES int32 = 180000 // 3 minutes for testing and debugging purposes
)

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

func intermission() {
	// notify all clients
	clients.send(enet.PACKET_FLAG_RELIABLE, 1, N_TIMELEFT, 0)

	// start 5 second timer
	end := time.After(5 * time.Second)

	// TODO: send server messages with some top stats

	// wait for timer to finish
	<-end

	// start new 10 minutes timer
	state.TimeLeft = TEN_MINUTES
	go countDown()

	// change map
	state.changeMap(mapRotation.nextMap(state.Map))
}
