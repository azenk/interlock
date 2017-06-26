package semaphore

import (
  "testing"
  "github.com/hashicorp/consul/api"
  "errors"
  "fmt"
)

type mockClient struct {}
type mockKVAPI struct {
  values map[string]api.KVPair
  index uint64
}

func (c *mockClient) KV() ConsulKVAPI {
  kv := new(mockKVAPI)
  kv.values = make(map[string]api.KVPair)
  return kv
}

func (kv *mockKVAPI) Get(key string, wropt *api.QueryOptions) (*api.KVPair, *api.QueryMeta, error) {
  var err error
  data, ok := kv.values[key]
  if !ok {
    return nil, new(api.QueryMeta), err
  }
  return &data, new(api.QueryMeta), err
}

func (kv *mockKVAPI) CAS(data *api.KVPair, wropt *api.WriteOptions) (bool, *api.WriteMeta, error) {
  var err error
  var success bool
  cur_value, ok := kv.values[data.Key]
  if !ok && data.ModifyIndex == 0 {
    data.CreateIndex = kv.index
    kv.index += 1
    kv.values[data.Key] = *data
    success = true
  } else if ok && data.ModifyIndex == cur_value.ModifyIndex {
    data.ModifyIndex = kv.index
    kv.index += 1
    kv.values[data.Key] = *data
    success = true
  } else {
    err = errors.New(fmt.Sprintf(
      "Unable to check and set value: ok -> %b, new_value: %s, cur_value: %s",
      ok, data, cur_value))
  }
  return success, new(api.WriteMeta), err
}

func TestConsulAcquire(t *testing.T) {
  s := NewSemaphoreConsul(new(mockClient), "/test/key", 1)
  ok, err := s.Acquire("test")
	if !ok {
		t.Error("Failed to acquire semaphore")
		t.Error(err)
	}
	ok, err = s.Acquire("test2")
	if ok {
		t.Error("Able to acquire semaphore more than once")
		t.Error(err)
	}
}
