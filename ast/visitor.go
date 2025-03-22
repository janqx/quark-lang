package ast

type Visitor interface {
	VisitChunk(node *Chunk)

	VisitStatementList(node *StatementList)
	VisitContinueStatement(node *ContinueStatement)
	VisitBreakStatement(node *BreakStatement)
	VisitBlockStatement(node *BlockStatement)
	VisitReturnStatement(node *ReturnStatement)
	VisitIfStatement(node *IfStatement)
	VisitForStatement(node *ForStatement)
	VisitFunctionDeclareStatement(node *FunctionDeclareStatement)
	VisitImportStatement(node *ImportStatement)
	VisitExportStatement(node *ExportStatement)
	VisitAssignStatement(node *AssignStatement)
	VisitCallFunctionStatement(node *CallFunctionStatement)
	VisitExpressionStatement(node *ExpressionStatement)
	VisitDebuggerStatement(node *DebuggerStatement)

	VisitExpressionList(node *ExpressionList)
	VisitNullLiteralExpression(node *NullLiteralExpression)
	VisitTrueLiteralExpression(node *TrueLiteralExpression)
	VisitFalseLiteralExpression(node *FalseLiteralExpression)
	VisitIntLiteralExpression(node *IntLiteralExpression)
	VisitFloatLiteralExpression(node *FloatLiteralExpression)
	VisitStringLiteralExpression(node *StringLiteralExpression)
	VisitListLiteralExpression(node *ListLiteralExpression)
	VisitDictLiteralExpression(node *DictLiteralExpression)
	VisitIdentifierExpression(node *IdentifierExpression)
	VisitIndexAccessExpression(node *IndexAccessExpression)
	VisitAttributeAccessExpression(node *AttributeAccessExpression)
	VisitFunctionDeclareExpression(node *FunctionDeclareExpression)
	VisitCallFunctionExpression(node *CallFunctionExpression)
	VisitUnaryExpression(node *UnaryExpression)
	VisitBinaryExpression(node *BinaryExpression)
	VisitTernaryExpression(node *TernaryExpression)
}

type EmptyVisitor struct{}

func (c *EmptyVisitor) VisitChunk(node *Chunk) {}

func (c *EmptyVisitor) VisitStatementList(node *StatementList)                       {}
func (c *EmptyVisitor) VisitContinueStatement(node *ContinueStatement)               {}
func (c *EmptyVisitor) VisitBreakStatement(node *BreakStatement)                     {}
func (c *EmptyVisitor) VisitBlockStatement(node *BlockStatement)                     {}
func (c *EmptyVisitor) VisitReturnStatement(node *ReturnStatement)                   {}
func (c *EmptyVisitor) VisitIfStatement(node *IfStatement)                           {}
func (c *EmptyVisitor) VisitForStatement(node *ForStatement)                         {}
func (c *EmptyVisitor) VisitFunctionDeclareStatement(node *FunctionDeclareStatement) {}
func (c *EmptyVisitor) VisitImportStatement(node *ImportStatement)                   {}
func (c *EmptyVisitor) VisitExportStatement(node *ExportStatement)                   {}
func (c *EmptyVisitor) VisitAssignStatement(node *AssignStatement)                   {}
func (c *EmptyVisitor) VisitCallFunctionStatement(node *CallFunctionStatement)       {}
func (c *EmptyVisitor) VisitExpressionStatement(node *ExpressionStatement)           {}
func (c *EmptyVisitor) VisitDebuggerStatement(node *DebuggerStatement)               {}

func (c *EmptyVisitor) VisitExpressionList(node *ExpressionList)                       {}
func (c *EmptyVisitor) VisitNullLiteralExpression(node *NullLiteralExpression)         {}
func (c *EmptyVisitor) VisitTrueLiteralExpression(node *TrueLiteralExpression)         {}
func (c *EmptyVisitor) VisitFalseLiteralExpression(node *FalseLiteralExpression)       {}
func (c *EmptyVisitor) VisitIntLiteralExpression(node *IntLiteralExpression)           {}
func (c *EmptyVisitor) VisitFloatLiteralExpression(node *FloatLiteralExpression)       {}
func (c *EmptyVisitor) VisitStringLiteralExpression(node *StringLiteralExpression)     {}
func (c *EmptyVisitor) VisitListLiteralExpression(node *ListLiteralExpression)         {}
func (c *EmptyVisitor) VisitDictLiteralExpression(node *DictLiteralExpression)         {}
func (c *EmptyVisitor) VisitIdentifierExpression(node *IdentifierExpression)           {}
func (c *EmptyVisitor) VisitIndexAccessExpression(node *IndexAccessExpression)         {}
func (c *EmptyVisitor) VisitAttributeAccessExpression(node *AttributeAccessExpression) {}
func (c *EmptyVisitor) VisitFunctionDeclareExpression(node *FunctionDeclareExpression) {}
func (c *EmptyVisitor) VisitCallFunctionExpression(node *CallFunctionExpression)       {}
func (c *EmptyVisitor) VisitUnaryExpression(node *UnaryExpression)                     {}
func (c *EmptyVisitor) VisitBinaryExpression(node *BinaryExpression)                   {}
func (c *EmptyVisitor) VisitTernaryExpression(node *TernaryExpression)                 {}
