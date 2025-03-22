package quark

import (
	"fmt"

	"github.com/janqx/quark-lang/v1/ast"
	"github.com/janqx/quark-lang/v1/tokenize"
)

type loopState struct {
	continues []int
	breaks    []int
}

type Compiler struct {
	ctx                *Context
	parent             *Compiler
	loops              []*loopState
	loopIndex          int
	currentSymbolTable *SymbolTable
	currentFunction    *CompiledFunctionObject
	compiled           *compiled
}

func NewCompiler(ctx *Context, parent *Compiler) *Compiler {
	return &Compiler{
		ctx:    ctx,
		parent: parent,
	}
}

func (c *Compiler) mark() int {
	return len(c.currentFunction.Instructions)
}

func (c *Compiler) emit1(opcode Opcode) {
	c.emit2(opcode, InvalidOperand)
}

func (c *Compiler) emit2(opcode Opcode, operand Operand) {
	c.addInstruction(NewInstruction(opcode, operand))
}
func (c *Compiler) addInstruction(inst Instruction) {
	c.currentFunction.Instructions = append(c.currentFunction.Instructions, inst)
}

func (c *Compiler) setInstructionOperand(index int, operand Operand) {
	c.currentFunction.Instructions[index] = NewInstruction(c.currentFunction.Instructions[index].Opcode(), operand)
}

func (c *Compiler) addContinueMark(mark int) {
	c.loops[c.loopIndex].continues = append(c.loops[c.loopIndex].continues, mark)
}

func (c *Compiler) addBreakMark(mark int) {
	c.loops[c.loopIndex].breaks = append(c.loops[c.loopIndex].breaks, mark)
}

func (c *Compiler) pushLoopState() {
	c.loopIndex++
	c.loops[c.loopIndex] = &loopState{}
}

func (c *Compiler) popLoopState() {
	c.loopIndex--
}

func (c *Compiler) Compile(chunk *ast.Chunk) (*compiled, error) {
	result := c.compileChunk(chunk)
	return result, c.ctx.err
}

func (c *Compiler) compileChunk(chunk *ast.Chunk) *compiled {
	c.loops = make([]*loopState, 32)
	c.loopIndex = -1

	if c.ctx.Mode == ModeREPL {
		c.currentSymbolTable = c.ctx.globalSymbolTable
	} else {
		c.currentSymbolTable = c.ctx.globalSymbolTable.Push(TypeFunction)
	}

	c.currentFunction = &CompiledFunctionObject{
		Name:           "<compiled-function entry>",
		Instructions:   make([]Instruction, 0),
		ParameterNames: []string{},
		SymbolTable:    c.currentSymbolTable,
	}

	c.compiled = &compiled{
		entryFunction:     c.currentFunction,
		compiledFunctions: []*CompiledFunctionObject{c.currentFunction},
	}

	chunk.Accept(c)

	return c.compiled
}

func (c *Compiler) VisitChunk(node *ast.Chunk) {
	node.Statements.Accept(c)
	c.emit1(OpLoadNull)
	c.emit1(OpReturn)
}

func (c *Compiler) VisitStatementList(node *ast.StatementList) {
	for _, s := range node.List {
		s.Accept(c)
	}
}

func (c *Compiler) VisitContinueStatement(node *ast.ContinueStatement) {
	c.addContinueMark(c.mark())
	c.emit1(OpJump)
}

func (c *Compiler) VisitBreakStatement(node *ast.BreakStatement) {
	c.addBreakMark(c.mark())
	c.emit1(OpJump)
}

func (c *Compiler) VisitBlockStatement(node *ast.BlockStatement) {
	c.currentSymbolTable = c.currentSymbolTable.Push(TypeBlock)
	node.Statements.Accept(c)
	c.currentSymbolTable = c.currentSymbolTable.Pop()
}

func (c *Compiler) VisitReturnStatement(node *ast.ReturnStatement) {
	node.Expressions.Accept(c)
	c.emit2(OpReturn, Operand(node.Expressions.Count()))
}

func (c *Compiler) VisitIfStatement(node *ast.IfStatement) {
	jumpNextMark := -1
	jumpElseMark := -1
	quitIfMarks := make([]int, 0)

	node.Condition.Accept(c)

	jumpNextMark = c.mark()
	c.emit2(OpJumpIfFalse, InvalidOperand)

	node.ThenBody.Accept(c)

	quitIfMarks = append(quitIfMarks, c.mark())
	c.emit2(OpJump, InvalidOperand)

	if len(node.Elifs) > 0 {
		for i, elif := range node.Elifs {
			c.setInstructionOperand(jumpNextMark, Operand(c.mark()))
			elif.Condition.Accept(c)

			jumpNextMark = c.mark()
			c.emit2(OpJumpIfFalse, InvalidOperand)

			if i == len(node.Elifs)-1 {
				jumpElseMark = jumpNextMark
			}

			elif.Body.Accept(c)

			quitIfMarks = append(quitIfMarks, c.mark())
			c.emit2(OpJump, InvalidOperand)
		}
	} else {
		jumpElseMark = jumpNextMark
	}

	c.setInstructionOperand(jumpElseMark, Operand(c.mark()))
	if node.ElseBody != nil {
		node.ElseBody.Accept(c)
	}

	for _, mark := range quitIfMarks {
		c.setInstructionOperand(mark, Operand(c.mark()))
	}

	c.emit1(OpNop)
}

func (c *Compiler) VisitForStatement(node *ast.ForStatement) {
	c.pushLoopState()

	if node.Init != nil {
		node.Init.Accept(c)
	}

	startLoopMark := c.mark()

	if node.Condition != nil {
		node.Condition.Accept(c)
	} else {
		c.emit1(OpLoadTrue)
	}

	c.addBreakMark(c.mark())
	c.emit1(OpJumpIfFalse)

	node.Body.Accept(c)

	if node.Increment != nil {
		node.Increment.Accept(c)
	}

	c.emit2(OpJump, Operand(startLoopMark))

	quitLoopMark := c.mark()
	c.emit1(OpNop)

	for _, mark := range c.loops[c.loopIndex].continues {
		c.setInstructionOperand(mark, Operand(startLoopMark))
	}

	for _, mark := range c.loops[c.loopIndex].breaks {
		c.setInstructionOperand(mark, Operand(quitLoopMark))
	}

	c.popLoopState()
}

func (c *Compiler) VisitFunctionDeclareStatement(node *ast.FunctionDeclareStatement) {
	symbol := c.currentSymbolTable.FindSymbol(node.Name)
	if symbol == nil {
		symbol = c.currentSymbolTable.AddLocalSymbol(node.Name)
	}

	name := node.Name
	prev := c.currentFunction

	c.currentSymbolTable = c.currentSymbolTable.Push(TypeFunction)

	fn := &CompiledFunctionObject{}
	fn.Name = name
	fn.Instructions = make([]Instruction, 0)
	fn.ParameterNames = node.ParameterNames[:]
	fn.SymbolTable = c.currentSymbolTable

	c.compiled.compiledFunctions = append(c.compiled.compiledFunctions, fn)
	c.currentFunction = fn
	for _, name := range fn.ParameterNames {
		c.currentSymbolTable.AddLocalSymbol(name)
	}
	node.Body.Accept(c)
	c.emit1(OpLoadNull)
	c.emit1(OpReturn)

	c.currentSymbolTable = c.currentSymbolTable.Pop()
	c.currentFunction = prev
	c.emit2(OpLoadConst, Operand(c.ctx.appendConstant(fn)))
	c.emit1(OpClosure)

	switch symbol.Scope {
	case ScopeLocal:
		c.emit2(OpStoreLocal, Operand(symbol.Index))
	case ScopeOuter:
		c.emit2(OpStoreOuter, Operand(symbol.Index))
	case ScopeGlobal:
		c.emit2(OpStoreGlobal, Operand(symbol.Index))
	}
}

func (c *Compiler) VisitImportStatement(node *ast.ImportStatement) {
	c.emit2(OpImport, Operand(c.ctx.addStringConstant(node.Path)))
}

func (c *Compiler) VisitExportStatement(node *ast.ExportStatement) {
	node.Module.Accept(c)
	c.emit1(OpExport)
}

func (c *Compiler) VisitAssignStatement(node *ast.AssignStatement) {
	node.Expressions.Accept(c)
	for i := len(node.Assignables) - 1; i >= 0; i-- {
		node.Assignables[i].Accept(c)
	}
}

func (c *Compiler) VisitCallFunctionStatement(node *ast.CallFunctionStatement) {
	node.Args.Accept(c)
	node.Callable.Accept(c)
	c.emit2(OpCall, Operand(node.Args.Count()))
	c.emit1(OpRemoveTop)
}

func (c *Compiler) VisitExpressionStatement(node *ast.ExpressionStatement) {
	node.Expression.Accept(c)
	c.emit1(OpRemoveTop)
}

func (c *Compiler) VisitDebuggerStatement(node *ast.DebuggerStatement) {
	c.emit1(OpDebugger)
}

func (c *Compiler) VisitExpressionList(node *ast.ExpressionList) {
	for _, e := range node.List {
		e.Accept(c)
	}
}

func (c *Compiler) VisitNullLiteralExpression(node *ast.NullLiteralExpression) {
	c.emit1(OpLoadNull)
}

func (c *Compiler) VisitTrueLiteralExpression(node *ast.TrueLiteralExpression) {
	c.emit1(OpLoadTrue)

}

func (c *Compiler) VisitFalseLiteralExpression(node *ast.FalseLiteralExpression) {
	c.emit1(OpLoadFalse)
}

func (c *Compiler) VisitIntLiteralExpression(node *ast.IntLiteralExpression) {
	c.emit2(OpLoadConst, Operand(c.ctx.addIntConstant(node.Value)))
}

func (c *Compiler) VisitFloatLiteralExpression(node *ast.FloatLiteralExpression) {
	c.emit2(OpLoadConst, Operand(c.ctx.addFloatConstant(node.Value)))
}

func (c *Compiler) VisitStringLiteralExpression(node *ast.StringLiteralExpression) {
	c.emit2(OpLoadConst, Operand(c.ctx.addStringConstant(node.Value)))
}

func (c *Compiler) VisitListLiteralExpression(node *ast.ListLiteralExpression) {
	node.Value.Accept(c)
	c.emit2(OpBuildList, Operand(node.Value.Count()))
}

func (c *Compiler) VisitDictLiteralExpression(node *ast.DictLiteralExpression) {
	for k, v := range node.Value {
		c.emit2(OpLoadConst, Operand(c.ctx.addStringConstant(k)))
		v.Accept(c)
	}
	c.emit2(OpBuildDict, Operand(len(node.Value)))
}

func (c *Compiler) VisitIdentifierExpression(node *ast.IdentifierExpression) {
	symbol := c.currentSymbolTable.FindSymbol(node.Name)
	if symbol == nil {
		if !node.Assign {
			panic(fmt.Errorf("undeclared identifier: '%s'", node.Name))
		}
		if c.currentSymbolTable.Parent == nil {
			symbol = c.ctx.addGlobalSymbol(node.Name)
		} else {
			symbol = c.currentSymbolTable.AddLocalSymbol(node.Name)
		}
	}
	switch symbol.Scope {
	case ScopeLocal:
		if node.Assign {
			c.emit2(OpStoreLocal, Operand(symbol.Index))
		} else {
			c.emit2(OpLoadLocal, Operand(symbol.Index))
		}
	case ScopeOuter:
		if node.Assign {
			c.emit2(OpStoreOuter, Operand(symbol.Index))
		} else {
			c.emit2(OpLoadOuter, Operand(symbol.Index))
		}
	case ScopeGlobal:
		if node.Assign {
			c.emit2(OpStoreGlobal, Operand(symbol.Index))
		} else {
			c.emit2(OpLoadGlobal, Operand(symbol.Index))
		}
	}
}

func (c *Compiler) VisitIndexAccessExpression(node *ast.IndexAccessExpression) {
	node.Value.Accept(c)
	node.Index.Accept(c)
	if node.Assign {
		c.emit1(OpStoreIndex)
	} else {
		c.emit1(OpLoadIndex)
	}
}

func (c *Compiler) VisitAttributeAccessExpression(node *ast.AttributeAccessExpression) {
	node.Value.Accept(c)
	c.emit2(OpLoadConst, Operand(c.ctx.addStringConstant(node.Name)))
	if node.Assign {
		c.emit1(OpStoreAttribute)
	} else {
		c.emit1(OpLoadAttribute)
	}
}

func (c *Compiler) VisitFunctionDeclareExpression(node *ast.FunctionDeclareExpression) {
	name := fmt.Sprintf("<closure #%d>", len(c.compiled.compiledFunctions))
	prev := c.currentFunction

	c.currentSymbolTable = c.currentSymbolTable.Push(TypeFunction)

	fn := &CompiledFunctionObject{}
	fn.Name = name
	fn.Instructions = make([]Instruction, 0)
	fn.ParameterNames = node.ParameterNames[:]
	fn.SymbolTable = c.currentSymbolTable

	c.compiled.compiledFunctions = append(c.compiled.compiledFunctions, fn)
	c.currentFunction = fn
	for _, name := range fn.ParameterNames {
		c.currentSymbolTable.AddLocalSymbol(name)
	}
	node.Body.Accept(c)
	c.emit1(OpLoadNull)
	c.emit1(OpReturn)

	c.currentSymbolTable = c.currentSymbolTable.Pop()
	c.currentFunction = prev
	c.emit2(OpLoadConst, Operand(c.ctx.appendConstant(fn)))
	c.emit1(OpClosure)
}

func (c *Compiler) VisitCallFunctionExpression(node *ast.CallFunctionExpression) {
	node.Args.Accept(c)
	node.Callable.Accept(c)
	c.emit2(OpCall, Operand(node.Args.Count()))
}

func (c *Compiler) VisitUnaryExpression(node *ast.UnaryExpression) {
	node.Expression.Accept(c)
	var op Opcode
	switch node.Op {
	case tokenize.TokenPlus:
		op = OpUnaryPlus
	case tokenize.TokenMinus:
		op = OpUnaryMinus
	case tokenize.TokenBitNot:
		op = OpUnaryBitNot
	case tokenize.TokenNot:
		op = OpUnaryNot
	default:
		panic(fmt.Errorf("invalid unary operator: %s", tokenize.TokenTypeToString[node.Op]))
	}
	c.emit1(op)
}

func (c *Compiler) VisitBinaryExpression(node *ast.BinaryExpression) {
	if node.Op == tokenize.TokenLogicAnd {
		node.Left.Accept(c)
		mark := c.mark()
		c.emit1(OpJumpIfFalseOrPop)
		node.Right.Accept(c)
		c.setInstructionOperand(mark, Operand(c.mark()))
		return
	} else if node.Op == tokenize.TokenLogicOr {
		node.Left.Accept(c)
		mark := c.mark()
		c.emit1(OpJumpIfTrueOrPop)
		node.Right.Accept(c)
		c.setInstructionOperand(mark, Operand(c.mark()))
		return
	}

	node.Left.Accept(c)
	node.Right.Accept(c)
	var op Opcode
	switch node.Op {
	case tokenize.TokenPlus:
		op = OpBinaryAdd
	case tokenize.TokenMinus:
		op = OpBinarySub
	case tokenize.TokenMul:
		op = OpBinaryMul
	case tokenize.TokenDiv:
		op = OpBinaryDiv
	case tokenize.TokenMod:
		op = OpBinaryMod
	case tokenize.TokenLT:
		op = OpBinaryLT
	case tokenize.TokenLTE:
		op = OpBinaryLTE
	case tokenize.TokenGT:
		op = OpBinaryGT
	case tokenize.TokenGTE:
		op = OpBinaryGTE
	case tokenize.TokenEQ:
		op = OpBinaryEQ
	case tokenize.TokenNEQ:
		op = OpBinaryNEQ
	case tokenize.TokenBitAnd:
		op = OpBinaryBitAnd
	case tokenize.TokenBitOr:
		op = OpBinaryBitOr
	case tokenize.TokenBitXor:
		op = OpBinaryBitXor
	case tokenize.TokenBitLhs:
		op = OpBinaryBitLhs
	case tokenize.TokenBitRhs:
		op = OpBinaryBitRhs
	default:
		panic(fmt.Errorf("invalid binary operator: %s", tokenize.TokenTypeToString[node.Op]))
	}
	c.emit1(op)
}

func (c *Compiler) VisitTernaryExpression(node *ast.TernaryExpression) {
	node.Cond.Accept(c)
	mark1 := c.mark()
	c.emit1(OpJumpIfFalse)
	node.X.Accept(c)
	mark2 := c.mark()
	c.emit1(OpJump)
	c.setInstructionOperand(mark1, Operand(c.mark()))
	node.Y.Accept(c)
	c.setInstructionOperand(mark2, Operand(c.mark()))
}
