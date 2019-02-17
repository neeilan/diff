package main

import "fmt"

type Expression interface {
	diff() Expression
	isZero() bool
	prune() Expression
	String() string
}

func isOne(expr Expression) bool {
	num, ok := expr.(*Number)
	if !ok {
		return false
	}
	return num.value == 1
}
/*-------------------------------*/
// Sum
/*-------------------------------*/
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

/*-------------------------------*/
// Product
/*-------------------------------*/

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

	if isOne(pr.f) {
		return pr.g
	} else if isOne(pr.g) {
		return pr.f
	}

	return pr
}

func (pr *Product) String() string {
	return "(" + pr.f.String() + " * " + pr.g.String() + ")"
}

/*-------------------------------*/
// Variable
/*-------------------------------*/

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

/*-------------------------------*/
// Number
/*-------------------------------*/

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

/*-------------------------------*/
// ToNumericPower
/*-------------------------------*/

// Need a more abstract function but here we go
type ToNumericPower struct {
	operand  Expression
	exponent Expression
}

func toNumericPower(expr Expression, exponent Expression) *ToNumericPower {
	if _, isNum := exponent.(*Number); !isNum {
		panic("Can only raise to powers of numbers currently.")
	}
	return &ToNumericPower{operand: expr, exponent: exponent}
}

func (pow *ToNumericPower) diff() Expression {
	exponent, _ := pow.exponent.(*Number)

	fPrimeOfG := newProduct(exponent, toNumericPower(pow.operand, newNum(exponent.value-1)))
	gPrime := pow.operand.diff()
	return newProduct(fPrimeOfG, gPrime)
}

func (pow *ToNumericPower) isZero() bool {
	// TODO: should we add an numerical exponent check?
	return pow.operand.isZero()
}

func (pow *ToNumericPower) prune() Expression {
	pow.operand = pow.operand.prune()
	pow.exponent = pow.exponent.prune()

	if pow.exponent.isZero() {
		return newNum(1)
	}

	if pow.isZero() {
		return newNum(0)
	} else if isOne(pow.exponent) {
		return pow.operand
	}

	return pow
}

func (pow *ToNumericPower) String() string {
	exponent, _ := pow.exponent.(*Number)
	return fmt.Sprintf("("+pow.operand.String()+"^%f)", exponent.value)
}

/*-------------------------------*/
// Log
/*-------------------------------*/

type Log struct {
	operand Expression
}

func newLog (f Expression) *Log {
	return &Log{f}
}

func (l *Log) isZero() bool {
	return isOne(l.operand)
}

func (l *Log) diff() Expression {
	f := l.operand
	switch v:= f.(type) {
		case *ToNumericPower:
			exp := v.exponent
			base := v.operand
			simplified := newProduct(exp, newLog(base))
			return simplified.diff()
		case *Product:
			term1 := newLog(v.f)
			term2 := newLog(v.g)
			return newSum(term1, term2).diff()
	}
	return newProduct(f.diff(), toNumericPower(f, newNum(-1)))
}

func (l *Log) prune() Expression {
	f := l.operand
	switch v := f.(type) {
		case *ToNumericPower:
			exp := v.exponent
			base := v.operand
			return newProduct(exp.prune(), newLog(base.prune()))
		case *Product:
			return newSum(v.f.prune(), v.g.prune())
	}
	return newLog(f.prune())

}

func (l *Log) String() string {
	return "log (" + l.operand.String() + ")"
}

func main() {
	expr := newSum(newSum(newProduct(newNum(1.1), newNum(3.3)), newProduct(newNum(2.2), newVar("x"))), toNumericPower(newVar("x"), newNum(3)))
	fmt.Println("Expression: ", expr)

	derivative := expr.diff()
	fmt.Println("Derivative: ", derivative)
	fmt.Println("Derivative (pruned): ", derivative.prune())

	derivative2 := derivative.diff()
	fmt.Println("2nd derivative: ", derivative2)
	fmt.Println("2nd derivative (pruned): ", derivative2.prune())

	recip := toNumericPower(newVar("x"), newNum(-1))
	derivativeRecip := recip.diff()
	fmt.Println("Derivative of reciprocal:", derivativeRecip)
	fmt.Println("Derivative of reciprocal (pruned):", derivativeRecip.prune())

	logTest := newLog(toNumericPower(newVar("x"), newNum(3)))
	derivativeLogCubic := logTest.diff()
	fmt.Println("Log of cubic:", logTest.prune())
	fmt.Println("Derivative of log of cubic: ", derivativeLogCubic)
	fmt.Println("Derivative of log of cubic, pruned: ", derivativeLogCubic.prune())
	fmt.Println("Derivative of pruned log of cubic: ", logTest.prune().diff())
	fmt.Println("Derivative of pruned log of cubic, pruned: ", logTest.prune().diff().prune())

}