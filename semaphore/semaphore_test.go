package semaphore

import (
	"testing"
	"strings"
)

func TestLoad(t *testing.T) {
	json_sem := "{ \"index\": 1234, \"count\": 1, \"max\": 3, \"members\": [\"a\"] }"
	jr := strings.NewReader(json_sem)
	s := Load(jr)
	if s.Index != 1234 || s.Count != 1 || s.Max != 3 {
		t.Fail()
	}
}
