package gomap

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
		data:           make(map[K]*mapValue[V]),
		ttl:            0,
		lastAccessTime: time.Now(),
		command:        make(chan CommandMap[K, V]),
	}

	go m.executeCommands()

	return m
}

type mapData[K, V comparable] struct {
	data           map[K]*mapValue[V]
	ttl            time.Duration
	lastAccessTime time.Time
	command        chan CommandMap[K, V]
}

func (m *mapData[K, V]) Set(key K, value V) {
	m.command <- &setCommand[K, V]{key: key, value: value}
}

func (m *mapData[K, V]) Get(key K) (value V, ok bool) {
	response := make(chan *getResponse[V])
	m.command <- &getCommand[K, V]{key: key, response: response}

	res := <-response
	if res.found {
		return res.value, true
	}
	return value, false
}

func (m *mapData[K, V]) Delete(key K) {
	m.command <- &deleteCommand[K, V]{key: key}
}

func (m *mapData[K, V]) delete(key K) {
	delete(m.data, key)
}

func (m *mapData[K, V]) Keys() []K {
	keys := make(chan []K)
	m.command <- &getKeysCommand[K, V]{response: keys}
	return <-keys
}

func (m *mapData[K, V]) Values() []V {
	value := make(chan []V)
	m.command <- &getValuesCommand[K, V]{response: value}
	return <-value
}

func (m *mapData[K, V]) Len() int {
	return len(m.data)
}

func (m *mapData[K, V]) ExpireKey(key K, ttl time.Duration) {
	m.command <- &expireKeyCommand[K, V]{key: key, ttl: ttl}
}

func (m *mapData[K, V]) Expire(ttl time.Duration) {
	m.ttl = ttl
	m.lastAccessTime = time.Now()
}

func (m *mapData[K, V]) TTLKey(key K) time.Duration {
	ttl := make(chan time.Duration)
	m.command <- &ttlKeyCommand[K, V]{key: key, response: ttl}
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
			m.clearExpiredData()
			cmd.Execute(m)
		case <-ticker.C:
			m.clearExpiredData()
		}
	}
}

func (m *mapData[K, V]) clearExpiredData() {
	if m.IsExpired() {
		m.data = make(map[K]*mapValue[V])
		m.ttl = 0
		m.updateLastAccessTime()
		return
	}

	for k, v := range m.data {
		if v.IsExpired() {
			delete(m.data, k)
		}
	}
}

func (m *mapData[K, V]) updateLastAccessTime() {
	m.lastAccessTime = time.Now()
}

type mapValue[V comparable] struct {
	value          V
	ttl            time.Duration
	lastAccessTime time.Time
}

func newMapValue[V comparable](value V, ttl time.Duration) *mapValue[V] {
	return &mapValue[V]{
		value:          value,
		ttl:            ttl,
		lastAccessTime: time.Now(),
	}
}

func (m *mapValue[V]) Value() V {
	return m.value
}

func (m *mapValue[V]) TTL() time.Duration {
	return m.ttl
}

func (m *mapValue[V]) SetValue(value V) {
	m.value = value

	if m.IsExpired() {
		m.ttl = 0
	}

	m.lastAccessTime = time.Now()
}

func (m *mapValue[V]) Expire(ttl time.Duration) {
	m.ttl = ttl
}

func (m *mapValue[V]) IsExpired() bool {
	return m.ttl > 0 && time.Since(m.lastAccessTime) > m.ttl
}
