package handler

import (
	"time"

	"github.com/trinhdaiphuc/go-memcache/resp"
)

type ExpiredHandler struct {
}

func NewExpiredHandler() Handler {
	return &ExpiredHandler{}
}

func (e *ExpiredHandler) Handle(ctx Context, args []resp.Expression) resp.Expression {
	key := args[0].Value().(string)
	ttl := args[1].Value().(int)
	ctx.Map.ExpireKey(key, time.Second*time.Duration(ttl))
	ctx.Has.Expire(key, time.Second*time.Duration(ttl))
	return resp.NewIntegerExpression(1)
}
