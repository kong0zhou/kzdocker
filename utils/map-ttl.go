package utils

import (
	"kzdocker/log"
	"sync"
	"time"
)

type item struct {
	value     interface{}
	timeToDel int64 //当当前时间大于该时间时，会被删除，时间戳
}

// MapT map with ttl
type MapT struct {
	m sync.Map
}

// NewMapT new MapT
func NewMapT() *MapT {
	var m MapT
	go func() {
		for now := range time.Tick(time.Second) {
			m.m.Range(func(key, value interface{}) bool {
				v, ok := value.(item)
				if !ok {
					log.Error(`mapT value is not an item`)
					return false
				}
				if now.Unix() < v.timeToDel {
					return true
				}
				log.Info(`开始删除`)
				m.delete(key)
				return true
			})
		}
	}()
	return &m
}

func (t *MapT) delete(key interface{}) {
	if t == nil {
		log.Error(`t *MapT is nil`)
		return
	}
	t.m.Delete(key)
}

// Store sets the value for a key.
func (t *MapT) Store(key interface{}, value interface{}, ttl time.Duration) {
	if t == nil {
		log.Error(`t *MapT is nil`)
		return
	}
	var item item
	item.value = value
	item.timeToDel = time.Now().Add(ttl).Unix()
	t.m.Store(key, item)
}

// Delete deletes the value for a key.
func (t *MapT) Delete(key interface{}) {
	t.delete(key)
}

// Load returns the value stored in the map for a key, or nil if no value is present.
// The ok result indicates whether value was found in the map.
func (t *MapT) Load(key interface{}) (value interface{}, ok bool) {
	if t == nil {
		log.Error(`t *MapT is nil`)
		return nil, false
	}
	v, ok := t.m.Load(key)
	if !ok {
		return nil, false
	}
	if v == nil {
		return nil, false
	}
	item, ok := v.(item)
	if !ok {
		log.Error(`value's type is not item`)
		return nil, false
	}
	return item.value, ok
}
