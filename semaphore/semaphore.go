package semaphore

import (
	"io"
	"fmt"
	"encoding/json"
	"bytes"
)

type Semaphore struct {
	Index uint64   `json:"-"`
	Count int      `json:"count"`
	Max   int      `json:"max"`
	Holders []string  `json:"holders"`
}

func (s Semaphore) String() string {
	return fmt.Sprintf("Semaphore - index: %d, count: %d, max %d, holders: %s", s.Index, s.Count, s.Max, s.Holders)

}

func Load(in io.Reader) (*Semaphore, error) {
	s := new(Semaphore)
	d := json.NewDecoder(in)
	err := d.Decode(s)
	return s, err
}

func (s *Semaphore)ToJSON() (string, error) {
	buf := new(bytes.Buffer)
	e := json.NewEncoder(buf)
	e.SetIndent("","")
	err := e.Encode(s)
	return buf.String(), err
}
