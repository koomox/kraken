package kraken

import "time"

type Store interface {
	Put(string, interface{}, time.Duration)
	Get(string) interface{}
	Remove(string)
	GetWithExpire(key string) (interface{}, time.Time)
}

type Cache interface {
	Put(string, interface{})
	Get(string) interface{}
	Remove(string)
	Cleanup()
}
