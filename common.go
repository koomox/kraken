package kraken

import "time"

type Cache interface {
	Put(interface{}, interface{}, time.Duration)
	Get(interface{}) interface{}
	Remove(interface{})
	GetWithExpire(interface{}) (interface{}, time.Time)
	Values() []interface{}
	ToJSON() ([]byte, error)
	CallbackFunc(func(interface{}))
	CancelFunc(func(interface{})bool)
}

type Tree interface {
	Put(interface{}, interface{})
	Get(interface{}) interface{}
	Remove(interface{})
	Cleanup()
	Values() []interface{}
	ToJSON() ([]byte, error)
	CallbackFunc(func(interface{}))
	CancelFunc(func(interface{})bool)
}