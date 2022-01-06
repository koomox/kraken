package kraken

import "time"

type Store interface {
	Put(string, interface{}, time.Duration)
	Get(string) interface{}
	Remove(string)
	GetWithExpire(key string) (interface{}, time.Time)
	Values() []interface{}
	ToJSON() ([]byte, error)
	CallbackFunc(func(interface{}))
}

type Cache interface {
	Put(string, interface{})
	Get(string) interface{}
	Remove(string)
	Cleanup()
	Values() []interface{}
	ToJSON() ([]byte, error)
	CallbackFunc(func(interface{}))
}

type Bucket interface {
	Put(int, interface{})
	Get(int) interface{}
	Remove(int)
	Cleanup()
	Values() []interface{}
	ToJSON() ([]byte, error)
	CallbackFunc(func(interface{}))
	CancelFunc(func(interface{}) bool)
}
