package hashmap

import (
	"time"

	"github.com/trinhdaiphuc/go-memcache/gomap"
)

type KeyValue[K, V comparable] struct {
	Key   K
	Value V
}

type HashMap[K, V comparable] interface {
	Get(key K) (gomap.Map[K, V], bool)
	Set(key K, keyValues ...KeyValue[K, V])
	Delete(key K)
	Keys() []K
	Values() []gomap.Map[K, V]
	Len() int
	TTL(key K) time.Duration
	Expire(key K, ttl time.Duration)
	IsExpired() bool
}

type hashMap[K, V comparable] struct {
	data map[K]gomap.Map[K, V]
}

func NewHashMap[K, V comparable]() HashMap[K, V] {
	return &hashMap[K, V]{
		data: make(map[K]gomap.Map[K, V]),
	}
}

func (h *hashMap[K, V]) Get(key K) (gomap.Map[K, V], bool) {
	v, ok := h.data[key]
	return v, ok
}

func (h *hashMap[K, V]) Set(key K, keyValues ...KeyValue[K, V]) {
	// TODO implement me
	panic("implement me")
}

func (h *hashMap[K, V]) Delete(key K) {
	// TODO implement me
	panic("implement me")
}

func (h *hashMap[K, V]) Keys() []K {
	// TODO implement me
	panic("implement me")
}

func (h *hashMap[K, V]) Values() []gomap.Map[K, V] {
	// TODO implement me
	panic("implement me")
}

func (h *hashMap[K, V]) Len() int {
	// TODO implement me
	panic("implement me")
}

func (h *hashMap[K, V]) TTL(key K) time.Duration {
	// TODO implement me
	panic("implement me")
}

func (h *hashMap[K, V]) Expire(key K, ttl time.Duration) {
	// TODO implement me
	panic("implement me")
}

func (h *hashMap[K, V]) IsExpired() bool {
	// TODO implement me
	panic("implement me")
}
