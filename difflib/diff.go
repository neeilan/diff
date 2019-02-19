package difflib

import "fmt"
import "math"

type Expression interface {
	diff() Expression
	eval(map[string]float64) float64
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

func (s *Sum) diff() Expression {
	fprime := s.f.diff()
	gprime := s.g.diff()

	return newSum(fprime, gprime)
}

func (s *Sum) eval(bindings map[string]float64) float64 {
	return s.f.eval(bindings) + s.g.eval(bindings)
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

func (pr *Product) eval(bindings map[string]float64) float64 {
	return pr.f.eval(bindings) * pr.g.eval(bindings)
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

func newVar(name string) *Variable {
	return &Variable{name: name}
}

func (*Variable) diff() Expression {
	return &Number{value: float64(1)}
}

func (v *Variable) eval(bindings map[string]float64) float64 {
	val, ok := bindings[v.name]
	if !ok {
		panic(fmt.Sprintf("Attempted to evaluate variable %s, not no binding provided! Available bindings: %v", v.name, bindings))
	}
	return val
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

func (n *Number) eval(bindings map[string]float64) float64 {
	return n.value
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
// Log
/*-------------------------------*/

type Log struct {
	operand Expression
}

func newLog(f Expression) *Log {
	return &Log{f}
}

func (l *Log) isZero() bool {
	return isOne(l.operand)
}

func (l *Log) diff() Expression {
	f := l.operand
	switch v := f.(type) {
	case *ToGenericPower:
		exp := v.exponent
		base := v.operand
		simplified := newProduct(exp, newLog(base))
		return simplified.diff()
	case *Product:
		term1 := newLog(v.f)
		term2 := newLog(v.g)
		return newSum(term1, term2).diff()
	}
	return newProduct(f.diff(), toGenericPower(f, newNum(-1)))
}

func (l *Log) eval(bindings map[string]float64) float64 {
	return math.Log(l.operand.eval(bindings))
}

func (l *Log) prune() Expression {
	f := l.operand
	switch v := f.(type) {
	case *ToGenericPower:
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

/*-------------------------------*/
// ToGenericPower
/*-------------------------------*/

// Need a more abstract function but here we go
type ToGenericPower struct {
	operand  Expression
	exponent Expression
}

func toGenericPower(expr Expression, exponent Expression) *ToGenericPower {
	//if _, isNum := exponent.(*Number); !isNum {
	//	panic("Can only raise to powers of numbers currently.")
	//}
	return &ToGenericPower{operand: expr, exponent: exponent}
}

func (pow *ToGenericPower) diff() Expression {
	pow.operand = pow.operand.prune()
	pow.exponent = pow.exponent.prune()
	// functional forms:
	// If c is a constant, and f(x), g(x) are arbitrary functions, we have
	// 1. c^f(x)
	// 2. f(x) ^ c
	// 3. f(x) ^ g(x) : requires logarithmic differentiation
	if base, baseIsNumber := pow.operand.(*Number); baseIsNumber {
		//case 1
		logC := newNum(math.Log(base.value))
		return newProduct(logC, newProduct(pow, pow.exponent.diff()))
	} else if exp, expIsNumber := pow.exponent.(*Number); expIsNumber { // case 2
		return newProduct(pow.exponent, newProduct(toGenericPower(pow.operand, newNum(exp.value-1)), pow.operand.diff()))
	} else { // case 3
		logDiffSubproblem := newProduct(newLog(pow.operand), pow.exponent).diff()
		return newProduct(pow, logDiffSubproblem)
	}
}

func (pow *ToGenericPower) eval(bindings map[string]float64) float64 {
	return math.Pow(pow.operand.eval(bindings), pow.exponent.eval(bindings))
}

func (pow *ToGenericPower) isZero() bool {
	return pow.operand.isZero()
}

func (pow *ToGenericPower) prune() Expression {
	pow.operand = pow.operand.prune()
	pow.exponent = pow.exponent.prune()

	if pow.exponent.isZero() && pow.operand.isZero() {
		panic("0 ^ 0 is not well-defined")
	} else if pow.exponent.isZero() {
		return newNum(1)
	} else if pow.operand.isZero() {
		return newNum(0)
	} else if isOne(pow.exponent) {
		return pow.operand
	}
	return toGenericPower(pow.operand.prune(), pow.exponent.prune())
}

func (pow *ToGenericPower) String() string {
	return pow.operand.String() + "^(" + pow.exponent.String() + ")"
}

/*-------------------------------*/
// Sine and Cosine
/*-------------------------------*/

type Sine struct {
	operand Expression
}

func newSine(f Expression) *Sine {
	return &Sine{f}
}

func (s *Sine) isZero() bool {
	// currently ignoring the repeated zeros of the sinusoid
	// TODO: what should we do about the set of zeros?
	return s.operand.isZero()
}

func (s *Sine) diff() Expression {
	f := s.operand
	return newProduct(newCosine(f), f.diff())
}

func (s *Sine) eval(bindings map[string]float64) float64 {
	return math.Sin(s.operand.eval(bindings))
}

type Cosine struct {
	operand Expression
}

func newCosine(f Expression) *Cosine {
	return &Cosine{f}
}

func (c *Cosine) isZero() bool {
	// TODO: what should we do about the zeros of cosine?
	return newSum(c.operand, newNum(math.Pi/2)).isZero()
}

func (s *Cosine) eval(bindings map[string]float64) float64 {
	return math.Cos(s.operand.eval(bindings))
}

func (c *Cosine) diff() Expression {
	f := c.operand
	return newProduct(newProduct(newSine(f), newNum(-1)), f.diff())
}

func (s *Sine) prune() Expression {
	f := s.operand
	f = f.prune()
	// TODO: add sin(a+b) = sin(a)cos(b) + sin(b)cos(a)
	return newSine(f)
}

func (c *Cosine) prune() Expression {
	f := c.operand
	f = f.prune()
	// TODO: add cos(a+b) = cos(a)cos(b) - sin(b)sin(a)
	return newCosine(f)
}

func (s *Sine) String() string {
	return "sin (" + s.operand.String() + ")"
}

func (c *Cosine) String() string {
	return "cos (" + c.operand.String() + ")"
}

func main() {
	// TODO: Stuff here should be moved into tests, and this main() function should be removed.

	expr := newSum(newSum(newProduct(newNum(1.1), newNum(3.3)), newProduct(newNum(2.2), newVar("x"))), toGenericPower(newVar("x"), newNum(3)))
	fmt.Println("Expression: ", expr)

	derivative := expr.diff()
	// fmt.Println("Derivative: ", derivative)
	fmt.Println("Derivative (pruned): ", derivative.prune())
	fmt.Println("Derivative, evaluated at x=3", derivative.eval(map[string]float64{"x": 3}))

	derivative2 := derivative.diff()
	// fmt.Println("2nd derivative: ", derivative2)
	fmt.Println("2nd derivative (pruned): ", derivative2.prune())

	recip := toGenericPower(newVar("x"), newNum(-1))
	derivativeRecip := recip.diff()
	// fmt.Println("Derivative of reciprocal:", derivativeRecip)
	fmt.Println("Derivative of reciprocal (pruned):", derivativeRecip.prune())

	logTest := newLog(toGenericPower(newVar("x"), newNum(3)))
	derivativeLogCubic := logTest.diff()
	fmt.Println("Log of cubic:", logTest.prune())
	// fmt.Println("Derivative of log of cubic: ", derivativeLogCubic)
	fmt.Println("Derivative of log of cubic, pruned: ", derivativeLogCubic.prune())
	// fmt.Println("Derivative of pruned log of cubic: ", logTest.prune().diff())
	fmt.Println("Derivative of pruned log of cubic, pruned: ", logTest.prune().diff().prune())
	fmt.Println("Derivative, evaluated at x=2", logTest.diff().prune().eval(map[string]float64{"x": 2}))

	genericPowerTest := toGenericPower(newSum(newVar("x"), newLog(newVar("x"))), newNum(2))
	derivativeGenericPowerTest := genericPowerTest.diff()
	fmt.Println("Generic Power function: ", genericPowerTest)
	// fmt.Println("Derivative of generic power function: ", derivativeGenericPowerTest)
	fmt.Println("Derivative of generic power function, pruned: ", derivativeGenericPowerTest.prune())

	genericPowerTest2 := toGenericPower(newNum(4), newVar("x"))
	derivativeGenericPowerTest2 := genericPowerTest2.diff()
	fmt.Println("Generic Power function 2: ", genericPowerTest2)
	// fmt.Println("Derivative of generic power function 2: ", derivativeGenericPowerTest2)
	fmt.Println("Derivative of generic power function 2, pruned: ", derivativeGenericPowerTest2.prune())

	genericPowerTest3 := toGenericPower(newVar("x"), newVar("x"))
	derivativeGenericPowerTest3 := genericPowerTest3.diff()
	fmt.Println("Generic Power function 3: ", genericPowerTest3)
	// fmt.Println("Derivative of generic power function 3: ", derivativeGenericPowerTest3)
	fmt.Println("Derivative of generic power function 3, pruned: ", derivativeGenericPowerTest3.prune())

	sineTest1 := newSum(newSine(toGenericPower(newVar("x"), newNum(5))),
		newCosine(newProduct(newVar("x"), newLog(newVar("x")))))
	derivativeSineTest1 := sineTest1.diff()
	fmt.Println("Sinusoidal function: ", sineTest1)
	// fmt.Println("Sinusoidal derivative: ", derivativeSineTest1)
	fmt.Println("Sinusoidal derivative, pruned: ", derivativeSineTest1.prune())

}
