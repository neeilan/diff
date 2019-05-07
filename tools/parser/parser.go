package main

import "fmt"
import "diff/tools/lexer"
import "diff/tools/token"
import df "diff/difflib"

func main() {
	New("cos(a+b)").Parse()
}

type Parser struct {
	l          *lexer.Lexer
	curr, peek token.Token
}

func New(s string) *Parser {
	p := &Parser{l: lexer.New(s)}
	p.advance() // initialize curr, peek
	p.advance()
	return p
}

func (p *Parser) Parse() df.Expression {
	e := p.parseExpression(0)
	fmt.Println(e)
	return e
}

func (p *Parser) advance() {
	p.curr = p.peek
	p.peek = p.l.NextToken()
}

func (p *Parser) parseExpression(precedence int) df.Expression {
	return &df.Number{}
}
