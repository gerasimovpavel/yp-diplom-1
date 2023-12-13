package scheduler

import "time"

func Schedule(f func(), interval time.Duration, done <-chan bool) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				f()
			case <-done:
				return
			}
		}
	}()
	return ticker
}
