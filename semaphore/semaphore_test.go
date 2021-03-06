package semaphore

import (
	"testing"
	"strings"
	"encoding/json"
	"reflect"
)

func TestSuccessfullLoad(t *testing.T) {
	json_sem := "{ \"index\": 1234, \"max\": 3, \"holders\": {\"a\": 0 }}"
	jr := strings.NewReader(json_sem)
	s, err := Load(jr)
	if err != nil || s.Index != 0 || s.Count() != 1 || s.Max != 3 {
		t.Fail()
	}
}

func TestFailedLoad(t *testing.T) {
	json_sem := "{ max\": 3, \"holders\": {\"a\": 0}}"
	jr := strings.NewReader(json_sem)
	_, err := Load(jr)
	if err == nil {
		t.Fail()
	}
}

func TestToJSON(t *testing.T) {
	s := New(3)
	s.Index = 1234
	s.Holders["a"] = 0
	json_repr,err := s.ToJSON()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	ref_json_repr := "{\"max\":3,\"holders\":{\"a\":0}}"

	var o1, o2 interface{}

	err = json.Unmarshal([]byte(json_repr), &o1)
	if err != nil {
		t.Log("Unable to unmarshal encoded json")
		t.Log(err)
		t.Fail()
	}

	err = json.Unmarshal([]byte(ref_json_repr), &o2)
	if err != nil {
		t.Log("Unable to unmarshal encoded json")
		t.Log(err)
		t.Fail()
	}

	if !reflect.DeepEqual(o1,o2) {
		t.Errorf("%s != %s",o1,o2)
	}
}

func TestAcquire(t *testing.T) {
	s := New(1)
	ok,_ := s.Acquire("Host1")
	if !ok {
		t.Errorf("Unable to acquire semaphore from Host1")
	}
	ok,_ = s.Acquire("Host2")
	if ok {
		t.Errorf("Incorrectly Able to acquire semaphore from Host2")
	}
}

func TestHolds(t *testing.T) {
	s := New(1)
	ok,_ := s.Acquire("Host1")
	if !ok {
		t.Errorf("Unable to acquire semaphore from Host1")
	}

	ok, _ = s.Holds("Host1")
	if !ok {
		t.Errorf("Semaphore doesn't claim that Host1 is a holder")
	}

	ok, _ = s.Holds("Host2")
	if ok {
		t.Errorf("Semaphore claims that Host2 is a holder")
	}
}

func TestRelease(t *testing.T) {
	s := New(2)
	s.Holders["test"] = 0

	ok, _ := s.Release("test")
	if !ok {
		t.Errorf("Failed to release semaphore")
	}

	ok, _ = s.Holds("test")
	if ok {
		t.Errorf("Failed to release semaphore")
	}

	ok, _ = s.Release("test")
	if !ok {
		t.Errorf("Repeated release of semaphore failed")
	}
}
