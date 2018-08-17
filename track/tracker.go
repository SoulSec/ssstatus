package track

import "time"

type TimeTracker struct {
	delayes  Delayes
	nextTime time.Time
	counter  int
}

func (t *TimeTracker) IsReady() bool {
	if time.Now().After(t.nextTime) {
		return true
	}
	return false
}

func (t *TimeTracker) SetNext() (time.Duration, time.Time) {
	t.counter++
	nextDelay := t.delayes.Delay()
	t.nextTime = time.Now().Add(nextDelay)
	return nextDelay, t.nextTime
}

func NewTracker(delayes Delayes) *TimeTracker {
	return &TimeTracker{delayes: delayes}
}

func (t *TimeTracker) HasBeenRan() bool {
	if t.counter > 0 {
		return true
	}
	return false
}
