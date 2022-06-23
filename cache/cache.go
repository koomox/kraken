package cache

import (
	"encoding/json"
	"github.com/koomox/kraken/redblacktree"
	"sync"
)

type Store struct {
	sync.RWMutex

	tree *redblacktree.Tree
}

type Element struct {
	Key     interface{}
	Payload interface{}
}

func NewWithStringComparator() *Store {
	return &Store{
		tree: redblacktree.NewWith(redblacktree.StringComparator),
	}
}

func NewWithIntComparator() *Store {
	return &Store{
		tree: redblacktree.NewWith(redblacktree.IntComparator),
	}
}

func NewWithInt64Comparator() *Store {
	return &Store{
		tree: redblacktree.NewWith(redblacktree.Int64Comparator),
	}
}

func (r *Store) Put(key, payload interface{}) {
	r.Lock()
	defer r.Unlock()

	r.tree.Put(key, &Element{
		Key:     key,
		Payload: payload,
	})
}

func (r *Store) Get(key interface{}) interface{} {
	r.RLock()
	defer r.RUnlock()

	if v, ok := r.tree.Get(key); ok {
		return v.(*Element).Payload
	}
	return nil
}

func (r *Store) Remove(key interface{}) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.tree.Get(key); ok {
		r.tree.Remove(key)
	}
}

func (r *Store) Cleanup() {
	r.Lock()
	defer r.Unlock()

	it := r.tree.Iterator()
	for it.Next() {
		v := it.Value().(*Element)
		r.tree.Remove(v.Key)
	}
}

func (r *Store) Values() (m []interface{}) {
	it := r.tree.Iterator()
	for it.Next() {
		v := it.Value().(*Element)
		m = append(m, v.Payload)
	}
	return
}

func (r *Store) ToJSON() ([]byte, error) {
	var m []interface{}
	it := r.tree.Iterator()
	for it.Next() {
		v := it.Value().(*Element)
		m = append(m, v.Payload)
	}

	return json.Marshal(m)
}

func (r *Store) CallbackFunc(callbackFunc func(interface{})) {
	it := r.tree.Iterator()
	for it.Next() {
		callbackFunc(it.Value())
	}
}

func (r *Store) CancelFunc(callbackFunc func(interface{})bool) {
	it := r.tree.Iterator()
	for it.Next() {
		callbackFunc(it.Value())
	}
}