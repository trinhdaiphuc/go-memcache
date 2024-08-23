package memcache

import "time"

type Map[K, V comparable] interface {
	Set(key K, value V)
	Get(key K) (V, bool)
	Delete(key K)
	Keys() []K
	Values() []V
	Len() int
	TTLKey(key K) time.Duration
	TTL() time.Duration
	ExpireKey(key K, ttl time.Duration)
	Expire(ttl time.Duration)
	IsExpired() bool
}

func NewMap[K, V comparable]() Map[K, V] {
	m := &mapData[K, V]{
		data:           make(map[K]*MapValue[V]),
		ttl:            0,
		lastAccessTime: time.Now(),
		command:        make(chan CommandMap[K, V]),
	}

	go m.executeCommands()

	return m
}

type mapData[K, V comparable] struct {
	data           map[K]*MapValue[V]
	ttl            time.Duration
	lastAccessTime time.Time
	command        chan CommandMap[K, V]
}

func (m *mapData[K, V]) Set(key K, value V) {
	m.command <- &SetMapCommand[K, V]{Key: key, Value: value}
}

func (m *mapData[K, V]) Get(key K) (value V, ok bool) {
	response := make(chan *GetResponse[V])
	m.command <- &GetMapCommand[K, V]{Key: key, Response: response}

	res := <-response
	if res.Found {
		return res.Value, true
	}
	return value, false
}

func (m *mapData[K, V]) Delete(key K) {
	m.command <- &DeleteMapCommand[K, V]{Key: key}
}

func (m *mapData[K, V]) Keys() []K {
	keys := make(chan []K)
	m.command <- &GetKeysCommand[K, V]{Response: keys}
	return <-keys
}

func (m *mapData[K, V]) Values() []V {
	value := make(chan []V)
	m.command <- &GetValuesCommand[K, V]{Response: value}
	return <-value
}

func (m *mapData[K, V]) Len() int {
	return len(m.data)
}

func (m *mapData[K, V]) ExpireKey(key K, ttl time.Duration) {
	m.command <- &ExpireMapCommand[K, V]{Key: key, TTL: ttl}
}

func (m *mapData[K, V]) Expire(ttl time.Duration) {
	m.ttl = ttl
	m.lastAccessTime = time.Now()
}

func (m *mapData[K, V]) TTLKey(key K) time.Duration {
	ttl := make(chan time.Duration)
	m.command <- &TTLMapCommand[K, V]{Key: key, Response: ttl}
	return <-ttl
}

func (m *mapData[K, V]) TTL() time.Duration {
	return m.ttl
}

func (m *mapData[K, V]) IsExpired() bool {
	return m.ttl > 0 && time.Since(m.lastAccessTime) > m.ttl
}

func (m *mapData[K, V]) executeCommands() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case cmd := <-m.command:
			cmd.Execute(m)
		case <-ticker.C:
			for k, v := range m.data {
				if v.IsExpired() {
					delete(m.data, k)
				}
			}
		}
	}
}

type MapValue[V comparable] struct {
	value          V
	ttl            time.Duration
	lastAccessTime time.Time
}

func NewMapValue[V comparable](value V, ttl time.Duration) *MapValue[V] {
	return &MapValue[V]{
		value:          value,
		ttl:            ttl,
		lastAccessTime: time.Now(),
	}
}

func (m *MapValue[V]) Value() V {
	return m.value
}

func (m *MapValue[V]) TTL() time.Duration {
	return m.ttl
}

func (m *MapValue[V]) SetValue(value V) {
	m.value = value
	m.lastAccessTime = time.Now()
	if m.IsExpired() {
		m.ttl = 0
	}
}

func (m *MapValue[V]) Expire(ttl time.Duration) {
	m.ttl = ttl
}

func (m *MapValue[V]) IsExpired() bool {
	return m.ttl > 0 && time.Since(m.lastAccessTime) > m.ttl
}
