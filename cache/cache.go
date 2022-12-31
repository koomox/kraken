package memory

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/koomox/kraken/redblacktree"
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
	r.Lock()
	defer r.Unlock()

	r.tree.Put(key, &Element{
		Key:     key,
		Payload: payload,
		Expired: time.Now().Add(ttl),
	})
}

func (r *Store) Get(key interface{}) interface{} {
	r.RLock()
	defer r.RUnlock()

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
	r.Lock()
	defer r.Unlock()

	if _, ok := r.tree.Get(key); ok {
		r.tree.Remove(key)
	}
}

func (r *Store) GetWithExpire(key interface{}) (payload interface{}, expired time.Time) {
	r.RLock()
	defer r.RUnlock()

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

func (r *Store) Cleanup() {
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
		m = append(m, it.Value().(*Element).Payload)
	}
	return
}

func (r *Store) ToJSON() ([]byte, error) {
	var m []interface{}
	it := r.tree.Iterator()
	for it.Next() {
		m = append(m, it.Value().(*Element).Payload)
	}

	return json.Marshal(m)
}

func (r *Store) CallbackFunc(callbackFunc func(interface{})) {
	it := r.tree.Iterator()
	for it.Next() {
		callbackFunc(it.Value().(*Element).Payload)
	}
}

func (r *Store) CancelFunc(callbackFunc func(interface{}) bool) {
	it := r.tree.Iterator()
	for it.Next() {
		if callbackFunc(it.Value().(*Element).Payload) {
			break
		}
	}
}
