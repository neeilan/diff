package main

import "fmt"
import "diff/tools/lexer"

func main() {
	New("a + b").Parse()
}

type Parser struct {
	l *lexer.Lexer
}

func New(s string) *Parser {
	return &Parser{l: lexer.New(s)}
}

func (p *Parser) Parse() {
	fmt.Println(p.l.NextToken())
}
