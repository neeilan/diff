package main

import "fmt"
import "diff/tools/token"

type Lexer struct {
	input                  string
	position, readPosition int
	ch                     byte
}

func main() {
	input := "(1+ 3 + 4.5"
	l := NewLexer(input)
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		fmt.Println(tok)
	}
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.ch {
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '^':
		tok = newToken(token.POWER, l.ch)
	case 0:
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.Lookup(tok.Literal)
		} else if isDigit(l.ch) {
			tok.Type = token.NUM
			tok.Literal = l.readNumber()
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar() // advance
	return tok
}

// PRIVATE

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // NULL
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) readIdentifier() string {
	initialPosition := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[initialPosition:l.position]
}

func (l *Lexer) readNumber() string {
	initialPosition := l.position
	numDecimals := 0
	for isDigit(l.ch) || l.ch == byte('.') {
		if l.ch == byte('.') {
			numDecimals += 1
		}
		if numDecimals > 1 {
			panic("Found more than one decimal point in a number")
		}
		l.readChar()
	}

	return l.input[initialPosition:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
