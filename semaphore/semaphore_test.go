package semaphore

import (
	"testing"
	"strings"
)

func TestSuccessfullLoad(t *testing.T) {
	json_sem := "{ \"index\": 1234, \"count\": 1, \"max\": 3, \"holders\": [\"a\"] }"
	jr := strings.NewReader(json_sem)
	s, err := Load(jr)
	if err != nil || s.Index != 0 || s.Count != 1 || s.Max != 3 {
		t.Fail()
	}
}

func TestFailedLoad(t *testing.T) {
	json_sem := "{ count\": 1, \"max\": 3, \"holders\": [\"a\"] }"
	jr := strings.NewReader(json_sem)
	_, err := Load(jr)
	if err == nil {
		t.Fail()
	}
}
