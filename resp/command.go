package resp

import "bufio"

type RESP struct {
}

func NewRESP() *RESP {
	return &RESP{}
}

func (r *RESP) Read(reader *bufio.Reader) (Expression, error) {
	var expression Expression
	b, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	switch b {
	case String:
		expression = &SimpleStringExpression{}
	case Error:
		expression = &ErrorExpression{}
	case Integer:
		expression = &IntegerExpression{}
	case BulkString:
		expression = &BulkStringExpression{}
	case Array:
		expression = &ArrayExpression{}
	case CR:
		expression = &CRExpression{}
	}
	err = expression.Read(reader)
	if err != nil {
		return nil, err
	}
	return expression, nil
}

type Command struct {
	Simple Expression
	Array  *ArrayExpression
}

func NewRESPCommand() *Command {
	return &Command{}
}

func (r *Command) Serialize() string {
	if r.Simple != nil {
		return r.Simple.Serialize()
	}
	return r.Array.Serialize()
}

func (r *Command) Read(reader *bufio.Reader) error {
	expression, err := NewRESP().Read(reader)
	if err != nil {
		return err
	}
	if arrayExpression, ok := expression.(*ArrayExpression); ok {
		r.Array = arrayExpression
	} else {
		r.Simple = expression
	}
	return nil
}
