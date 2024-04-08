package lexer

import (
	"github.com/EmilLaursen/wiig/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	tok = token.Ch(string(l.ch))
	switch {

	case tok.Type == token.ASSIGN:
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{
				Type:    token.EQ,
				Literal: string(ch) + string(l.ch),
			}
		}
	case tok.Type == token.BANG:
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: string(ch) + string(l.ch),
			}
		}

	case tok.Type != token.ILLEGAL:
		// we gooood
	// case '=':
	// 	tok = token.Token{Type: token.ASSIGN, Literal: "="}
	// case ';':
	// 	tok = token.Token{Type: token.SEMICOLON, Literal: ";"}
	// case '(':
	// 	tok = token.Token{Type: token.LPAREN, Literal: "("}
	// case ')':
	// 	tok = token.Token{Type: token.RPAREN, Literal: ")"}
	// case ',':
	// 	tok = token.Token{Type: token.COMMA, Literal: ","}
	// case '+':
	// 	tok = token.Token{Type: token.PLUS, Literal: "+"}
	// case '{':
	// 	tok = token.Token{Type: token.LBRACE, Literal: "{"}
	// case '}':
	// 	tok = token.Token{Type: token.RBRACE, Literal: "}"}
	// case 0:
	// 	tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		switch {
		case l.ch == '"':
			tok = token.Str(l.readString())
		case l.ch == '[':
			tok.Type = token.LBRACKET
			tok.Literal = string(l.ch)
		case l.ch == ']':
			tok.Type = token.RBRACKET
			tok.Literal = string(l.ch)
		case isLetter(l.ch):
			return token.Ident(l.readIdentifier())
		case isDigit(l.ch):
			return token.Num(l.readNumber())
		}
		// tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
	}
	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readWhile(predicate func(ch byte) bool) string {
	pos := l.position
	for predicate(l.ch) && l.ch != 0 {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() string {
	return l.readWhile(isDigit)
}

func (l *Lexer) readIdentifier() string {
	return l.readWhile(isLetter)
}

func (l *Lexer) readString() string {
	// TODO: handle escaping, " \n \t \r etc
	pos := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[pos:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	// only supports ascii
	// TODO: add unicode support (and emojis)
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
