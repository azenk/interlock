package semaphore

import (
	"github.com/hashicorp/consul/api"
	"bytes"
)

type ConsulClient interface {
	KV() ConsulKVAPI
}

type ConsulKVAPI interface {
	CAS(*api.KVPair, *api.WriteOptions) (bool, *api.WriteMeta, error)
	Get(key string, q *api.QueryOptions) (*api.KVPair, *api.QueryMeta, error)
}
type SemaphoreConsul struct {
	client ConsulClient
	kv ConsulKVAPI
	key string
	max int
}

func NewSemaphoreConsul(client ConsulClient, key string, max int) *SemaphoreConsul {
	s := new(SemaphoreConsul)
	s.client = client
	s.kv = client.KV()
	s.key = key
	s.max = max
	return s
}

func (s *SemaphoreConsul) operation(op func (*SemaphoreData,string) (bool,error), id string) (bool,error) {
	kvapi := s.kv
	query_opts := new(api.QueryOptions)
	// Get current value
	kvdata, _, err := kvapi.Get(s.key, query_opts)
	if err != nil {
		return false, err
	}

	var sem_data *SemaphoreData

	if kvdata == nil {
		kvdata = new(api.KVPair)
		kvdata.Key = s.key
		sem_data = New(s.max)
	} else {
		sem_data, err = Load(bytes.NewBuffer(kvdata.Value))
		if err != nil {
			return false, err
		}
	}

	result, err := op(sem_data, id)
	if !result || err != nil {
		return false, err
	}

	val, err := sem_data.ToJSON()
	if err != nil {
		return false, err
	}
	kvdata.Value = []byte(val)

	wropt := new(api.WriteOptions)
	success, _, err := kvapi.CAS(kvdata, wropt)
	return success, err
}

func (s *SemaphoreConsul) Acquire(id string) (bool,error) {
	return s.operation((*SemaphoreData).Acquire, id)
}

func (s *SemaphoreConsul) Release(id string) (bool,error) {
	return s.operation((*SemaphoreData).Release, id)
}

func (s *SemaphoreConsul) Holds(id string) (bool,error) {
	return s.operation((*SemaphoreData).Holds, id)
}
