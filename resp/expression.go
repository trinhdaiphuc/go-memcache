package resp

import (
	"bufio"
	"strconv"
)

const (
	String     = '+'
	Error      = '-'
	Array      = '*'
	BulkString = '$'
	Integer    = ':'
	CR         = '\r'
	LF         = '\n'
)

type Expression interface {
	Serialize() string
	String() string
	Value() interface{}
	Read(reader *bufio.Reader) error
}

type SimpleStringExpression struct {
	value string
}

func NewSimpleStringExpression(value string) Expression {
	return &SimpleStringExpression{value: value}
}

func (s *SimpleStringExpression) Serialize() string {
	return string(String) + s.value + "\r\n"
}

func (s *SimpleStringExpression) Read(reader *bufio.Reader) error {
	line, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	s.value = string(line)
	return nil
}

func (s *SimpleStringExpression) Value() interface{} {
	return s.Value
}

func (s *SimpleStringExpression) String() string {
	return "String: " + s.value
}

type ErrorExpression struct {
	value string
}

func NewErrorExpression(value string) Expression {
	return &ErrorExpression{value: value}
}

func (e *ErrorExpression) Serialize() string {
	return string(Error) + e.value + "\r\n"
}

func (e *ErrorExpression) Read(reader *bufio.Reader) error {
	line, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	e.value = string(line)
	return nil
}

func (e *ErrorExpression) String() string {
	return "Error: " + e.value
}

func (e *ErrorExpression) Value() interface{} {
	return e.value
}

type IntegerExpression struct {
	value int
}

func NewIntegerExpression(value int) Expression {
	return &IntegerExpression{value: value}
}

func (i *IntegerExpression) Serialize() string {
	return string(Integer) + strconv.Itoa(i.value) + "\r\n"
}

func (i *IntegerExpression) Read(reader *bufio.Reader) error {
	line, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	value, err := strconv.Atoi(string(line))
	if err != nil {
		return err
	}
	i.value = value
	return nil
}

func (i *IntegerExpression) String() string {
	return "Integer: " + strconv.Itoa(i.value)
}

func (i *IntegerExpression) Value() interface{} {
	return i.value
}

type BulkStringExpression struct {
	value string
}

func NewBulkStringExpression(value string) Expression {
	return &BulkStringExpression{value: value}
}

func (b *BulkStringExpression) Serialize() string {
	return string(BulkString) + strconv.Itoa(len(b.value)) + "\r\n" + b.value + "\r\n"
}

func (b *BulkStringExpression) Read(reader *bufio.Reader) error {
	line, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	size, err := strconv.Atoi(string(line))
	if err != nil {
		return err
	}
	buf := make([]byte, size)
	_, err = reader.Read(buf)
	if err != nil {
		return err
	}
	b.value = string(buf)
	return nil
}

func (b *BulkStringExpression) String() string {
	return "BulkString: " + b.value
}

func (b *BulkStringExpression) Value() interface{} {
	return b.value
}

type NullBulkStringExpression struct {
}

func NewNullBulkStringExpression() Expression {
	return &NullBulkStringExpression{}
}

func (n *NullBulkStringExpression) Serialize() string {
	return string(BulkString) + "-1\r\n"
}

func (n *NullBulkStringExpression) Read(reader *bufio.Reader) error {
	return nil
}

func (n *NullBulkStringExpression) String() string {
	return "NullBulkString"
}

func (n *NullBulkStringExpression) Value() interface{} {
	return nil
}

type NullArrayExpression struct {
}

type CRExpression struct {
}

func (c *CRExpression) Serialize() string {
	return string(CR)
}

func (c *CRExpression) Read(reader *bufio.Reader) error {
	_, err := reader.ReadByte()
	return err
}

func (c *CRExpression) String() string {
	return "CR"
}

func (c *CRExpression) Value() interface{} {
	return ""
}

func (n *NullArrayExpression) Serialize() string {
	return string(Array) + "-1\r\n"
}

type ArrayExpression struct {
	Expressions []Expression
}

func (a *ArrayExpression) Serialize() string {
	var serialized string
	serialized += string(Array) + strconv.Itoa(len(a.Expressions)) + "\r\n"
	for _, expression := range a.Expressions {
		serialized += expression.Serialize()
	}
	return serialized
}

func (a *ArrayExpression) Read(reader *bufio.Reader) error {
	line, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	size, err := strconv.Atoi(string(line))
	if err != nil {
		return err
	}
	for i := 0; i < size; {
		expression, err := NewRESP().Read(reader)
		if err != nil {
			return err
		}
		if _, ok := expression.(*CRExpression); ok {
			continue
		}
		a.Expressions = append(a.Expressions, expression)
		i++
	}
	return nil
}

func (a *ArrayExpression) Value() interface{} {
	return ""
}

func (a *ArrayExpression) String() string {
	str := "Array: " + strconv.Itoa(len(a.Expressions)) + "\n"
	for _, expression := range a.Expressions {
		str += expression.String() + "\n"
	}
	return str
}
