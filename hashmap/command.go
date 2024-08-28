package hashmap

import "github.com/trinhdaiphuc/go-memcache/gomap"

type CommandHashMap[K, V comparable] interface {
	Execute(data *hashMap[K, V])
}

type setCommand[K, V comparable] struct {
	key       K
	keyValues []KeyValue[K, V]
}

func (c *setCommand[K, V]) Execute(hashMap *hashMap[K, V]) {
	mapData, ok := hashMap.data[c.key]
	if !ok {
		mapData = gomap.NewMap[K, V]()
		hashMap.data[c.key] = mapData
	}

	for _, kv := range c.keyValues {
		mapData.Set(kv.Key, kv.Value)
	}
}

type getResponse[K, V comparable] struct {
	mapValue gomap.Map[K, V]
	found    bool
}

type getCommand[K, V comparable] struct {
	key      K
	response chan *getResponse[K, V]
}

func (c *getCommand[K, V]) Execute(hashMap *hashMap[K, V]) {
	mapData, ok := hashMap.data[c.key]
	if !ok {
		c.response <- &getResponse[K, V]{found: false}
	} else {
		c.response <- &getResponse[K, V]{mapValue: mapData, found: true}
	}
	close(c.response)
}
