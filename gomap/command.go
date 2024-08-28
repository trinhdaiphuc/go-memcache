package gomap

import "time"

type CommandMap[K, V comparable] interface {
	Execute(data *mapData[K, V])
}

type setCommand[K, V comparable] struct {
	key   K
	value V
}

func (c *setCommand[K, V]) Execute(mapData *mapData[K, V]) {
	v, ok := mapData.data[c.key]
	if !ok {
		mapData.data[c.key] = newMapValue[V](c.value, 0)
		return
	}

	v.SetValue(c.value)
}

type getCommand[K, V comparable] struct {
	key      K
	response chan *getResponse[V]
}

type getResponse[V comparable] struct {
	value V
	found bool
}

func (c *getCommand[K, V]) Execute(mapData *mapData[K, V]) {
	v, ok := mapData.data[c.key]
	if !ok {
		c.response <- &getResponse[V]{found: false}
	} else {
		c.response <- &getResponse[V]{value: v.Value(), found: true}
	}
	close(c.response)
}

type deleteCommand[K, V comparable] struct {
	key K
}

func (c *deleteCommand[K, V]) Execute(mapData *mapData[K, V]) {
	delete(mapData.data, c.key)
}

type getKeysCommand[K, V comparable] struct {
	response chan []K
}

func (c *getKeysCommand[K, V]) Execute(mapData *mapData[K, V]) {
	keys := make([]K, 0, len(mapData.data))
	for k := range mapData.data {
		keys = append(keys, k)
	}
	c.response <- keys
	close(c.response)
}

type getValuesCommand[K, V comparable] struct {
	response chan []V
}

func (c *getValuesCommand[K, V]) Execute(mapData *mapData[K, V]) {
	values := make([]V, 0, len(mapData.data))
	for _, v := range mapData.data {
		values = append(values, v.Value())
	}
	c.response <- values
	close(c.response)
}

type expireKeyCommand[K, V comparable] struct {
	key K
	ttl time.Duration
}

func (c *expireKeyCommand[K, V]) Execute(mapData *mapData[K, V]) {
	v, ok := mapData.data[c.key]
	if !ok {
		return
	}
	v.Expire(c.ttl)
}

type ttlKeyCommand[K, V comparable] struct {
	key      K
	response chan time.Duration
}

func (c *ttlKeyCommand[K, V]) Execute(mapData *mapData[K, V]) {
	v, ok := mapData.data[c.key]
	if !ok {
		c.response <- 0
	} else {
		c.response <- v.TTL()
	}
	close(c.response)
}
