package difflib

import "testing"

var emptyBindings = make(map[string]float64)

func TestDifffNumber(t *testing.T) {
	num := newNum(2.5)
	if val := num.eval(emptyBindings); val != 2.5 {
		t.Fatalf("Got %f as the eval() value of Number 2.5", val)
	}

	der := num.diff()
	if val := der.eval(emptyBindings); val != 0 {
		t.Errorf("Got %f when evaluating [Number %v].diff(), instead of 0", val, num)
	}
}

func TestDiffSumOfProducts(t *testing.T) {
	// 1.1 * 3.3 + 2.3x
	sum := newSum(
		newProduct(
			newNum(1.1),
			newNum(3.3),
		),
		newProduct(
			newNum(-2),
			newVar("x"),
		),
	)

	if val := sum.eval(bindX(1)); val != 1.63 {
		t.Fatalf("Got %f when evaluating sum %v at x=2. Expected 1.63", val, sum.prune())
	}

	// Should == 2.3
	der := sum.diff()
	if val := der.eval(bindX(2)); val != -2 {
		t.Errorf("Got %f when evaluating [Sum %v].diff() at x=2. Expected -2. Derivative: [%v]", val, sum, der.prune())
	}

}

func bindX(x float64) map[string]float64 {
	return map[string]float64{"x": x}
}
