package memcache

import "time"

type CommandMap[K, V comparable] interface {
	Execute(data *mapData[K, V])
}

type SetMapCommand[K, V comparable] struct {
	Key   K
	Value V
}

func (c *SetMapCommand[K, V]) Execute(mapData *mapData[K, V]) {
	v, ok := mapData.data[c.Key]
	if !ok {
		mapData.data[c.Key] = NewMapValue[V](c.Value, 0)
		return
	}

	v.SetValue(c.Value)
}

type GetMapCommand[K, V comparable] struct {
	Key      K
	Response chan *GetResponse[V]
}

type GetResponse[V comparable] struct {
	Value V
	Found bool
}

func (c *GetMapCommand[K, V]) Execute(mapData *mapData[K, V]) {
	v, ok := mapData.data[c.Key]
	if !ok || v.IsExpired() {
		c.Response <- &GetResponse[V]{Found: false}
	} else {
		c.Response <- &GetResponse[V]{Value: v.Value(), Found: true}
	}
	close(c.Response)
}

type DeleteMapCommand[K, V comparable] struct {
	Key K
}

func (c *DeleteMapCommand[K, V]) Execute(mapData *mapData[K, V]) {
	delete(mapData.data, c.Key)
}

type GetKeysCommand[K, V comparable] struct {
	Response chan []K
}

func (c *GetKeysCommand[K, V]) Execute(mapData *mapData[K, V]) {
	keys := make([]K, 0, len(mapData.data))
	for k := range mapData.data {
		keys = append(keys, k)
	}
	c.Response <- keys
	close(c.Response)
}

type GetValuesCommand[K, V comparable] struct {
	Response chan []V
}

func (c *GetValuesCommand[K, V]) Execute(mapData *mapData[K, V]) {
	values := make([]V, 0, len(mapData.data))
	for _, v := range mapData.data {
		values = append(values, v.Value())
	}
	c.Response <- values
	close(c.Response)
}

type ExpireMapCommand[K, V comparable] struct {
	Key K
	TTL time.Duration
}

func (c *ExpireMapCommand[K, V]) Execute(mapData *mapData[K, V]) {
	v, ok := mapData.data[c.Key]
	if !ok {
		return
	}
	v.Expire(c.TTL)
}

type TTLMapCommand[K, V comparable] struct {
	Key      K
	Response chan time.Duration
}

func (c *TTLMapCommand[K, V]) Execute(mapData *mapData[K, V]) {
	v, ok := mapData.data[c.Key]
	if !ok {
		c.Response <- 0
	} else {
		c.Response <- v.TTL()
	}
	close(c.Response)
}
