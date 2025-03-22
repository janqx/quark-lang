package parser

import (
	"fmt"
	"strings"

	"github.com/janqx/quark-lang/v1/ast"
	"github.com/janqx/quark-lang/v1/tokenize"
)

type Parser struct {
	filename string
	source   []byte
	lexer    *Lexer
	token    *tokenize.Token
	err      error
}

func NewParser(filename string, source []byte) *Parser {
	return &Parser{
		filename: filename,
		source:   source,
	}
}

func (p *Parser) Parse() (*ast.Chunk, error) {
	chunk := p.parse()
	if p.err != nil {
		return nil, fmt.Errorf("in file \"%s\"\n%s", p.filename, p.err.Error())
	}
	return chunk, nil
}

func (p *Parser) parse() *ast.Chunk {
	p.lexer = NewLexer(p.filename, strings.NewReader(string(p.source)))
	p.token = nil
	p.err = nil
	p.next()
	chunk := &ast.Chunk{}
	chunk.Statements = p.parseStatementList()
	p.expect(tokenize.TokenEof)
	return chunk
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.token.Type {
	case tokenize.TokenDebugger:
		p.next()
		return &ast.DebuggerStatement{}
	case tokenize.TokenContinue:
		p.next()
		return &ast.ContinueStatement{}
	case tokenize.TokenBreak:
		p.next()
		return &ast.BreakStatement{}
	case tokenize.TokenReturn:
		return p.parseReturnStatement()
	case tokenize.TokenIf:
		return p.parseIfStatement()
	case tokenize.TokenFor:
		return p.parseForStatement()
	case tokenize.TokenFunction:
		return p.parseFunctionDeclareStatement()
	case tokenize.TokenOpenBrace:
		return p.parseBlockStatement()
	// case tokenize.TokenImport:
	case tokenize.TokenExport:
		return p.parseExportStatement()
	case tokenize.TokenNewline:
		p.next()
		return ast.SingletonEmptyStatement
	case tokenize.TokenEof:
		return ast.SingletonEmptyStatement
	default:
		// assign or call?
		return p.parseOtherStatement()
	}
}

// statement+
func (p *Parser) parseStatementList() *ast.StatementList {
	result := &ast.StatementList{
		List: []ast.Statement{p.parseStatement()},
	}
	for {
		if p.test(tokenize.TokenEof, tokenize.TokenCloseBrace) {
			break
		}
		result.List = append(result.List, p.parseStatement())
	}
	return result
}

func (p *Parser) parseReturnStatement() ast.Statement {
	p.expect(tokenize.TokenReturn)
	var expressions *ast.ExpressionList = nil
	if p.test(tokenize.TokenCloseBrace) || p.token.IsNewLine() {
		expressions = &ast.ExpressionList{
			List: []ast.Expression{&ast.NullLiteralExpression{}},
		}
	} else {
		expressions = p.parseExpressionList(false)
	}
	return &ast.ReturnStatement{Expressions: expressions}
}

/*
if cond {
} else if cond {
} else if cond {
} else {
}
*/
func (p *Parser) parseIfStatement() ast.Statement {
	p.expect(tokenize.TokenIf)
	result := &ast.IfStatement{}
	result.Condition = p.parseExpression()
	result.ThenBody = p.parseBlockStatement()
	result.Elifs = make([]ast.Elif, 0)
	for p.test(tokenize.TokenElse) {
		p.next()
		if p.test(tokenize.TokenIf) {
			p.next()
			elif := ast.Elif{}
			elif.Condition = p.parseExpression()
			elif.Body = p.parseBlockStatement()
			result.Elifs = append(result.Elifs, elif)
		} else {
			result.ElseBody = p.parseBlockStatement()
			break
		}
	}
	return result
}

/*
for { }
for cond { }
for init?; cond?; post? { }
*/
func (p *Parser) parseForStatement() ast.Statement {
	p.expect(tokenize.TokenFor)
	result := &ast.ForStatement{}

	var x ast.Statement = nil

	if p.test(tokenize.TokenOpenBrace) {
		goto L_body
	}

	if !p.test(tokenize.TokenSemiColon) {
		x = p.parseStatement()

		if p.test(tokenize.TokenOpenBrace) {
			switch x := x.(type) {
			case *ast.ExpressionStatement:
				result.Condition = x.Expression
				goto L_body
			default:
				p.errorMessage("for syntax error, expected: ast.ExpressionStatement, but got: %T", x)
			}
		} else {
			result.Init = x
		}
	}

	p.expect(tokenize.TokenSemiColon)

	if !p.test(tokenize.TokenSemiColon) {
		result.Condition = p.parseExpression()
	}

	p.expect(tokenize.TokenSemiColon)

	if !p.test(tokenize.TokenOpenBrace) {
		result.Increment = p.parseStatement()
	}

L_body:
	result.Body = p.parseBlockStatement()
	return result
}

func (p *Parser) parseFunctionDeclareStatement() ast.Statement {
	p.expect(tokenize.TokenFunction)
	result := &ast.FunctionDeclareStatement{}
	result.Name = p.expect(tokenize.TokenIdentifier).Value.(string)
	p.expect(tokenize.TokenOpenParen)
	if !p.test(tokenize.TokenCloseParen) {
		result.ParameterNames = p.parseIdentifierList()
	}
	p.expect(tokenize.TokenCloseParen)
	result.Body = p.parseBlockStatement()
	return result
}

func (p *Parser) parseBlockStatement() ast.Statement {
	p.expect(tokenize.TokenOpenBrace)
	result := &ast.BlockStatement{Statements: ast.EmptyStatementList}
	if !p.test(tokenize.TokenCloseBrace) {
		result.Statements = p.parseStatementList()
	}
	p.expect(tokenize.TokenCloseBrace)
	return result
}

// export expression
func (p *Parser) parseExportStatement() ast.Statement {
	p.expect(tokenize.TokenExport)
	result := &ast.ExportStatement{
		Module: p.parseExpression(),
	}
	return result
}

var __cast_assignable = func(expression ast.Expression) ast.Expression {
	switch expression := expression.(type) {
	case *ast.IdentifierExpression:
		return &ast.IdentifierExpression{
			Assign: true,
			Name:   expression.Name,
		}
	case *ast.IndexAccessExpression:
		return &ast.IndexAccessExpression{
			Assign: true,
			Value:  expression.Value,
			Index:  expression.Index,
		}
	case *ast.AttributeAccessExpression:
		return &ast.AttributeAccessExpression{
			Assign: true,
			Value:  expression.Value,
			Name:   expression.Name,
		}
	}
	return expression
}

// assign or call
// a, b, c... = 1, 2, 3...
// func(a, b, c...)
func (p *Parser) parseOtherStatement() ast.Statement {
	expression := p.parseExpression()

	// assignable
	if p.test(tokenize.TokenComma, tokenize.TokenAssign) {
		assign := &ast.AssignStatement{
			Assignables: make([]ast.Expression, 0),
			Expressions: ast.EmptyExpressionList,
		}

		// parse more assignables
		assign.Assignables = append(assign.Assignables, __cast_assignable(expression))
		for p.test(tokenize.TokenComma) {
			p.next()
			assign.Assignables = append(assign.Assignables, __cast_assignable(p.parseExpression()))
		}
		p.expect(tokenize.TokenAssign)

		// parse expressions
		assign.Expressions = p.parseExpressionList(false)

		return assign
	}

	return &ast.ExpressionStatement{
		Expression: expression,
	}
}

func (p *Parser) parseExpression() ast.Expression {
	return p.parseTernaryExpression()
}

// cond ? x : y
func (p *Parser) parseTernaryExpression() ast.Expression {
	cond := p.parseLogicOrExpression()
	if p.test(tokenize.TokenQuestion) {
		p.next()
		x := p.parseTernaryExpression()

		p.expect(tokenize.TokenColon)
		y := p.parseTernaryExpression()

		return &ast.TernaryExpression{
			Cond: cond,
			X:    x,
			Y:    y,
		}
	}
	return cond
}

func (p *Parser) parseLogicOrExpression() ast.Expression {
	left := p.parseLogicAndExpression()
	for p.test(tokenize.TokenLogicOr) {
		p.next()
		right := p.parseLogicAndExpression()
		left = &ast.BinaryExpression{
			Op:    tokenize.TokenLogicOr,
			Left:  left,
			Right: right,
		}
	}
	return left
}

// x && y
func (p *Parser) parseLogicAndExpression() ast.Expression {
	left := p.parseBitOrExpression()
	for p.test(tokenize.TokenLogicAnd) {
		p.next()
		right := p.parseBitOrExpression()
		left = &ast.BinaryExpression{
			Op:    tokenize.TokenLogicAnd,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *Parser) parseBitOrExpression() ast.Expression {
	left := p.parseBitXorExpression()
	for p.test(tokenize.TokenBitOr) {
		p.next()
		right := p.parseBitXorExpression()
		left = &ast.BinaryExpression{
			Op:    tokenize.TokenBitOr,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *Parser) parseBitXorExpression() ast.Expression {
	left := p.parseBitAndExpression()
	for p.test(tokenize.TokenBitXor) {
		p.next()
		right := p.parseBitAndExpression()
		left = &ast.BinaryExpression{
			Op:    tokenize.TokenBitXor,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *Parser) parseBitAndExpression() ast.Expression {
	left := p.parseEqualityExpression()
	for p.test(tokenize.TokenBitAnd) {
		p.next()
		right := p.parseEqualityExpression()
		left = &ast.BinaryExpression{
			Op:    tokenize.TokenBitAnd,
			Left:  left,
			Right: right,
		}
	}
	return left
}

// x op=(== | !=) y
func (p *Parser) parseEqualityExpression() ast.Expression {
	left := p.parseRelationalExpression()
	for {
		op := p.token
		if !p.test(tokenize.TokenEQ, tokenize.TokenNEQ) {
			break
		}
		p.next()
		right := p.parseRelationalExpression()
		left = &ast.BinaryExpression{
			Op:    op.Type,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *Parser) parseRelationalExpression() ast.Expression {
	left := p.parseBitShiftExpression()
	for {
		op := p.token
		if !p.test(tokenize.TokenLT, tokenize.TokenLTE, tokenize.TokenGT, tokenize.TokenGTE) {
			break
		}
		p.next()
		right := p.parseBitShiftExpression()
		left = &ast.BinaryExpression{
			Op:    op.Type,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *Parser) parseBitShiftExpression() ast.Expression {
	left := p.parseAdditiveExpression()
	for {
		op := p.token
		if !p.test(tokenize.TokenBitLhs, tokenize.TokenBitRhs) {
			break
		}
		p.next()
		right := p.parseAdditiveExpression()
		left = &ast.BinaryExpression{
			Op:    op.Type,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *Parser) parseAdditiveExpression() ast.Expression {
	left := p.parseMultiplicativeExpression()
	for {
		op := p.token
		if !p.test(tokenize.TokenPlus, tokenize.TokenMinus) {
			break
		}
		p.next()
		right := p.parseMultiplicativeExpression()
		left = &ast.BinaryExpression{
			Op:    op.Type,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *Parser) parseMultiplicativeExpression() ast.Expression {
	left := p.parseUnaryExpression()
	for {
		op := p.token
		if !p.test(tokenize.TokenMul, tokenize.TokenDiv, tokenize.TokenMod) {
			break
		}
		p.next()
		right := p.parseUnaryExpression()
		left = &ast.BinaryExpression{
			Op:    op.Type,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	op := p.token
	if p.test(tokenize.TokenPlus, tokenize.TokenMinus, tokenize.TokenNot, tokenize.TokenBitNot) {
		p.next()
		return &ast.UnaryExpression{
			Op:         op.Type,
			Expression: p.parseUnaryExpression(),
		}
	}
	return p.parsePrimaryExpression()
}

func (p *Parser) parsePrimaryExpression() ast.Expression {
	left := p.parseAtomExpression()
loop:
	for {
		switch p.token.Type {
		case tokenize.TokenOpenParen:
			p.next()
			var args *ast.ExpressionList
			if p.test(tokenize.TokenCloseParen) {
				args = ast.EmptyExpressionList
			} else {
				args = p.parseExpressionList(false)
			}
			p.expect(tokenize.TokenCloseParen)
			left = &ast.CallFunctionExpression{
				Callable: left,
				Args:     args,
			}
		case tokenize.TokenOpenBracket:
			p.next()
			left = &ast.IndexAccessExpression{
				Value:  left,
				Index:  p.parseExpression(),
				Assign: false,
			}
			p.expect(tokenize.TokenCloseBracket)
		case tokenize.TokenDot:
			p.next()
			left = &ast.AttributeAccessExpression{
				Value:  left,
				Name:   p.expect(tokenize.TokenIdentifier).Value.(string),
				Assign: false,
			}
		default:
			break loop
		}
	}
	return left
}

func (p *Parser) parseAtomExpression() ast.Expression {
	if p.empty() {
		panic(fmt.Errorf("<expression> expected near '<eof>"))
	}

	switch token := p.token; token.Type {
	case tokenize.TokenNull:
		p.next()
		return &ast.NullLiteralExpression{}
	case tokenize.TokenTrue:
		p.next()
		return &ast.TrueLiteralExpression{}
	case tokenize.TokenFalse:
		p.next()
		return &ast.FalseLiteralExpression{}
	case tokenize.TokenLiteralInt:
		p.next()
		return &ast.IntLiteralExpression{Value: token.Value.(int64)}
	case tokenize.TokenLiteralFloat:
		p.next()
		return &ast.FloatLiteralExpression{Value: token.Value.(float64)}
	case tokenize.TokenLiteralString:
		p.next()
		proto := token.Value.(string)
		return &ast.StringLiteralExpression{
			Value: proto[1 : len(proto)-1],
			Proto: proto,
		}
	case tokenize.TokenIdentifier:
		p.next()
		return &ast.IdentifierExpression{Name: token.Value.(string), Assign: false}
	case tokenize.TokenOpenParen:
		p.next()
		expression := p.parseExpression()
		p.expect(tokenize.TokenCloseParen)
		return expression
	case tokenize.TokenOpenBracket:
		p.next()
		result := &ast.ListLiteralExpression{
			Value: ast.EmptyExpressionList,
		}
		if !p.test(tokenize.TokenCloseBracket) {
			result.Value = p.parseExpressionList(true)
		}
		p.expect(tokenize.TokenCloseBracket)
		return result
	case tokenize.TokenOpenBrace:
		p.next()
		result := &ast.DictLiteralExpression{}
		if p.test(tokenize.TokenCloseBrace) {
			result.Value = map[string]ast.Expression{}
		} else {
			result.Value = p.parseDictLiteral()
		}
		p.expect(tokenize.TokenCloseBrace)
		return result
	case tokenize.TokenFunction:
		p.next()
		p.expect(tokenize.TokenOpenParen)
		result := &ast.FunctionDeclareExpression{}
		if !p.test(tokenize.TokenCloseParen) {
			result.ParameterNames = p.parseIdentifierList()
		}
		p.expect(tokenize.TokenCloseParen)
		result.Body = p.parseBlockStatement()
		return result
	default:
		panic(fmt.Errorf("unexpected token near '%s'", token.Type.String()))
	}
}

// expression (',' expression)*
func (p *Parser) parseExpressionList(skipNewline bool) *ast.ExpressionList {
	if skipNewline {
		p.skipNewline()
	}
	result := &ast.ExpressionList{
		List: []ast.Expression{p.parseExpression()},
	}
	for p.test(tokenize.TokenComma) {
		p.next()
		if skipNewline {
			p.skipNewline()
		}
		result.List = append(result.List, p.parseExpression())
	}
	if skipNewline {
		p.skipNewline()
	}
	return result
}

// identifier (',' identifier)*
func (p *Parser) parseIdentifierList() []string {
	result := []string{p.expect(tokenize.TokenIdentifier).Value.(string)}
	for p.test(tokenize.TokenComma) {
		p.next()
		result = append(result, p.expect(tokenize.TokenIdentifier).Value.(string))
	}
	return result
}

// identifier ':' expression (',' identifier ':' expression)*
func (p *Parser) parseDictLiteral() map[string]ast.Expression {
	result := make(map[string]ast.Expression)
	p.skipNewline()
	key := p.expect(tokenize.TokenIdentifier).Value.(string)
	p.expect(tokenize.TokenColon)
	result[key] = p.parseExpression()
	for p.test(tokenize.TokenComma) {
		p.next()
		p.skipNewline()
		key = p.expect(tokenize.TokenIdentifier).Value.(string)
		p.expect(tokenize.TokenColon)
		result[key] = p.parseExpression()
	}
	p.skipNewline()
	return result
}

func (p *Parser) next() *tokenize.Token {
	p.token = p.lexer.Next()
	return p.token
}

func (p *Parser) test(tt ...tokenize.TokenType) bool {
	for _, t := range tt {
		if t == p.token.Type {
			return true
		}
	}
	return false
}

func (p *Parser) errorMessage(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

func (p *Parser) expect(tt tokenize.TokenType) *tokenize.Token {
	result := p.token
	if p.token.Type != tt {
		p.errorMessage("'%s' expected, but got: '%s'", tt.String(), p.token.Type.String())
	}
	p.next()
	return result
}

func (p *Parser) empty() bool {
	return p.token.Type == tokenize.TokenEof
}

func (p *Parser) skipNewline() {
	for p.token.IsNewLine() {
		p.next()
	}
}
