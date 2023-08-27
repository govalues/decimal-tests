package benchmarks

import (
	"math"
	"testing"

	cd "github.com/cockroachdb/apd/v3"
	gv "github.com/govalues/decimal"
	ss "github.com/shopspring/decimal"
)

var corpus = []struct {
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

func roundShopspring(d ss.Decimal) (string, bool) {
	// Check if number fits uint64 coefficient
	prec := int32(d.NumDigits())
	scale := int32(-d.Exponent())
	if prec-scale > gv.MaxScale {
		return "", false
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
		return "", false
	}
	return d.String(), true
}

func roundCockroachdb(d *cd.Decimal) (string, bool) {
	// Check if number fits uint64 coefficient
	prec := int32(d.NumDigits())
	scale := int32(-d.Exponent)
	if prec-scale > gv.MaxPrec {
		return "", false
	}
	// Rounding
	switch {
	case scale >= prec && scale > gv.MaxScale: // no integer part
		scale = gv.MaxScale
		_, err := cd.BaseContext.Quantize(d, d, -scale)
		if err != nil {
			return "", false
		}
	case prec > scale && prec > gv.MaxPrec: // there is an integer part
		scale = scale - (prec - gv.MaxPrec)
		_, err := cd.BaseContext.Quantize(d, d, -scale)
		if err != nil {
			return "", false
		}
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
		return "", false
	}
	return d.Text('f'), true
}

func quoGovalues(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", err
	}
	e, err := gv.New(ecoef, escale)
	if err != nil {
		return "", err
	}
	f, err := d.Quo(e)
	if err != nil {
		return "", err
	}
	return f.Trim(0).String(), nil
}

func divShopspring(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Div(e)
	return roundShopspring(f)
}

func quoCockroachdb(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Quo(f, d, e)
	if err != nil {
		return "", false
	}
	return roundCockroachdb(f)
}

func FuzzDecimalQuo(f *testing.F) {
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
		gotCD, ok := quoCockroachdb(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		// GoValues
		wantGV, err := quoGovalues(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("quo_govalues(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if wantGV != gotCD {
			t.Errorf("quo_cockroachdb(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotCD, wantGV)
		}
		// ShopSpring
		gotSS, ok := divShopspring(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		if wantGV != gotSS {
			t.Errorf("div_shopspring(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotSS, wantGV)
		}
	})
}

func mulShopspring(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Mul(e)
	return roundShopspring(f)
}

func mulCockroachdb(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Mul(f, d, e)
	if err != nil {
		return "", false
	}
	return roundCockroachdb(f)
}

func mulGovalues(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", err
	}
	e, err := gv.New(ecoef, escale)
	if err != nil {
		return "", err
	}
	f, err := d.Mul(e)
	if err != nil {
		return "", err
	}
	return f.Trim(0).String(), nil
}

func FuzzDecimalMul(f *testing.F) {
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
		gotCD, ok := mulCockroachdb(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		// GoValues
		wantGV, err := mulGovalues(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("mul_govalues(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if wantGV != gotCD {
			t.Errorf("mul_cockroachdb(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotCD, wantGV)
		}
		// ShopSpring
		gotSS, ok := mulShopspring(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		if wantGV != gotSS {
			t.Errorf("mul_shopspring(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotSS, wantGV)
		}
	})
}

func addShopspring(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Add(e)
	return roundShopspring(f)
}

func addCockroachdb(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Add(f, d, e)
	if err != nil {
		return "", false
	}
	return roundCockroachdb(f)
}

func addGovalues(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", err
	}
	e, err := gv.New(ecoef, escale)
	if err != nil {
		return "", err
	}
	f, err := d.Add(e)
	if err != nil {
		return "", err
	}
	return f.Trim(0).String(), nil
}

func FuzzDecimalAdd(f *testing.F) {
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
		gotCD, ok := addCockroachdb(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		// GoValues
		wantGV, err := addGovalues(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("add_govalues(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if wantGV != gotCD {
			t.Errorf("add_cockroachdb(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotCD, wantGV)
		}
		// ShopSpring
		gotSS, ok := addShopspring(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		if wantGV != gotSS {
			t.Errorf("add_shopspring(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotSS, wantGV)
		}
	})
}
