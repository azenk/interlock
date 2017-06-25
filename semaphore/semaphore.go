package semaphore

import (
	"io"
	"fmt"
	"encoding/json"
	"bytes"
	"time"
)

type Semaphore interface {
	Acquire(string) (bool,error)
	Holds(string) (bool,error)
	Release(string) (bool,error)
}

type SemaphoreData struct {
	Index uint64   `json:"-"`
	Max   int      `json:"max"`
	Holders map[string]int64  `json:"holders"`
}

func New(max int) *SemaphoreData {
	s := new(SemaphoreData)
	s.Max = max
	s.Holders = make(map[string]int64, max)
	return s
}

func (s SemaphoreData) String() string {
	return fmt.Sprintf("SemaphoreData - index: %d, count: %d, max %d, holders: %s", s.Index, s.Count(), s.Max, s.Holders)

}

func Load(in io.Reader) (*SemaphoreData, error) {
	s := new(SemaphoreData)
	d := json.NewDecoder(in)
	err := d.Decode(s)
	return s, err
}

func (s *SemaphoreData)Count() int {
	return len(s.Holders)
}

func (s *SemaphoreData)ToJSON() (string, error) {
	buf := new(bytes.Buffer)
	e := json.NewEncoder(buf)
	e.SetIndent("","")
	err := e.Encode(s)
	return buf.String(), err
}

func (s *SemaphoreData)Acquire(id string) (bool,error) {
	_, ok := s.Holders[id]
	if ok {
		return true, nil
	}

	if s.Count() == s.Max {
		return false, nil
	}

	s.Holders[id] = time.Now().Unix()
	return true, nil
}

// Return true if id in list of holders, false otherwise
func (s *SemaphoreData)Holds(id string) (bool,error) {
	_,ok := s.Holders[id]
	return ok, nil
}

// Remove holder entry from semaphore if it's present
func (s *SemaphoreData)Release(id string) (bool,error) {
	delete(s.Holders, id)
	ok,err := s.Holds(id)
	return !ok, err
}
