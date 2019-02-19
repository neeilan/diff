package difflib

import (
	"math"
	"testing"
)

var emptyBindings = make(map[string]float64)

func TestDiffNumber(t *testing.T) {
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

func TestDiffLog(t *testing.T) {
	simpleLog := newLog(newVar("x"))
	logDerivative := simpleLog.diff()
	x0 := 4.5
	if val := logDerivative.eval(bindX(x0)); val != 1/x0 {
		t.Errorf("Got %f when evaluation [%v].diff() at x=%f. Expected %f. Derivative: [%v]", val, simpleLog, x0, 1/x0, logDerivative.prune())
	}

	logCubed := newLog(toGenericPower(newVar("x"), newNum(3)))
	logCubedDerivative := logCubed.diff()
	if val := logCubedDerivative.eval(bindX(x0)); val != 3/x0 {
		t.Errorf("Got %f when evaluation [%v].diff() at x=%f. Expected %f. Derivative: [%v]", val, logCubed, x0, 3/x0, logCubedDerivative.prune())
	}
}
// TODO: Write or find a float equality function (ie within tolerance)
func TestDiffGenericPower(t *testing.T) {
	recip := toGenericPower(newVar("x"), newNum(-1))
	derivativeRecip := recip.diff()

	if val := derivativeRecip.eval(bindX(0.5)); val != -4 {
		t.Errorf("Got %f when evaluating [%v].diff() at x=0.5. Expected -4. Derivative: [%v]", val, recip, derivativeRecip.prune())
	}

	square := toGenericPower(newVar("x"), newNum(2))
	squareDiff := square.diff()
	x0 := 6.8
	if val := squareDiff.eval(bindX(x0)); val != 2*x0 {
		t.Errorf("Got %f when evaluating [%v].diff() at x=%f. Expected %f. Derivative: [%v]", val, square, x0,
			2*x0, squareDiff.prune())
	}

	exponential := toGenericPower(newNum(4), newVar("x"))
	exponentialDiff := exponential.diff()
	expectedAnswer := 17213.270665
	if val := exponentialDiff.eval(bindX(x0)); val != expectedAnswer {
		t.Errorf("Got %f when evaluating [%v].diff() at x=%f. Expected %f. Derivative: [%v]", val, exponential, x0,
			expectedAnswer, exponentialDiff.prune())
	}


	xToTheX := toGenericPower(newVar("x"), newVar("x"))
	xToTheXDiff := xToTheX.diff()
	expectedAnswer = 1336550.933484
	if val := xToTheXDiff.eval(bindX(x0)); val != expectedAnswer {
		t.Errorf("Got %f when evaluating [%v].diff() at x=%f. Expected %f. Derivative: [%v]", val, xToTheX, x0,
			expectedAnswer, xToTheXDiff.prune())
	}
}

func TestSineCosine(t *testing.T) {
	sine := newSine(newVar("x"))
	derivativeSine := sine.diff()

	x0 := math.Pi/2
	if val := derivativeSine.eval(bindX(x0)); val != 0.0 {
		t.Errorf("Got %f when evaluating [%v].diff() at x=%f. Expected %f. Derivative: [%v]", val, sine, x0,
			0.0, derivativeSine.prune())
	}

	cosine := newCosine(newVar("x"))
	derivativeCosine := cosine.diff()
	if val := derivativeCosine.eval(bindX(x0)); val != -1 {
		t.Errorf("Got %f when evaluating [%v].diff() at x=%f. Expected %f. Derivative: [%v]", val, cosine, x0,
			-1.0, derivativeCosine.prune())
	}


}

func bindX(x float64) map[string]float64 {
	return map[string]float64{"x": x}
}

