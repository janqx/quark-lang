package tokenize

import "fmt"

type TokenType uint8

const (
	TokenEof TokenType = iota

	TokenNewline // \n

	// separators
	TokenOpenBracket  // [
	TokenCloseBracket // ]
	TokenOpenParen    // (
	TokenCloseParen   // )
	TokenOpenBrace    // {
	TokenCloseBrace   // }
	TokenComma        // ,
	TokenSemiColon    // ;
	TokenColon        // :
	TokenQuestion     // ?
	TokenDot          // .

	// operators
	TokenAssign   // =
	TokenPlus     // +
	TokenMinus    // -
	TokenBitNot   // ~
	TokenBitAnd   // &
	TokenBitOr    // |
	TokenBitXor   // ^
	TokenBitLhs   // <<
	TokenBitRhs   // >>
	TokenNot      // !
	TokenMul      // *
	TokenDiv      // /
	TokenMod      // %
	TokenLogicAnd // &&
	TokenLogicOr  // ||
	TokenLT       // <
	TokenLTE      // <=
	TokenGT       // >
	TokenGTE      // >=
	TokenEQ       // ==
	TokenNEQ      // !=

	// keywords
	TokenNull
	TokenTrue
	TokenFalse
	TokenIf
	TokenElse
	TokenFor
	TokenBreak
	TokenContinue
	TokenFunction
	TokenReturn
	TokenClass
	TokenThis
	TokenSuper
	TokenImport
	TokenExport
	TokenDebugger

	// identitie
	TokenIdentifier

	// literal
	TokenLiteralInt
	TokenLiteralFloat
	TokenLiteralString
)

var TokenTypeToString = map[TokenType]string{
	TokenEof: "<eof>",

	TokenNewline: "<newline>", // \n

	// operators
	TokenOpenBracket:  "[",
	TokenCloseBracket: "]",
	TokenOpenParen:    "(",
	TokenCloseParen:   ")",
	TokenOpenBrace:    "{",
	TokenCloseBrace:   "}",
	TokenComma:        ",",
	TokenSemiColon:    ";",
	TokenColon:        ":",
	TokenQuestion:     "?",
	TokenDot:          ".",
	TokenAssign:       "=",
	TokenPlus:         "+",
	TokenMinus:        "-",
	TokenBitNot:       "~",
	TokenBitAnd:       "&",
	TokenBitOr:        "|",
	TokenBitXor:       "^",
	TokenBitLhs:       "<<",
	TokenBitRhs:       ">>",
	TokenNot:          "!",
	TokenMul:          "*",
	TokenDiv:          "/",
	TokenMod:          "%",
	TokenLogicAnd:     "&&",
	TokenLogicOr:      "||",
	TokenLT:           "<",
	TokenLTE:          "<=",
	TokenGT:           ">",
	TokenGTE:          ">=",
	TokenEQ:           "==",
	TokenNEQ:          "!=",

	// keywords
	TokenNull:     "null",
	TokenTrue:     "true",
	TokenFalse:    "false",
	TokenIf:       "if",
	TokenElse:     "else",
	TokenFor:      "for",
	TokenBreak:    "break",
	TokenContinue: "continue",
	TokenFunction: "fn",
	TokenReturn:   "return",
	TokenClass:    "class",
	TokenThis:     "this",
	TokenSuper:    "super",
	TokenImport:   "__import__",
	TokenExport:   "export",
	TokenDebugger: "debugger",

	// identitie
	TokenIdentifier: "<identifier>",

	// literal
	TokenLiteralInt:    "<literal-int>",
	TokenLiteralFloat:  "<literal-float>",
	TokenLiteralString: "<literal-string>",
}

var SeparatorToTokenType = map[rune]TokenType{
	'[': TokenOpenBracket,  // [
	']': TokenCloseBracket, // ]
	'(': TokenOpenParen,    // (
	')': TokenCloseParen,   // )
	'{': TokenOpenBrace,    // {
	'}': TokenCloseBrace,   // }
	',': TokenComma,        // ,
	';': TokenSemiColon,    // ;
	':': TokenColon,        // :
	'?': TokenQuestion,     // ?
	'.': TokenDot,          // .
}

var SingleOperatorToTokenType = map[rune]TokenType{
	'+': TokenPlus,
	'-': TokenMinus,
	'*': TokenMul,
	'/': TokenDiv,
	'%': TokenMod,
	'~': TokenBitNot,
	'^': TokenBitXor,
}

var OperatorToTokenType = map[string]TokenType{
	"=":  TokenAssign,   // =
	"+":  TokenPlus,     // +
	"-":  TokenMinus,    // -
	"~":  TokenBitNot,   // ~
	"&":  TokenBitAnd,   // &
	"|":  TokenBitOr,    // |
	"^":  TokenBitXor,   // ^
	"<<": TokenBitLhs,   // <<
	">>": TokenBitRhs,   // >>
	"!":  TokenNot,      // !
	"*":  TokenMul,      // *
	"/":  TokenDiv,      // /
	"%":  TokenMod,      // %
	"&&": TokenLogicAnd, // &&
	"||": TokenLogicOr,  // ||
	"<":  TokenLT,       // <
	"<=": TokenLTE,      // <=
	">":  TokenGT,       // >
	">=": TokenGTE,      // >=
	"==": TokenEQ,       // ==
	"!=": TokenNEQ,      // !=
}

var KeywordToTokenType = map[string]TokenType{
	"null":       TokenNull,
	"true":       TokenTrue,
	"false":      TokenFalse,
	"if":         TokenIf,
	"else":       TokenElse,
	"for":        TokenFor,
	"break":      TokenBreak,
	"continue":   TokenContinue,
	"fn":         TokenFunction,
	"return":     TokenReturn,
	"class":      TokenClass,
	"this":       TokenThis,
	"super":      TokenSuper,
	"__import__": TokenImport,
	"export":     TokenExport,
	"debugger":   TokenDebugger,
}

func (tt TokenType) String() string {
	return TokenTypeToString[tt]
}

type Token struct {
	Type     TokenType
	Value    interface{} // int64 float64 string
	Position *Position
}

func (t *Token) IsNewLine() bool {
	return t.Type == TokenNewline
}

func (t *Token) Clone() *Token {
	return &Token{
		Type:     t.Type,
		Value:    t.Value,
		Position: t.Position,
	}
}

func (t *Token) String() (s string) {
	s = "("
	if t.Type == TokenIdentifier {
		s += fmt.Sprintf("<identifier %s>", t.Value.(string))
	} else if t.Type == TokenLiteralInt {
		s += fmt.Sprintf("<literal-int %d>", t.Value.(int64))
	} else if t.Type == TokenLiteralFloat {
		s += fmt.Sprintf("<literal-float %f>", t.Value.(float64))
	} else if t.Type == TokenLiteralString {
		s += fmt.Sprintf("<literal-string %s>", t.Value.(string))
	} else {
		s += TokenTypeToString[t.Type]
	}
	s += fmt.Sprintf(" %d, %d)", t.Position.Line, t.Position.Column)
	return
}
