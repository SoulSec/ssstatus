package track

import "time"

type Delayes interface {
	Delay() time.Duration
}
