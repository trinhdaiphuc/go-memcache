package handler

import "github.com/trinhdaiphuc/go-memcache/resp"

type PingHandler struct {
}

func (p *PingHandler) Handle(ctx Context, args []resp.Expression) resp.Expression {
	return resp.NewSimpleStringExpression("PONG")
}

func NewPingHandler() Handler {
	return &PingHandler{}
}
