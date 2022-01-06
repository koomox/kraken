package memory

import (
	"github.com/koomox/kraken/redblacktree"
	"sync"
	"time"
)

type Store struct {
	sync.RWMutex

	tree *redblacktree.Tree
}

type Element struct {
	Expired time.Time
	Key     string
	Payload interface{}
}

func NewStore() *Store {
	return &Store{
		tree: redblacktree.NewWithStringComparator(),
	}
}

func (r *Store) Put(key string, payload interface{}, ttl time.Duration) {
	r.tree.Put(key, &Element{
		Key:     key,
		Payload: payload,
		Expired: time.Now().Add(ttl),
	})
}

func (r *Store) Get(key string) interface{} {
	v, ok := r.tree.Get(key)
	if !ok {
		return nil
	}
	element := v.(*Element)
	if time.Since(element.Expired) > 0 {
		r.tree.Remove(key)
		return nil
	}

	return element.Payload
}

func (r *Store) Remove(key string) {
	if _, ok := r.tree.Get(key); ok {
		r.tree.Remove(key)
	}
}

func (r *Store) GetWithExpire(key string) (payload interface{}, expired time.Time) {
	v, ok := r.tree.Get(key)
	if !ok {
		return
	}
	element := v.(*Element)
	if time.Since(element.Expired) > 0 {
		r.tree.Remove(key)
		return
	}

	return element.Payload, element.Expired
}

func (r *Store) cleanup() {
	it := r.tree.Iterator()
	for it.Next() {
		v := it.Value().(*Element)
		if time.Since(v.Expired) > 0 {
			r.tree.Remove(v.Key)
		}
	}
}
