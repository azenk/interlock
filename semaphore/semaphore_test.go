package semaphore

import (
	"testing"
	"strings"
	"encoding/json"
	"reflect"
)

func TestSuccessfullLoad(t *testing.T) {
	json_sem := "{ \"index\": 1234, \"count\": 1, \"max\": 3, \"holders\": {\"a\": 0 }}"
	jr := strings.NewReader(json_sem)
	s, err := Load(jr)
	if err != nil || s.Index != 0 || s.Count != 1 || s.Max != 3 {
		t.Fail()
	}
}

func TestFailedLoad(t *testing.T) {
	json_sem := "{ count\": 1, \"max\": 3, \"holders\": {\"a\": 0}}"
	jr := strings.NewReader(json_sem)
	_, err := Load(jr)
	if err == nil {
		t.Fail()
	}
}

func TestToJSON(t *testing.T) {
	s := Semaphore{1234,1,3, map[string]uint64{"a": 0}}
	json_repr,err := s.ToJSON()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	ref_json_repr := "{\"count\":1,\"max\":3,\"holders\":{\"a\":0}}"

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
