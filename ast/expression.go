package ast

import (
	"fmt"
	"strings"

	"github.com/janqx/quark-lang/v1/tokenize"
)

type Expression interface {
	Node
	expressionNode()
}

type ExpressionImpl struct {
	Expression
	start tokenize.Position
	end   tokenize.Position
}

func (node *ExpressionImpl) Start() tokenize.Position {
	return node.start
}

func (node *ExpressionImpl) End() tokenize.Position {
	return node.end
}

func (node *ExpressionImpl) String() string {
	panic(fmt.Errorf("not implemented"))
}

func (node *ExpressionImpl) Accept(visitor Visitor) {
	panic(fmt.Errorf("not implemented"))
}

func (node *ExpressionImpl) expressionNode() {
	panic(fmt.Errorf("not implemented"))
}

var EmptyExpressionList = &ExpressionList{List: []Expression{}}

type ExpressionList struct {
	ExpressionImpl
	List []Expression
}

func (node *ExpressionList) Count() int {
	return len(node.List)
}

func (node *ExpressionList) String() string {
	es := make([]string, 0)
	for _, e := range node.List {
		es = append(es, e.String())
	}
	return strings.Join(es, ",")
}

func (node *ExpressionList) Accept(visitor Visitor) {
	visitor.VisitExpressionList(node)
}

type NullLiteralExpression struct {
	ExpressionImpl
}

func (node *NullLiteralExpression) String() string {
	return "null"
}

func (node *NullLiteralExpression) Accept(visitor Visitor) {
	visitor.VisitNullLiteralExpression(node)
}

type TrueLiteralExpression struct {
	ExpressionImpl
}

func (node *TrueLiteralExpression) String() string {
	return "true"
}

func (node *TrueLiteralExpression) Accept(visitor Visitor) {
	visitor.VisitTrueLiteralExpression(node)
}

type FalseLiteralExpression struct {
	ExpressionImpl
}

func (node *FalseLiteralExpression) String() string {
	return "false"
}

func (node *FalseLiteralExpression) Accept(visitor Visitor) {
	visitor.VisitFalseLiteralExpression(node)
}

type IntLiteralExpression struct {
	ExpressionImpl
	Value int64
}

func (node *IntLiteralExpression) String() string {
	return fmt.Sprintf("%d", node.Value)
}

func (node *IntLiteralExpression) Accept(visitor Visitor) {
	visitor.VisitIntLiteralExpression(node)
}

type FloatLiteralExpression struct {
	ExpressionImpl
	Value float64
}

func (node *FloatLiteralExpression) String() string {
	return fmt.Sprintf("%f", node.Value)
}

func (node *FloatLiteralExpression) Accept(visitor Visitor) {
	visitor.VisitFloatLiteralExpression(node)
}

type StringLiteralExpression struct {
	ExpressionImpl
	Value string
	Proto string
}

func (node *StringLiteralExpression) String() string {
	return node.Proto
}

func (node *StringLiteralExpression) Accept(visitor Visitor) {
	visitor.VisitStringLiteralExpression(node)
}

type ListLiteralExpression struct {
	ExpressionImpl
	Value *ExpressionList
}

func (node *ListLiteralExpression) String() string {
	return "[" + node.Value.String() + "]"
}

func (node *ListLiteralExpression) Accept(visitor Visitor) {
	visitor.VisitListLiteralExpression(node)
}

type DictLiteralExpression struct {
	ExpressionImpl
	Value map[string]Expression
}

func (node *DictLiteralExpression) String() string {
	s := "{"
	if node.Value != nil {
		i := 0
		for key, value := range node.Value {
			s += key + ":" + value.String()
			if i < len(node.Value)-1 {
				s += ","
			}
			i++
		}
	}
	return s + "}"
}

func (node *DictLiteralExpression) Accept(visitor Visitor) {
	visitor.VisitDictLiteralExpression(node)
}

type IdentifierExpression struct {
	ExpressionImpl
	Name   string
	Assign bool
}

func (node *IdentifierExpression) String() string {
	return node.Name
}

func (node *IdentifierExpression) Accept(visitor Visitor) {
	visitor.VisitIdentifierExpression(node)
}

type IndexAccessExpression struct {
	ExpressionImpl
	Value  Expression
	Index  Expression
	Assign bool
}

func (node *IndexAccessExpression) String() string {
	return node.Value.String() + "[" + node.Index.String() + "]"
}

func (node *IndexAccessExpression) Accept(visitor Visitor) {
	visitor.VisitIndexAccessExpression(node)
}

type AttributeAccessExpression struct {
	ExpressionImpl
	Value  Expression
	Name   string
	Assign bool
}

func (node *AttributeAccessExpression) String() string {
	return node.Value.String() + "." + node.Name
}

func (node *AttributeAccessExpression) Accept(visitor Visitor) {
	visitor.VisitAttributeAccessExpression(node)
}

type FunctionDeclareExpression struct {
	ExpressionImpl
	ParameterNames []string
	Body           Statement
}

func (node *FunctionDeclareExpression) String() string {
	s := "fn("
	if node.ParameterNames != nil {
		s += strings.Join(node.ParameterNames, ",")
	}
	s += ")" + node.Body.String()

	return s
}

func (node *FunctionDeclareExpression) Accept(visitor Visitor) {
	visitor.VisitFunctionDeclareExpression(node)
}

type CallFunctionExpression struct {
	ExpressionImpl
	Callable Expression
	Args     *ExpressionList
}

func (node *CallFunctionExpression) String() string {
	return node.Callable.String() + "(" + node.Args.String() + ")"
}

func (node *CallFunctionExpression) Accept(visitor Visitor) {
	visitor.VisitCallFunctionExpression(node)
}

type UnaryExpression struct {
	ExpressionImpl
	Op         tokenize.TokenType
	Expression Expression
}

func (node *UnaryExpression) String() string {
	return node.Op.String() + node.Expression.String()
}

func (node *UnaryExpression) Accept(visitor Visitor) {
	visitor.VisitUnaryExpression(node)
}

type BinaryExpression struct {
	ExpressionImpl
	Op    tokenize.TokenType
	Left  Expression
	Right Expression
}

func (node *BinaryExpression) String() string {
	return node.Left.String() + node.Op.String() + node.Right.String()
}

func (node *BinaryExpression) Accept(visitor Visitor) {
	visitor.VisitBinaryExpression(node)
}

type TernaryExpression struct {
	ExpressionImpl
	Cond Expression
	X    Expression
	Y    Expression
}

func (node *TernaryExpression) String() string {
	return node.Cond.String() + "?" + node.X.String() + ":" + node.Y.String()
}

func (node *TernaryExpression) Accept(visitor Visitor) {
	visitor.VisitTernaryExpression(node)
}
