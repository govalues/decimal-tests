package benchmarks

import (
	"math"
	"testing"

	cd "github.com/cockroachdb/apd/v3"
	gv "github.com/govalues/decimal"
	ss "github.com/shopspring/decimal"
)

var (
	corpus = []struct {
		scale int
		coef  int64
	}{
		{0, math.MaxInt64},
		{1, math.MaxInt64},
		{2, math.MaxInt64},
		{0, 5000000000000000000},
		{0, 1000000000000000000},
		{19, math.MaxInt64},
		{0, 9},
		{0, 7},
		{0, 6},
		{0, 3},
		{0, 2},
		{0, 1},
		{18, 3000000000000000003},
		{18, 3000000000000000000},
		{18, 2000000000000000002},
		{18, 2000000000000000000},
		{18, 1000000000000000001},
		{18, 1000000000000000000},
		{19, 3},
		{19, 2},
		{19, 1},
		{0, 0},
		{19, 0},
		{0, math.MinInt64},
		{1, math.MinInt64},
		{2, math.MinInt64},
		{0, -5000000000000000000},
		{0, -1000000000000000000},
		{19, math.MinInt64},
		{0, -9},
		{0, -7},
		{0, -6},
		{0, -3},
		{0, -2},
		{0, -1},
		{18, -3000000000000000003},
		{18, -3000000000000000000},
		{18, -2000000000000000002},
		{18, -2000000000000000000},
		{18, -1000000000000000001},
		{18, -1000000000000000000},
		{19, -3},
		{19, -2},
		{19, -1},
	}
)

func round_shopspring(d ss.Decimal) (string, bool) {
	// Check if number fits uint64 coefficient
	prec := int32(d.NumDigits())
	scale := int32(-d.Exponent())
	if prec-scale > gv.MaxScale {
		return "<overflow>", false
	}
	// Rounding
	switch {
	case scale >= prec && scale > gv.MaxScale: // no integer part
		scale = gv.MaxScale
		d = d.RoundBank(scale)
	case prec > scale && prec > gv.MaxPrec: // there is an integer part
		scale = scale - (prec - gv.MaxPrec)
		d = d.RoundBank(scale)
	}
	// Check if rounding added 1 extra digit
	prec = int32(d.NumDigits())
	scale = int32(-d.Exponent())
	if prec-scale > gv.MaxScale {
		return "<overflow>", false
	}
	return d.String(), true
}

func round_cockroachdb(d *cd.Decimal) (string, bool) {
	// Check if number fits uint64 coefficient
	prec := int32(d.NumDigits())
	scale := int32(-d.Exponent)
	if prec-scale > gv.MaxPrec {
		return "<overflow>", false
	}
	// Rounding
	switch {
	case scale >= prec && scale > gv.MaxScale: // no integer part
		scale = gv.MaxScale
		cd.BaseContext.Quantize(d, d, -scale)
	case prec > scale && prec > gv.MaxPrec: // there is an integer part
		scale = scale - (prec - gv.MaxPrec)
		cd.BaseContext.Quantize(d, d, -scale)
	}
	// Negative Zeros
	if d.IsZero() {
		d.Abs(d)
	}
	// Trailing Zeros
	d.Reduce(d)
	// Check if rounding added 1 extra digit
	prec = int32(d.NumDigits())
	scale = int32(-d.Exponent)
	if prec-scale > gv.MaxPrec {
		return "<overflow>", false
	}
	return d.Text('f'), true
}

func quo_govalues(dcoef int64, dscale int, ecoef int64, escale int) string {
	d := gv.New(dcoef, dscale)
	e := gv.New(ecoef, escale)
	f := d.Quo(e).Reduce()
	return f.String()
}

func div_shopspring(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Div(e)
	return round_shopspring(f)
}

func quo_cockroachdb(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(0, 0)
	cd.BaseContext.Quo(f, d, e)
	return round_cockroachdb(f)
}

func FuzzDecimal_Quo(f *testing.F) {

	ss.DivisionPrecision = 38
	cd.BaseContext.Precision = 38
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for _, e := range corpus {
			f.Add(d.coef, d.scale, e.coef, e.scale)
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, ecoef int64, escale int) {
		if dscale > gv.MaxScale || dscale < 0 {
			t.Skip()
			return
		}
		if escale > gv.MaxScale || escale < 0 {
			t.Skip()
			return
		}
		if ecoef == 0 {
			t.Skip()
			return
		}
		// Cockroach DB
		gotCD, ok := quo_cockroachdb(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		// GoValues
		wantGV := quo_govalues(dcoef, dscale, ecoef, escale)
		if wantGV != gotCD {
			t.Errorf("quo_cockroachdb(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotCD, wantGV)
		}
		// ShopSpring
		gotSS, ok := div_shopspring(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		if wantGV != gotSS {
			t.Errorf("div_shopspring(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotSS, wantGV)
		}
	})
}

func mul_shopspring(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Mul(e)
	return round_shopspring(f)
}

func mul_cockroachdb(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(0, 0)
	cd.BaseContext.Mul(f, d, e)
	return round_cockroachdb(f)
}

func mul_govalues(dcoef int64, dscale int, ecoef int64, escale int) string {
	d := gv.New(dcoef, dscale)
	e := gv.New(ecoef, escale)
	f := d.Mul(e).Reduce()
	return f.String()
}

func FuzzDecimal_Mul(f *testing.F) {

	ss.DivisionPrecision = 19
	cd.BaseContext.Precision = 38
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for _, e := range corpus {
			f.Add(d.coef, d.scale, e.coef, e.scale)
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, ecoef int64, escale int) {
		if dscale > gv.MaxScale || dscale < 0 {
			t.Skip()
			return
		}
		if escale > gv.MaxScale || escale < 0 {
			t.Skip()
			return
		}
		// Cockroach DB
		gotCD, ok := mul_cockroachdb(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		// GoValues
		wantGV := mul_govalues(dcoef, dscale, ecoef, escale)
		if wantGV != gotCD {
			t.Errorf("mul_cockroachdb(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotCD, wantGV)
		}
		// ShopSpring
		gotSS, ok := mul_shopspring(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		if wantGV != gotSS {
			t.Errorf("mul_shopspring(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotSS, wantGV)
		}
	})
}

func add_shopspring(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Add(e)
	return round_shopspring(f)
}

func add_cockroachdb(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(0, 0)
	cd.BaseContext.Add(f, d, e)
	return round_cockroachdb(f)
}

func add_govalues(dcoef int64, dscale int, ecoef int64, escale int) string {
	d := gv.New(dcoef, dscale)
	e := gv.New(ecoef, escale)
	f := d.Add(e).Reduce()
	return f.String()
}

func FuzzDecimal_Add(f *testing.F) {

	ss.DivisionPrecision = 19
	cd.BaseContext.Precision = 38
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for _, e := range corpus {
			f.Add(d.coef, d.scale, e.coef, e.scale)
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, ecoef int64, escale int) {
		if dscale > gv.MaxScale || dscale < 0 {
			t.Skip()
			return
		}
		if escale > gv.MaxScale || escale < 0 {
			t.Skip()
			return
		}
		// Cockroach DB
		gotCD, ok := add_cockroachdb(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		// GoValues
		wantGV := add_govalues(dcoef, dscale, ecoef, escale)
		if wantGV != gotCD {
			t.Errorf("add_cockroachdb(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotCD, wantGV)
		}
		// ShopSpring
		gotSS, ok := add_shopspring(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		if wantGV != gotSS {
			t.Errorf("add_shopspring(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotSS, wantGV)
		}
	})
}
