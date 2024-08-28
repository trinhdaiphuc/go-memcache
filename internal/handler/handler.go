package handler

import (
	"github.com/trinhdaiphuc/go-memcache/gomap"
	"github.com/trinhdaiphuc/go-memcache/hashmap"
	"github.com/trinhdaiphuc/go-memcache/resp"
)

const (
	PING    = "PING"
	GET     = "GET"
	SET     = "SET"
	EXPIRED = "EXPIRE"
)

type Func func([]resp.Expression) resp.Expression

type Context struct {
	Map gomap.Map[string, string]
	Has hashmap.HashMap[string, string]
}

type Map map[string]Handler

type Handler interface {
	Handle(ctx Context, args []resp.Expression) resp.Expression
}

func NewContext(m gomap.Map[string, string], h hashmap.HashMap[string, string]) Context {
	return Context{
		Map: m,
		Has: h,
	}
}

func NewMap() Map {
	return Map{
		PING:    NewPingHandler(),
		GET:     NewGetHandler(),
		SET:     NewSetHandler(),
		EXPIRED: NewExpiredHandler(),
	}
}
