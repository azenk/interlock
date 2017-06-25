package semaphore

import (
	"io"
	"fmt"
	"encoding/json"
	"bytes"
	"time"
)

type Semaphore struct {
	Index uint64   `json:"-"`
	Max   int      `json:"max"`
	Holders map[string]int64  `json:"holders"`
}

func New(max int) *Semaphore {
	s := new(Semaphore)
	s.Max = max
	s.Holders = make(map[string]int64, max)
	return s
}

func (s Semaphore) String() string {
	return fmt.Sprintf("Semaphore - index: %d, count: %d, max %d, holders: %s", s.Index, s.Count(), s.Max, s.Holders)

}

func Load(in io.Reader) (*Semaphore, error) {
	s := new(Semaphore)
	d := json.NewDecoder(in)
	err := d.Decode(s)
	return s, err
}

func (s *Semaphore)Count() int {
	return len(s.Holders)
}

func (s *Semaphore)ToJSON() (string, error) {
	buf := new(bytes.Buffer)
	e := json.NewEncoder(buf)
	e.SetIndent("","")
	err := e.Encode(s)
	return buf.String(), err
}

func (s *Semaphore)Acquire(id string) bool {
	_, ok := s.Holders[id]
	if ok {
		return true
	}

	if s.Count() == s.Max {
		return false
	}

	s.Holders[id] = time.Now().Unix()
	return true
}

// Return true if id in list of holders, false otherwise
func (s *Semaphore)Holds(id string) bool {
	_,ok := s.Holders[id]
	return ok
}

// Remove holder entry from semaphore if it's present
func (s *Semaphore)Release(id string) bool {
	delete(s.Holders, id)
	return !s.Holds(id)
}
