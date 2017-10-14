package math_test

import (
	"fmt"
	gmath "math"

	"github.com/ericlagergren/decimal"
	"github.com/ericlagergren/decimal/math"
)

var one = decimal.New(1, 0)

type phiGenerator struct{ prec int32 }

func (p phiGenerator) Next() bool {
	return true
}

func (p phiGenerator) Term() math.Term {
	return math.Term{A: one, B: one}
}

func (p phiGenerator) Lentz() (f, Δ, C, D, eps *decimal.Big) {
	f = new(decimal.Big)
	Δ = new(decimal.Big)
	C = new(decimal.Big)
	D = new(decimal.Big)
	eps = decimal.New(1, p.prec)

	// Add a little extra precision to C and D so we get an "exact" result after
	// rounding.
	C.Context.SetPrecision(p.prec + 1)
	D.Context.SetPrecision(p.prec + 1)
	return f, Δ, C, D, eps
}

// Phi sets z to the golden ratio, φ, and returns z.
func Phi(z *decimal.Big) *decimal.Big {
	g := phiGenerator{prec: z.Context.Precision()}
	return math.Lentz(z, g)
}

// This example demonstrates using Lentz by calculating the golden ratio, φ.
func ExampleLentz_phi() {
	z := new(decimal.Big)
	Phi(z)
	p := (1 + gmath.Sqrt(5)) / 2

	fmt.Printf(`
Go     : %g
Decimal: %s`, p, z)
	// Output:
	//
	// Go     : 1.618033988749895
	// Decimal: 1.618033988749895
}
