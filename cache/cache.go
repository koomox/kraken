package cache

import (
	"github.com/koomox/kraken/redblacktree"
	"sync"
)

type Store struct {
	sync.RWMutex

	tree *redblacktree.Tree
}

type Element struct {
	Key     string
	Payload interface{}
}

func NewStore() *Store {
	return &Store{
		tree: redblacktree.NewWithStringComparator(),
	}
}

func (r *Store) Put(key string, payload interface{}) {
	r.tree.Put(key, &Element{
		Key:     key,
		Payload: payload,
	})
}

func (r *Store) Get(key string) interface{} {
	if v, ok := r.tree.Get(key); ok {
		return v.(*Element).Payload
	}
	return nil
}

func (r *Store) Remove(key string) {
	if _, ok := r.tree.Get(key); ok {
		r.tree.Remove(key)
	}
}

func (r *Store) Cleanup() {
	it := r.tree.Iterator()
	for it.Next() {
		v := it.Value().(*Element)
		r.tree.Remove(v.Key)
	}
}
