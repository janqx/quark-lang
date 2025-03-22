package parser

import (
	"fmt"
	"io"
	"strconv"
	"unicode"

	"github.com/janqx/quark-lang/v1/tokenize"
)

const EOF rune = -1

type Lexer struct {
	filename        string
	reader          io.RuneReader
	ch              rune
	offset          int
	line, column    int
	lines           map[int]int // number of columns per line
	currentToken    *tokenize.Token
	lookaheadToken  *tokenize.Token
	currentPosition *tokenize.Position
}

func NewLexer(filename string, reader io.RuneReader) *Lexer {
	l := &Lexer{
		filename: filename,
		reader:   reader,
		ch:       0,
		offset:   0,
		line:     1,
		column:   1,
		lines:    make(map[int]int),
	}
	l.currentToken = nil
	l.lookaheadToken = nil
	l.advance()
	return l
}

func (l *Lexer) advance() rune {
	var err error
	l.ch, _, err = l.reader.ReadRune()
	if err == nil {
		l.offset++
		if l.ch == '\n' {
			l.lines[l.line] = l.column
			l.line++
			l.column = 1
		} else {
			l.column++
		}
	} else {
		l.ch = EOF
	}
	return l.ch
}

func (l *Lexer) skipComment() {
	first := l.ch
	l.advance()
	for l.ch != EOF {
		if first == '*' {
			if l.ch == '*' {
				if l.advance() == '/' {
					l.advance()
					break
				}
			} else {
				l.advance()
			}
		} else {
			if l.ch == '\n' {
				l.advance()
				break
			} else {
				l.advance()
			}
		}
	}
}

func (l *Lexer) lexNumber() *tokenize.Token {
	first := l.ch
	hex := false
	token := l.makeToken(tokenize.TokenLiteralInt)
	s := []rune{l.ch}
	l.advance()
	if first == '0' && (l.ch == 'x' || l.ch == 'X') {
		hex = true
		s = append(s, l.ch)
		l.advance()
	}
	for l.ch != EOF {
		if hex {
			if (l.ch >= '0' && l.ch <= '9') || (l.ch >= 'a' && l.ch <= 'f') || (l.ch >= 'A' && l.ch <= 'F') {
				s = append(s, l.ch)
				l.advance()
			} else {
				break
			}
		} else {
			if l.ch >= '0' && l.ch <= '9' {
				s = append(s, l.ch)
				l.advance()
			} else {
				break
			}
		}
	}
	if !hex && l.ch == '.' {
		s = append(s, '.')
		l.advance()
		for l.ch != EOF {
			if l.ch >= '0' && l.ch <= '9' {
				s = append(s, l.ch)
				l.advance()
			} else {
				break
			}
		}
		value, err := strconv.ParseFloat(string(s), 64)
		if err != nil {
			panic(err)
		}
		token.Type = tokenize.TokenLiteralFloat
		token.Value = float64(value)
	} else {
		var value int64
		var err error
		if hex {
			value, err = strconv.ParseInt(string(s[2:]), 16, 64)
		} else {
			value, err = strconv.ParseInt(string(s), 10, 64)
		}
		if err != nil {
			panic(err)
		}
		token.Value = int64(value)
	}
	return token
}

func (l *Lexer) readEscape() (rune, error) {
	l.advance() // skip '\'
	switch l.ch {
	case 'n':
		return '\n', nil
	case 'r':
		return '\r', nil
	case 't':
		return '\t', nil
	case 'v':
		return '\v', nil
	case 'b':
		return '\b', nil
	case 'f':
		return '\f', nil
	case 'a':
		return '\a', nil
	case '\\':
		return '\\', nil
	case '\'':
		return '\'', nil
	case '"':
		return '"', nil
	case '0':
		return 0, nil
	default:
		return EOF, fmt.Errorf("illegal escape sequence")
	}
}

func (l *Lexer) lexSimpleString() *tokenize.Token {
	first := l.ch
	s := []rune{l.ch}
	token := l.makeToken(tokenize.TokenLiteralString)
	l.advance()
	for l.ch != EOF {
		if l.ch == '\\' {
			ch, err := l.readEscape()
			if err != nil {
				panic(err)
			}
			l.ch = ch
			l.advance()
			continue
		}
		s = append(s, l.ch)
		if l.ch == first {
			l.advance()
			token.Value = string(s)
			return token
		} else if l.ch == '\n' {
			break
		}
		l.advance()
	}
	panic(fmt.Errorf("%s:%d string literal not terminated", l.filename, token.Position.Line))
}

func (l *Lexer) lexLongString() *tokenize.Token {
	s := []rune{l.ch}
	token := l.makeToken(tokenize.TokenLiteralString)
	l.advance()
	for l.ch != EOF {
		s = append(s, l.ch)
		if l.ch == '`' {
			l.advance()
			token.Value = string(s)
			return token
		}
		l.advance()
	}
	panic(fmt.Errorf("%s:%d string literal not terminated", l.filename, token.Position.Line))
}

func (l *Lexer) makePosition() *tokenize.Position {
	return &tokenize.Position{
		Filename: l.filename,
		Offset:   l.offset,
		Line:     l.line,
		Column:   l.column,
	}
}

func (l *Lexer) makeToken(tokenType tokenize.TokenType) *tokenize.Token {
	return &tokenize.Token{
		Type:     tokenType,
		Value:    nil,
		Position: l.currentPosition,
	}
}

func (l *Lexer) scan() *tokenize.Token {
	for {
		l.currentPosition = l.makePosition()
		switch l.ch {
		case ' ', '\t', '\v', '\f', '\r':
			l.advance()
		case '\n':
			l.advance()
			return l.makeToken(tokenize.TokenNewline)
		case '/':
			l.advance()
			if l.ch == '/' || l.ch == '*' {
				l.skipComment()
			} else {
				return l.makeToken(tokenize.TokenDiv)
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return l.lexNumber()
		case 0x27, '"':
			return l.lexSimpleString()
		case '`':
			return l.lexLongString()
		case '[', ']', '(', ')', '{', '}', ',', ';', ':', '?', '.':
			ch := l.ch
			l.advance()
			return l.makeToken(tokenize.SeparatorToTokenType[ch])
		case '+', '-', '*', '%', '~', '^':
			ch := l.ch
			l.advance()
			return l.makeToken(tokenize.SingleOperatorToTokenType[ch])
		case '=':
			if l.advance() == '=' {
				l.advance()
				return l.makeToken(tokenize.TokenEQ)
			} else {
				return l.makeToken(tokenize.TokenAssign)
			}
		case '!':
			if l.advance() == '=' {
				l.advance()
				return l.makeToken(tokenize.TokenNEQ)
			} else {
				return l.makeToken(tokenize.TokenNot)
			}
		case '<':
			l.advance()
			if l.ch == '=' {
				l.advance()
				return l.makeToken(tokenize.TokenLTE)
			}
			if l.ch == '<' {
				l.advance()
				return l.makeToken(tokenize.TokenBitLhs)
			}
			return l.makeToken(tokenize.TokenLT)
		case '>':
			l.advance()
			if l.ch == '=' {
				l.advance()
				return l.makeToken(tokenize.TokenGTE)
			}
			if l.ch == '>' {
				l.advance()
				return l.makeToken(tokenize.TokenBitRhs)
			}
			return l.makeToken(tokenize.TokenGT)
		case '|':
			if l.advance() == '|' {
				l.advance()
				return l.makeToken(tokenize.TokenLogicOr)
			} else {
				return l.makeToken(tokenize.TokenBitOr)
			}
		case '&':
			if l.advance() == '&' {
				l.advance()
				return l.makeToken(tokenize.TokenLogicAnd)
			} else {
				return l.makeToken(tokenize.TokenBitAnd)
			}
		case EOF:
			return l.makeToken(tokenize.TokenEof)
		default:
			if l.ch == '_' || unicode.IsLetter(l.ch) {
				token := l.makeToken(tokenize.TokenIdentifier)
				s := []rune{l.ch}
				l.advance()
				for l.ch == '_' || unicode.IsLetter(l.ch) || unicode.IsDigit(l.ch) {
					s = append(s, l.ch)
					l.advance()
				}
				token.Value = string(s)
				if ktype, ok := tokenize.KeywordToTokenType[token.Value.(string)]; ok {
					token.Type = ktype
				} else {
					token.Type = tokenize.TokenIdentifier
				}
				return token
			} else {
				panic(fmt.Errorf("unrecognized character: '%c'", l.ch))
			}
		}
	}
}

func (l *Lexer) Next() *tokenize.Token {
	if l.lookaheadToken != nil {
		l.currentToken = l.lookaheadToken
		l.lookaheadToken = nil
	} else {
		l.currentToken = l.scan()
	}
	return l.currentToken
}

func (l *Lexer) Lookahead() *tokenize.Token {
	if l.lookaheadToken == nil {
		l.lookaheadToken = l.scan()
	}
	return l.lookaheadToken
}
