package handler

import "github.com/trinhdaiphuc/go-memcache/resp"

type SetHandler struct {
}

func NewSetHandler() Handler {
	return &SetHandler{}
}

func (s *SetHandler) Handle(ctx Context, args []resp.Expression) resp.Expression {
	key := args[0].Value().(string)
	value := args[1].Value().(string)
	ctx.Map.Set(key, value)
	return resp.NewSimpleStringExpression("OK")
}
