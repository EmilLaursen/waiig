package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/EmilLaursen/wiig/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

// String implements Node.
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

var _ Node = &Program{}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

var _ Statement = &LetStatement{}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

var _ Statement = &BlockStatement{}

func (n *BlockStatement) statementNode()       {}
func (n *BlockStatement) TokenLiteral() string { return n.Token.Literal }
func (n *BlockStatement) String() string {
	var out bytes.Buffer
	for _, stmt := range n.Statements {
		out.WriteString(stmt.String())
	}
	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

var _ Expression = &Identifier{}

func (n *Identifier) expressionNode()      {}
func (n *Identifier) TokenLiteral() string { return n.Token.Literal }
func (n *Identifier) String() string {
	return n.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

var _ Expression = &IntegerLiteral{}

func (n *IntegerLiteral) expressionNode()      {}
func (n *IntegerLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *IntegerLiteral) String() string       { return n.Token.Literal }

type Boolean struct {
	Token token.Token
	Value bool
}

var _ Expression = &Boolean{}

func (n *Boolean) expressionNode()      {}
func (n *Boolean) TokenLiteral() string { return n.Token.Literal }
func (n *Boolean) String() string       { return n.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

var _ Expression = &StringLiteral{}

func (n *StringLiteral) expressionNode()      {}
func (n *StringLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *StringLiteral) String() string       { return n.Token.Literal }

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

var _ Expression = &IfExpression{}

func (n *IfExpression) expressionNode()      {}
func (n *IfExpression) TokenLiteral() string { return n.Token.Literal }
func (n *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(n.Condition.String())
	out.WriteString(" ")
	out.WriteString(n.Consequence.String())
	if n.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(n.Alternative.String())
	}
	return out.String()
}

type FunctionLiteral struct {
	Token  token.Token
	Params []*Identifier
	Body   *BlockStatement
}

var _ Expression = &FunctionLiteral{}

func (n *FunctionLiteral) expressionNode()      {}
func (n *FunctionLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range n.Params {
		params = append(params, p.String())
	}
	out.WriteString(n.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	return out.String()
}

type ArrayLiteral struct {
	Token token.Token
	Elems []Expression
}

var _ Expression = &ArrayLiteral{}

func (n *ArrayLiteral) expressionNode()      {}
func (n *ArrayLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *ArrayLiteral) String() string {
	var out bytes.Buffer
	last := len(n.Elems) - 1
	out.WriteString("[")
	if last >= 0 {
		for _, el := range n.Elems[:last] {
			out.WriteString(el.String())
			out.WriteString(", ")
		}
		out.WriteString(n.Elems[last].String())
	}
	out.WriteString("]")
	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

var _ Expression = &HashLiteral{}

func (n *HashLiteral) expressionNode()      {}
func (n *HashLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for k, v := range n.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s:%s", k.String(), v.String()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

var _ Expression = &IndexExpression{}

func (n *IndexExpression) expressionNode()      {}
func (n *IndexExpression) TokenLiteral() string { return n.Token.Literal }
func (n *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(n.Left.String())
	out.WriteString("[")
	out.WriteString(n.Index.String())
	out.WriteString("])")
	return out.String()
}

type SliceExpression struct {
	Token      token.Token
	Left       Expression
	IndexLeft  Expression
	IndexRight Expression
}

var _ Expression = &SliceExpression{}

func (n *SliceExpression) expressionNode()      {}
func (n *SliceExpression) TokenLiteral() string { return n.Token.Literal }
func (n *SliceExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(n.Left.String())
	out.WriteString("[")
	out.WriteString(n.IndexLeft.String())
	out.WriteString(":")
	out.WriteString(n.IndexRight.String())
	out.WriteString("])")
	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

var _ Expression = &CallExpression{}

func (n *CallExpression) expressionNode()      {}
func (n *CallExpression) TokenLiteral() string { return n.Token.Literal }
func (n *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, p := range n.Arguments {
		args = append(args, p.String())
	}
	out.WriteString(n.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

var _ Expression = &PrefixExpression{}

func (n *PrefixExpression) expressionNode()      {}
func (n *PrefixExpression) TokenLiteral() string { return n.Token.Literal }
func (n *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(n.Operator)
	out.WriteString(n.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

var _ Expression = &InfixExpression{}

func (n *InfixExpression) expressionNode()      {}
func (n *InfixExpression) TokenLiteral() string { return n.Token.Literal }
func (n *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(n.Left.String())
	out.WriteString(" ")
	out.WriteString(n.Operator)
	out.WriteString(" ")
	out.WriteString(n.Right.String())
	out.WriteString(")")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

var _ Statement = &ReturnStatement{}

func (n *ReturnStatement) statementNode()       {}
func (n *ReturnStatement) TokenLiteral() string { return n.Token.Literal }
func (n *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(n.TokenLiteral())
	out.WriteString(" ")
	if n.ReturnValue != nil {
		out.WriteString(n.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

var _ Statement = &ExpressionStatement{}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	// var out bytes.Buffer
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
	// return out.String()
}
