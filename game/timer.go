package game

import "time"

type Timer struct {
	startedAt time.Time
	duration  int64
}

func (timer *Timer) start(duration int64) {
	timer.startedAt = time.Now()
	timer.duration = duration
}

func (timer *Timer) restart() {
	timer.startedAt = time.Now()
}

func (timer *Timer) ended() bool {
	return time.Since(timer.startedAt).Milliseconds() >= timer.duration
}
