package main

import (
	"fmt"
	"time"
)

// countdown runs for totalDur, emiting a tick for every tickInterval.
// On each tick, we receive the number of seconds that elapsed since the start of the countdown.
// After totalDur expires, countdown returns.
func countdown(totalDur time.Duration, tickInterval time.Duration) {
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	// Create a Timer that fires once after totalDur
	timer := time.NewTimer(totalDur)
	defer timer.Stop()

	// elapsedSeconds tracks how many seconds have elapsed
	elapsedSeconds := 0

	// loop until the timer expires
	for {
		select {
		case <-ticker.C:
			// One tick elapsed
			elapsedSeconds++
			fmt.Printf("\rElapsed: %2d seconds elapsed\n", elapsedSeconds)
		case <-timer.C:
			// Total duration reached
			fmt.Printf("\rElapsed: %2d sec - DONE!\n", elapsedSeconds)
			return
		}
	}
}

func main() {
	countdown(10*time.Second, 1*time.Second)
}
