package decimal_test

import (
	"fmt"
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

func roundShopspring(d ss.Decimal) (string, error) {
	// Check if number fits uint64 coefficient
	prec := int32(d.NumDigits())
	scale := int32(-d.Exponent())
	if prec-scale > gv.MaxScale {
		return "", fmt.Errorf("overflow (prec=%v, scale=%v)", prec, scale)
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
		return "", fmt.Errorf("overflow (prec=%v, scale=%v)", prec, scale)
	}
	return d.String(), nil
}

func roundCockroachdb(d *cd.Decimal) (string, error) {
	// Trailing Zeros
	d.Reduce(d)
	// Check if number fits uint64 coefficient
	prec := int32(d.NumDigits())
	scale := int32(-d.Exponent)
	if prec-scale > gv.MaxPrec {
		return "", fmt.Errorf("overflow (prec=%v, scale=%v)", prec, scale)
	}
	// Rounding
	switch {
	case scale >= prec && scale > gv.MaxScale: // no integer part
		scale = gv.MaxScale
		_, err := cd.BaseContext.Quantize(d, d, -scale)
		if err != nil {
			return "", err
		}
	case prec > scale && prec > gv.MaxPrec: // there is an integer part
		scale = scale - (prec - gv.MaxPrec)
		_, err := cd.BaseContext.Quantize(d, d, -scale)
		if err != nil {
			return "", err
		}
	}
	// Check if rounding added 1 extra digit
	prec = int32(d.NumDigits())
	scale = int32(-d.Exponent)
	if prec-scale > gv.MaxPrec {
		return "", fmt.Errorf("overflow (prec=%v, scale=%v)", prec, scale)
	}
	// Trailing Zeros
	d.Reduce(d)
	// Negative Zeros
	if d.IsZero() {
		d.Abs(d)
	}
	return d.Text('f'), nil
}

func quoGovalues(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", false
	}
	e, err := gv.New(ecoef, escale)
	if err != nil {
		return "", false
	}
	f, err := d.Quo(e)
	if err != nil {
		return "", false
	}
	return f.Trim(0).String(), true
}

func divShopspring(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Div(e)
	return roundShopspring(f)
}

func quoCockroachdb(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Quo(f, d, e)
	if err != nil {
		return "", err
	}
	return roundCockroachdb(f)
}

func FuzzDecimalQuo(f *testing.F) {
	ss.DivisionPrecision = 100
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for _, e := range corpus {
			f.Add(d.coef, d.scale, e.coef, e.scale)
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, ecoef int64, escale int) {
		// GoValues
		gotGV, ok := quoGovalues(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := quoCockroachdb(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("quoCockroachdb(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("quoGovalues(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := divShopspring(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("divShopspring(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("quoGovalues(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotGV, wantSS)
			return
		}
	})
}

func mulShopspring(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Mul(e)
	return roundShopspring(f)
}

func mulCockroachdb(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Mul(f, d, e)
	if err != nil {
		return "", err
	}
	return roundCockroachdb(f)
}

func mulGovalues(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", false
	}
	e, err := gv.New(ecoef, escale)
	if err != nil {
		return "", false
	}
	f, err := d.Mul(e)
	if err != nil {
		return "", false
	}
	return f.Trim(0).String(), true
}

func FuzzDecimalMul(f *testing.F) {
	ss.DivisionPrecision = 100
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for _, e := range corpus {
			f.Add(d.coef, d.scale, e.coef, e.scale)
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, ecoef int64, escale int) {
		// GoValues
		gotGV, ok := mulGovalues(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := mulCockroachdb(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("mulCockroachdb(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("mulGovalues(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := mulShopspring(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("mulShopspring(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("mulGovalues(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotGV, wantSS)
		}
	})
}

func addShopspring(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Add(e)
	return roundShopspring(f)
}

func addCockroachdb(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Add(f, d, e)
	if err != nil {
		return "", err
	}
	return roundCockroachdb(f)
}

func addGovalues(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", false
	}
	e, err := gv.New(ecoef, escale)
	if err != nil {
		return "", false
	}
	f, err := d.Add(e)
	if err != nil {
		return "", false
	}
	return f.Trim(0).String(), true
}

func FuzzDecimalAdd(f *testing.F) {
	ss.DivisionPrecision = 100
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for _, e := range corpus {
			f.Add(d.coef, d.scale, e.coef, e.scale)
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, ecoef int64, escale int) {
		// GoValues
		gotGV, ok := addGovalues(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := addCockroachdb(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("addCockroachdb(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("addGovalues(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := addShopspring(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("addShopspring(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("addGovalues(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotGV, wantSS)
		}
	})
}

func powGovalues(dcoef int64, dscale int, power int) (string, bool) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", false
	}
	f, err := d.Pow(power)
	if err != nil {
		return "", false
	}
	return f.Trim(0).String(), true
}

func powCockroachdb(dcoef int64, dscale int, power int) (string, error) {
	if dcoef == 0 && power == 0 {
		return "1", nil
	}
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(int64(power), 0)
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Pow(f, d, e)
	if err != nil {
		return "", err
	}
	return roundCockroachdb(f)
}

func powShopspring(dcoef int64, dscale int, power int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(int64(power), 0)
	f := d.Pow(e)
	return roundShopspring(f)
}

func FuzzDecimalPow(f *testing.F) {
	ss.DivisionPrecision = 100
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for p := -10; p <= 10; p++ {
			f.Add(d.coef, d.scale, p)
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, power int) {
		// GoValues
		gotGV, ok := powGovalues(dcoef, dscale, power)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach Db
		wantCD, err := powCockroachdb(dcoef, dscale, power)
		if err != nil {
			t.Errorf("powCockroachdb(%v, %v, %v) failed: %v", dcoef, dscale, power, err)
			return
		}
		if c, err := cmpULP(gotGV, wantCD); err != nil {
			t.Errorf("cmpULP(%v, %v) failed: %v", gotGV, wantCD, err)
			return
		} else if c != 0 {
			t.Errorf("powGovalues(%v, %v, %v) = %v, want %v", dcoef, dscale, power, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := powShopspring(dcoef, dscale, power)
		if err != nil {
			t.Errorf("powShopspring(%v, %v, %v) failed: %v", dcoef, dscale, power, err)
			return
		}
		if c, err := cmpULP(gotGV, wantSS); err != nil {
			t.Errorf("cmpULP(%v, %v) failed: %v", gotGV, wantSS, err)
		} else if c != 0 {
			t.Errorf("powGovalues(%v, %v, %v) = %v, want %v", dcoef, dscale, power, gotGV, wantSS)
		}
	})
}

// cmpULP compares decimals and returns 0 if they are within 1 ULP.
func cmpULP(s, t string) (int, error) {
	d, err := gv.Parse(s)
	if err != nil {
		return 0, err
	}
	e, err := gv.Parse(t)
	if err != nil {
		return 0, err
	}
	dist, err := d.SubAbs(e)
	if err != nil {
		return 0, err
	}
	ulp := d.ULP().Min(e.ULP())
	if dist.Cmp(ulp) <= 0 {
		return 0, nil
	}
	return d.Cmp(e), nil
}
