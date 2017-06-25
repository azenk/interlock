package semaphore

import (
	"io"
)

type Semaphore struct {
	Index uint64
	Count int
	Max   int
	Holders []string
}

func Load(in io.Reader) *Semaphore {
	return &Semaphore{0, 0, 0, make([]string, 0)}
}
