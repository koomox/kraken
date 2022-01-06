package kraken

import "time"

type Store interface {
	Put(string, interface{}, time.Duration)
	Get(string) interface{}
	Remove(string)
	GetWithExpire(key string) (interface{}, time.Time)
}
