package simulated

import (
	"fmt"
	"sync"
	"time"
)

func ExampleClock_After() {
	start := time.Date(2021, time.December, 9, 9, 47, 0, 0, time.UTC)
	c := NewClock(start)
	defer c.Stop()

	// get a channel that waits for 5 minutes of simulation time
	after := c.After(5 * time.Minute)

	// drive the simulation until there aren't any pending channels
	var next time.Duration
	for {
		var ok bool
		fmt.Println("Advancing", next)
		next, ok = c.Advance(next)

		if !ok {
			break
		}
	}

	// get the value from the channel buffer
	t := <-after
	fmt.Println("The time is now", t)

	// Output:
	// Advancing 0s
	// Advancing 5m0s
	// The time is now 2021-12-09 09:52:00 +0000 UTC
}

func ExampleClock_Every() {
	start := time.Date(2021, time.December, 9, 9, 47, 0, 0, time.UTC)
	c := NewClock(start)
	defer c.Stop()

	every := c.Every(time.Minute)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer every.Stop()

		for i := 0; i < 5; i++ {
			t := <-every.C()
			fmt.Println("The time is now", t)
		}
	}()

	// drive the simulation until there aren't any pending channels
	var next time.Duration
	for {
		var ok bool
		next, ok = c.Advance(next)

		if !ok {
			break
		}
	}

	// allow the goroutine time to print
	wg.Wait()

	// Output:
	// The time is now 2021-12-09 09:48:00 +0000 UTC
	// The time is now 2021-12-09 09:49:00 +0000 UTC
	// The time is now 2021-12-09 09:50:00 +0000 UTC
	// The time is now 2021-12-09 09:51:00 +0000 UTC
	// The time is now 2021-12-09 09:52:00 +0000 UTC
}
