package memory

import (
	"encoding/json"
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

func (r *Store) Put(key, payload interface{}, ttl time.Duration) {
	r.tree.Put(key, &Element{
		Key:     key,
		Payload: payload,
		Expired: time.Now().Add(ttl),
	})
}

func (r *Store) Get(key interface{}) interface{} {
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

func (r *Store) Remove(key interface{}) {
	if _, ok := r.tree.Get(key); ok {
		r.tree.Remove(key)
	}
}

func (r *Store) GetWithExpire(key interface{}) (payload interface{}, expired time.Time) {
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