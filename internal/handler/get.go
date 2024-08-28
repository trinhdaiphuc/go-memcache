package handler

import (
	"github.com/trinhdaiphuc/go-memcache/resp"
)

type GetHandler struct {
}

func NewGetHandler() Handler {
	return &GetHandler{}
}

func (g *GetHandler) Handle(ctx Context, args []resp.Expression) resp.Expression {
	key := args[0].Value().(string)
	value, ok := ctx.Map.Get(key)
	if !ok {
		return resp.NewNullBulkStringExpression()
	}
	return resp.NewBulkStringExpression(value)
}
