package math

import (
	"fmt"

	"github.com/ericlagergren/decimal"
)

func newDecimal(s string) *decimal.Big {
	x, ok := new(decimal.Big).SetString(s)
	if !ok {
		panic(fmt.Sprintf("bad input: %q", s))
	}
	return x
}

var Px = _Ln10

var (
	_E     = newDecimal("2.718281828459045235360287471352662497757247093699959574966967627724076630353547594571382178525166427")
	_Pi    = newDecimal("3.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117067")
	_Gamma = newDecimal("0.577215664901532860606512090082402431042159335939923598805767234884867726777664670936947063291746749")
	_Ln2   = newDecimal("0.693147180559945309417232121458176568075500134360255254120680009493393621969694715605863326996418687")
	_Ln10  = newDecimal("2.302585092994045684017991454684364207601101488628772976033327900967572609677352480235997205089598298")
)

// E sets z to the mathematical constant e.
func E(z *decimal.Big) *decimal.Big {
	prec := z.Context.Precision()
	if prec <= 100 {
		return z.Set(_E)
	}

	var (
		fac  = decimal.New(1, 0)
		incr = decimal.New(1, 0)
		sum  = decimal.New(2, 0)
		term = decimal.New(0, 0)
		prev = decimal.New(0, 0)
	)
	term.Context.SetPrecision(prec)

	for sum.Round(prec).Cmp(prev) != 0 {
		fac.Mul(fac, incr.Add(incr, one))
		prev.Copy(sum)
		sum.Add(sum, term.Quo(one, fac))
	}
	return sum
}

// Pi sets z to the mathematical constant π.
func Pi(z *decimal.Big) *decimal.Big {
	prec := z.Context.Precision()
	if prec <= 100 {
		return z.Set(_Pi)
	}

	var (
		lasts = decimal.New(0, 0)
		t     = decimal.New(3, 0)
		s     = decimal.New(3, 0)
		n     = decimal.New(1, 0)
		na    = decimal.New(0, 0)
		d     = decimal.New(0, 0)
		da    = decimal.New(24, 0)
	)
	lasts.Context.SetPrecision(prec)
	t.Context.SetPrecision(prec)

	for s.Round(prec).Cmp(lasts) != 0 {
		lasts.Set(s)
		n.Add(n, na)
		na.Add(na, eight)
		d.Add(d, da)
		da.Add(da, thirtyTwo)
		t.Mul(t, n)
		t.Quo(t, d)
		s.Add(s, t)
	}
	return s
}

/*
// Gamma sets z to the mathematical constant γ,
func Gamma(z *decimal.Big) *decimal.Big {
	prec := z.Context.Precision()
	if prec <= 100 {
		return z.Set(_Gamma)
	}

	// Antonino Machado Souza Filho and Georges Schwachheim. 1967.
	// Algorithm 309: Gamma function with arbitrary precision.
	// Commun. ACM 10, 8 (August 1967), 511-512.
	// DOI=http://dx.doi.org/10.1145/363534.363561

}

func loggamma(z, t *decimal.Big) *decimal.Big {
	var tmin *decimal.Big

	zcp := z.Context.Precision()
	if zcp >= 18 {
		tmin = decimal.New(int64(zcp), 0)
	} else {
		tmin = decimal.New(7, 0)
	}

	if t.Cmp(tmin) {
		return lgm(z, t)
	}

	f := new(decimal.Big).Copy(t)
	t0 := new(decimal.Big).Copy(t)

	for {
		t0.Add(t0, one)
		if t0.Comp(tmin) >= 0 {
			break
		}
		f.Mul(f, t0)
	}

	lgm(z, t0)

	tmp := z.Context.New(0, 0)
	Ln(tmp, f)

	return z.Sub(z, ln(tmp, f))
}

func lgm(z, w *decimal.Big) *decimal.Big {
	var c [20]*decimal.Big

	w0 := new(decimal.Big).Copy(w)
	den := new(decimal.Big).Copy(w) // den := w
	w2 := new(decimal.Big).Copy(w)  // w2 := w

	tmp := z.Context.New(0, 0)

	presum := new(decimal.Big)
	// presum := (w - .5) * ln(w) - w + const
	presum.Sub(w, ptFive)
	presum.Mul(presum, Ln(&tmp, w))
	presum.Sub(presum, tmp.Add(w, cnst))
}
*/
