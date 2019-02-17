package main

import "fmt"

type Expression interface {
	diff() Expression
	isZero() bool
	prune() Expression
	String() string
}

type Sum struct {
	// f + g
	f Expression
	g Expression
}

func newSum(f Expression, g Expression) *Sum {
	return &Sum{f: f, g: g}
}

func (pr *Sum) diff() Expression {
	fprime := pr.f.diff()
	gprime := pr.g.diff()

	return newSum(fprime, gprime)
}

func (s *Sum) isZero() bool {
	return s.f.isZero() && s.g.isZero()
}

func (s *Sum) prune() Expression {
	s.f = s.f.prune()
	s.g = s.g.prune()

	if s.f.isZero() && s.g.isZero() {
		return newNum(0)
	} else if s.f.isZero() {
		return s.g
	} else if s.g.isZero() {
		return s.f
	}
	return s
}

func (s *Sum) String() string {
	return "(" + s.f.String() + " + " + s.g.String() + ")"
}

type Product struct {
	// f * g
	f Expression
	g Expression
}

func newProduct(f Expression, g Expression) *Product {
	return &Product{f: f, g: g}
}

func (pr *Product) diff() Expression {
	f, g := pr.f, pr.g

	fPrime := f.diff()
	gPrime := g.diff()

	return newSum(newProduct(fPrime, g), newProduct(f, gPrime))
}

func (pr *Product) isZero() bool {
	return pr.f.isZero() && pr.g.isZero()
}

func (pr *Product) prune() Expression {
	pr.f = pr.f.prune()
	pr.g = pr.g.prune()

	if pr.f.isZero() || pr.g.isZero() {
		return newNum(0)
	}
	return pr
}

func (pr *Product) String() string {
	return "(" + pr.f.String() + " * " + pr.g.String() + ")"
}

type Variable struct {
	name string
}

func (*Variable) diff() Expression {
	return &Number{value: float64(1)}
}

func newVar(name string) *Variable {
	return &Variable{name: name}
}

func (*Variable) isZero() bool {
	return false
}

func (v *Variable) prune() Expression {
	return v
}

func (v *Variable) String() string {
	return v.name
}

type Number struct {
	value float64
}

func newNum(value float64) *Number {
	return &Number{value: value}
}

func (*Number) diff() Expression {
	return &Number{value: float64(0)}
}

func (n *Number) isZero() bool {
	return n.value == 0
}

func (n *Number) prune() Expression {
	return n
}

func (num *Number) String() string {
	return fmt.Sprintf("%f", num.value)
}

// Need a more abstract function but here we go
type Square struct {
	operand Expression
}

func square(expr Expression) *Square {
	return &Square{operand: expr}
}

func (pow *Square) diff() Expression {
	fPrimeOfG := newProduct(newNum(2), pow.operand)
	gPrime := pow.operand.diff()
	return newProduct(fPrimeOfG, gPrime)
}

func (pow *Square) isZero() bool {
	return pow.operand.isZero()
}

func (pow *Square) prune() Expression {
	if pow.isZero() {
		return newNum(0)
	}
	return pow
}

func (pow *Square) String() string {
	return "(" + pow.operand.String() + "^2)"
}

func main() {
	expr := newSum(newSum(newProduct(newNum(1.1), newNum(3.3)), newProduct(newNum(2.2), newVar("x"))), square(newVar("x")))
	fmt.Println("Expression: ", expr)

	derivative := expr.diff()
	fmt.Println("Derivative: ", derivative)
	fmt.Println("Derivative (pruned): ", derivative.prune())

	derivative2 := derivative.diff()
	fmt.Println("2nd derivative: ", derivative2)
	fmt.Println("2nd derivative (pruned): ", derivative2.prune())

}
