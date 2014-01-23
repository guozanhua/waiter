package main

import (
	"time"
)

func intermission() {
	// notify clients
	clients.send(true, 1, N_TIMELEFT, 0)

	// start 5 second timer
	end := time.After(5 * time.Second)

	// display some top stats

	// wait for timer to finish
	<-end

	// start new 10 minutes timer
	state.TimeLeft = TEN_MINUTES
	go countDown()

	// change map
	state.changeMap(mapRotation.nextMap(state.Map))
}
