package handler

import "github.com/trinhdaiphuc/go-memcache/resp"

type EchoHandler struct {
}

func (e *EchoHandler) Handle(ctx Context, args []resp.Expression) resp.Expression {
	value := args[0].Value().(string)
	return resp.NewSimpleStringExpression(value)
}

func NewEchoHandler() Handler {
	return &EchoHandler{}
}
