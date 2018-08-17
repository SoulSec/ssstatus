package track

import "time"

type ExpTrack struct {
	counter int
	base    int
}

func calculateExp(base, counter int) int {
	if counter == 0 {
		return 1
	}
	return base * calculateExp(base, counter-1)
}

func (e *ExpTrack) Delay() time.Duration {
	e.counter++
	return time.Duration(calculateExp(e.base, e.counter)) * time.Second
}

func NewExpTrack(base int) *ExpTrack {
	return &ExpTrack{
		base: base,
	}
}
