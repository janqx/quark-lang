package ast

import (
	"fmt"
	"strings"

	"github.com/janqx/quark-lang/v1/tokenize"
)

type Statement interface {
	Node
	statementNode()
}

type StatementImpl struct {
	Statement
	start tokenize.Position
	end   tokenize.Position
}

func (node *StatementImpl) Start() tokenize.Position {
	return node.start
}

func (node *StatementImpl) End() tokenize.Position {
	return node.end
}

func (node *StatementImpl) String() string {
	panic(fmt.Errorf("not implemented"))
}

func (node *StatementImpl) Accept(visitor Visitor) {
	panic(fmt.Errorf("not implemented"))
}

func (node *StatementImpl) statementNode() {
	panic(fmt.Errorf("not implemented"))
}

var SingletonEmptyStatement = &EmptyStatement{}

type EmptyStatement struct {
	StatementImpl
}

func (node *EmptyStatement) String() string {
	return ""
}

func (node *EmptyStatement) Accept(visitor Visitor) {
	// do nothing
}

type Chunk struct {
	StatementImpl
	Statements *StatementList
}

func (node *Chunk) String() string {
	return node.Statements.String()
}

func (node *Chunk) Accept(visitor Visitor) {
	visitor.VisitChunk(node)
}

var EmptyStatementList = &StatementList{List: []Statement{}}

type StatementList struct {
	StatementImpl
	List []Statement
}

func (node *StatementList) Count() int {
	return len(node.List)
}

func (node *StatementList) String() string {
	ss := make([]string, 0)
	for _, s := range node.List {
		ss = append(ss, s.String())
	}
	return strings.Join(ss, "\n")
}

func (node *StatementList) Accept(visitor Visitor) {
	visitor.VisitStatementList(node)
}

type ContinueStatement struct {
	StatementImpl
}

func (node *ContinueStatement) String() string {
	return "continue"
}

func (node *ContinueStatement) Accept(visitor Visitor) {
	visitor.VisitContinueStatement(node)
}

type BreakStatement struct {
	StatementImpl
}

func (node *BreakStatement) String() string {
	return "break"
}

func (node *BreakStatement) Accept(visitor Visitor) {
	visitor.VisitBreakStatement(node)
}

type BlockStatement struct {
	StatementImpl
	Statements *StatementList
}

func (node *BlockStatement) String() string {
	return "{" + node.Statements.String() + "}"
}

func (node *BlockStatement) Accept(visitor Visitor) {
	visitor.VisitBlockStatement(node)
}

type ReturnStatement struct {
	StatementImpl
	Expressions *ExpressionList
}

func (node *ReturnStatement) String() string {
	return "return " + node.Expressions.String()
}

func (node *ReturnStatement) Accept(visitor Visitor) {
	visitor.VisitReturnStatement(node)
}

type Elif struct {
	Condition Expression
	Body      Statement
}

type IfStatement struct {
	StatementImpl
	Condition Expression
	ThenBody  Statement
	Elifs     []Elif
	ElseBody  Statement
}

func (node *IfStatement) String() string {
	result := "if " + node.Condition.String() + " " + node.ThenBody.String()
	for _, elif := range node.Elifs {
		result += "else if " + elif.Condition.String() + " " + elif.Body.String()
	}
	if node.ElseBody != nil {
		result += "else " + node.ElseBody.String()
	}
	return result
}

func (node *IfStatement) Accept(visitor Visitor) {
	visitor.VisitIfStatement(node)
}

type ForStatement struct {
	StatementImpl
	Init      Statement
	Condition Expression
	Increment Statement
	Body      Statement
}

func (node *ForStatement) String() string {
	result := "for "
	if node.Init != nil {
		result += node.Init.String()
	}
	result += ";"
	if node.Condition != nil {
		result += node.Condition.String()
	}
	result += ";"
	if node.Increment != nil {
		result += node.Increment.String()
	}
	result += node.Body.String()
	return result
}

func (node *ForStatement) Accept(visitor Visitor) {
	visitor.VisitForStatement(node)
}

type FunctionDeclareStatement struct {
	StatementImpl
	Name           string
	ParameterNames []string
	Body           Statement
}

func (node *FunctionDeclareStatement) String() string {
	return "fn " + node.Name + "(" + strings.Join(node.ParameterNames, ",") + ")" + node.Body.String()
}

func (node *FunctionDeclareStatement) Accept(visitor Visitor) {
	visitor.VisitFunctionDeclareStatement(node)
}

type ImportStatement struct {
	StatementImpl
	Path string
	Name string
}

func (node *ImportStatement) String() string {
	return "import \"" + node.Path + "\" as " + node.Name
}

func (node *ImportStatement) Accept(visitor Visitor) {
	visitor.VisitImportStatement(node)
}

type ExportStatement struct {
	StatementImpl
	Module Expression
}

func (node *ExportStatement) String() string {
	return "export " + node.Module.String()
}

func (node *ExportStatement) Accept(visitor Visitor) {
	visitor.VisitExportStatement(node)
}

type AssignStatement struct {
	StatementImpl
	Assignables []Expression
	Expressions *ExpressionList
}

func (node *AssignStatement) String() string {
	assigns := make([]string, 0)
	for _, assign := range node.Assignables {
		assigns = append(assigns, assign.String())
	}
	return strings.Join(assigns, ",") + "=" + node.Expressions.String()
}

func (node *AssignStatement) Accept(visitor Visitor) {
	visitor.VisitAssignStatement(node)
}

type CallFunctionStatement struct {
	StatementImpl
	Callable Expression
	Args     *ExpressionList
}

func (node *CallFunctionStatement) String() string {
	return node.Callable.String() + "(" + node.Args.String() + ")"
}

func (node *CallFunctionStatement) Accept(visitor Visitor) {
	visitor.VisitCallFunctionStatement(node)
}

type ExpressionStatement struct {
	StatementImpl
	Expression Expression
}

func (node *ExpressionStatement) String() string {
	return node.Expression.String()
}

func (node *ExpressionStatement) Accept(visitor Visitor) {
	visitor.VisitExpressionStatement(node)
}

type DebuggerStatement struct {
	StatementImpl
}

func (node *DebuggerStatement) String() string {
	return "debugger"
}

func (node *DebuggerStatement) Accept(visitor Visitor) {
	visitor.VisitDebuggerStatement(node)
}
