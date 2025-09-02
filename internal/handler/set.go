package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/trinhdaiphuc/go-memcache/resp"
)

type SetHandler struct {
}

func NewSetHandler() Handler {
	return &SetHandler{}
}

const (
	EX = "EX"
	PX = "PX"
)

func (s *SetHandler) Handle(ctx Context, args []resp.Expression) resp.Expression {
	key := args[0].Value().(string)
	value := args[1].Value().(string)
	ctx.Map.Set(key, value)

	if len(args) > 2 {
		exOption := args[2].Value().(string)
		expiration, err := strconv.Atoi(args[3].Value().(string))
		if err != nil {
			return resp.NewErrorExpression("Invalid expiration value")
		}
		switch strings.ToUpper(exOption) {
		case EX:
			ctx.Map.ExpireKey(key, time.Duration(expiration)*time.Second)
		case PX:
			ctx.Map.ExpireKey(key, time.Duration(expiration)*time.Millisecond)
		}
	}
	return resp.NewSimpleStringExpression("OK")
}
