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

func FuzzDecimal_Add(f *testing.F) {
	ss.DivisionPrecision = 100
	ss.PowPrecisionNegativeExponent = 100
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for _, e := range corpus {
			f.Add(d.coef, d.scale, e.coef, e.scale)
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, ecoef int64, escale int) {
		// GoValues
		gotGV, ok := addGV(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := addCD(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("addCD(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("addGV(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := addSS(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("addSS(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("addGV(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotGV, wantSS)
		}
	})
}

func FuzzDecimal_Mul(f *testing.F) {
	ss.DivisionPrecision = 100
	ss.PowPrecisionNegativeExponent = 100
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for _, e := range corpus {
			f.Add(d.coef, d.scale, e.coef, e.scale)
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, ecoef int64, escale int) {
		// GoValues
		gotGV, ok := mulGV(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := mulCD(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("mulCD(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("mulGV(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := mulSS(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("mulSS(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("mulGV(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotGV, wantSS)
		}
	})
}

func FuzzDecimal_AddMul(f *testing.F) {
	ss.DivisionPrecision = 100
	ss.PowPrecisionNegativeExponent = 100
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for _, e := range corpus {
			for _, g := range corpus {
				f.Add(d.coef, d.scale, e.coef, e.scale, g.coef, g.scale)
			}
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) {
		// GoValues
		gotGV, ok := addMulGV(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := addMulCD(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if err != nil {
			t.Errorf("addMulCD(%v, %v, %v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, fcoef, fscale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("addMulGV(%v, %v, %v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, fcoef, fscale, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := addMulSS(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if err != nil {
			t.Errorf("addMulSS(%v, %v, %v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, fcoef, fscale, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("addMulGV(%v, %v, %v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, fcoef, fscale, gotGV, wantSS)
			return
		}
	})
}

func FuzzDecimal_AddQuo(f *testing.F) {
	ss.DivisionPrecision = 100
	ss.PowPrecisionNegativeExponent = 100
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for _, e := range corpus {
			for _, g := range corpus {
				f.Add(d.coef, d.scale, e.coef, e.scale, g.coef, g.scale)
			}
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) {
		// GoValues
		gotGV, ok := addQuoGV(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := addQuoCD(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if err != nil {
			t.Errorf("addQuoCD(%v, %v, %v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, fcoef, fscale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("addQuoGV(%v, %v, %v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, fcoef, fscale, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := addQuoSS(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if err != nil {
			t.Errorf("addQuoSS(%v, %v, %v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, fcoef, fscale, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("addQuoGV(%v, %v, %v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, fcoef, fscale, gotGV, wantSS)
			return
		}
	})
}

func FuzzDecimal_Quo(f *testing.F) {
	ss.DivisionPrecision = 100
	ss.PowPrecisionNegativeExponent = 100
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for _, e := range corpus {
			f.Add(d.coef, d.scale, e.coef, e.scale)
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, ecoef int64, escale int) {
		// GoValues
		gotGV, ok := quoGV(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := quoCD(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("quoCD(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("quoGV(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := divSS(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("divSS(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("quoGV(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotGV, wantSS)
			return
		}
	})
}

func FuzzDecimal_QuoRem(f *testing.F) {
	ss.DivisionPrecision = 100
	ss.PowPrecisionNegativeExponent = 100
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for _, e := range corpus {
			f.Add(d.coef, d.scale, e.coef, e.scale)
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, ecoef int64, escale int) {
		// GoValues
		gotQGV, gotRGV, ok := quoRemGV(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantQCD, wantRCD, err := quoRemCD(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("quoRemCD(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if gotQGV != wantQCD || gotRGV != wantRCD {
			t.Errorf("quoRemGV(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotQGV, wantQCD)
			return
		}
		// ShopSpring
		wantQSS, wantRSS, err := quoRemSS(dcoef, dscale, ecoef, escale)
		if err != nil {
			t.Errorf("quoRemSS(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			return
		}
		if gotQGV != wantQSS || gotRGV != wantRSS {
			t.Errorf("quoRemGV(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotQGV, wantQSS)
			return
		}
	})
}

func FuzzDecimal_Pow(f *testing.F) {
	ss.DivisionPrecision = 100
	ss.PowPrecisionNegativeExponent = 100
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for p := -10; p <= 10; p++ {
			f.Add(d.coef, d.scale, p)
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, power int) {
		// GoValues
		gotGV, ok := powGV(dcoef, dscale, power)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach Db
		wantCD, err := powCD(dcoef, dscale, power)
		if err != nil {
			t.Errorf("powCD(%v, %v, %v) failed: %v", dcoef, dscale, power, err)
			return
		}
		if c, err := cmpULP(gotGV, wantCD); err != nil {
			t.Errorf("cmpULP(%v, %v) failed: %v", gotGV, wantCD, err)
			return
		} else if c != 0 {
			t.Errorf("powGV(%v, %v, %v) = %v, want %v", dcoef, dscale, power, gotGV, wantCD)
			return
		}
		// ShopSpring
		if dcoef == 0 {
			t.Skip()
			return
		}
		wantSS, err := powSS(dcoef, dscale, power)
		if err != nil {
			t.Errorf("powSS(%v, %v, %v) failed: %v", dcoef, dscale, power, err)
			return
		}
		if c, err := cmpULP(gotGV, wantSS); err != nil {
			t.Errorf("cmpULP(%v, %v) failed: %v", gotGV, wantSS, err)
		} else if c != 0 {
			t.Errorf("powGV(%v, %v, %v) = %v, want %v", dcoef, dscale, power, gotGV, wantSS)
		}
	})
}

func FuzzDecimal_Sqrt(f *testing.F) {
	ss.DivisionPrecision = 100
	ss.PowPrecisionNegativeExponent = 100
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		f.Add(d.coef, d.scale)
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int) {
		// GoValues
		gotGV, ok := sqrtGV(dcoef, dscale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := sqrtCD(dcoef, dscale)
		if err != nil {
			t.Errorf("sqrtCD(%v, %v) failed: %v", dcoef, dscale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("sqrtGV(%v, %v) = %v, want %v", dcoef, dscale, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := sqrtSS(dcoef, dscale)
		if err != nil {
			t.Errorf("sqrtSS(%v, %v) failed: %v", dcoef, dscale, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("sqrtGV(%v, %v) = %v, want %v", dcoef, dscale, gotGV, wantSS)
		}
	})
}

func FuzzDecimal_Exp(f *testing.F) {
	ss.DivisionPrecision = 100
	ss.PowPrecisionNegativeExponent = 100
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		f.Add(d.coef, d.scale)
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int) {
		// GoValues
		gotGV, ok := expGV(dcoef, dscale)
		if !ok {
			t.Skip()
			return
		}
		// Prevent hunging
		if gotGV == "0" {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := expCD(dcoef, dscale)
		if err != nil {
			t.Errorf("expCD(%v, %v) failed: %v", dcoef, dscale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("expGV(%v, %v) = %v, want %v", dcoef, dscale, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := expSS(dcoef, dscale)
		if err != nil {
			t.Errorf("expSS(%v, %v) failed: %v", dcoef, dscale, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("expGV(%v, %v) = %v, want %v", dcoef, dscale, gotGV, wantSS)
		}
	})
}

func expGV(dcoef int64, dscale int) (string, bool) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", false
	}
	f, err := d.Exp()
	if err != nil {
		return "", false
	}
	return f.Trim(0).String(), true
}

//nolint:gosec
func expCD(dcoef int64, dscale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Exp(f, d)
	if err != nil {
		return "", err
	}
	return roundCD(f)
}

//nolint:gosec
func expSS(dcoef int64, dscale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e, err := d.ExpTaylor(100)
	if err != nil {
		return "", err
	}
	return roundSS(e)
}

func sqrtGV(dcoef int64, dscale int) (string, bool) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", false
	}
	f, err := d.Sqrt()
	if err != nil {
		return "", false
	}
	return f.Trim(0).String(), true
}

//nolint:gosec
func sqrtCD(dcoef int64, dscale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Sqrt(f, d)
	if err != nil {
		return "", err
	}
	return roundCD(f)
}

//nolint:gosec
func sqrtSS(dcoef int64, dscale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(5, -1)
	f, err := d.PowWithPrecision(e, 100)
	if err != nil {
		return "", err
	}
	return roundSS(f)
}

func quoGV(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
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

//nolint:gosec
func divSS(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Div(e)
	return roundSS(f)
}

//nolint:gosec
func quoCD(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Quo(f, d, e)
	if err != nil {
		return "", err
	}
	return roundCD(f)
}

func quoRemGV(dcoef int64, dscale int, ecoef int64, escale int) (string, string, bool) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", "", false
	}
	e, err := gv.New(ecoef, escale)
	if err != nil {
		return "", "", false
	}
	q, r, err := d.QuoRem(e)
	if err != nil {
		return "", "", false
	}
	return q.Trim(0).String(), r.Trim(0).String(), true
}

//nolint:gosec
func quoRemSS(dcoef int64, dscale int, ecoef int64, escale int) (string, string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	q, r := d.QuoRem(e, 0)
	qs, err := roundSS(q)
	if err != nil {
		return "", "", err
	}
	rs, err := roundSS(r)
	if err != nil {
		return "", "", err
	}
	return qs, rs, nil
}

//nolint:gosec
func quoRemCD(dcoef int64, dscale int, ecoef int64, escale int) (string, string, error) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	q := cd.New(0, 0)
	r := cd.New(0, 0)
	_, err := cd.BaseContext.QuoInteger(q, d, e)
	if err != nil {
		return "", "", err
	}
	_, err = cd.BaseContext.Rem(r, d, e)
	if err != nil {
		return "", "", err
	}
	qs, err := roundCD(q)
	if err != nil {
		return "", "", err
	}
	rs, err := roundCD(r)
	if err != nil {
		return "", "", err
	}
	return qs, rs, nil
}

//nolint:gosec
func mulSS(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Mul(e)
	return roundSS(f)
}

//nolint:gosec
func mulCD(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Mul(f, d, e)
	if err != nil {
		return "", err
	}
	return roundCD(f)
}

func mulGV(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
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

//nolint:gosec
func addSS(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Add(e)
	return roundSS(f)
}

//nolint:gosec
func addCD(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Add(f, d, e)
	if err != nil {
		return "", err
	}
	return roundCD(f)
}

func addGV(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
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

func addMulGV(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, bool) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", false
	}
	e, err := gv.New(ecoef, escale)
	if err != nil {
		return "", false
	}
	f, err := gv.New(fcoef, fscale)
	if err != nil {
		return "", false
	}
	g, err := d.AddMul(e, f)
	if err != nil {
		return "", false
	}
	return g.Trim(0).String(), true
}

//nolint:gosec
func addMulCD(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(fcoef, int32(-fscale))
	g := cd.New(0, 0)
	_, err := cd.BaseContext.Mul(g, e, f)
	if err != nil {
		return "", err
	}
	_, err = cd.BaseContext.Add(g, g, d)
	if err != nil {
		return "", err
	}
	return roundCD(g)
}

//nolint:gosec
func addMulSS(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := ss.New(fcoef, int32(-fscale))
	g := d.Add(e.Mul(f))
	return roundSS(g)
}

func addQuoGV(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, bool) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", false
	}
	e, err := gv.New(ecoef, escale)
	if err != nil {
		return "", false
	}
	f, err := gv.New(fcoef, fscale)
	if err != nil {
		return "", false
	}
	g, err := d.AddQuo(e, f)
	if err != nil {
		return "", false
	}
	return g.Trim(0).String(), true
}

//nolint:gosec
func addQuoCD(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(fcoef, int32(-fscale))
	g := cd.New(0, 0)
	_, err := cd.BaseContext.Quo(g, e, f)
	if err != nil {
		return "", err
	}
	_, err = cd.BaseContext.Add(g, g, d)
	if err != nil {
		return "", err
	}
	return roundCD(g)
}

//nolint:gosec
func addQuoSS(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := ss.New(fcoef, int32(-fscale))
	g := d.Add(e.Div(f))
	return roundSS(g)
}

func powGV(dcoef int64, dscale int, power int) (string, bool) {
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

//nolint:gosec
func powCD(dcoef int64, dscale int, power int) (string, error) {
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
	return roundCD(f)
}

//nolint:gosec
func powSS(dcoef int64, dscale int, power int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e, err := d.PowInt32(int32(power))
	if err != nil {
		return "", err
	}
	return roundSS(e)
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

//nolint:gosec
func roundSS(d ss.Decimal) (string, error) {
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

//nolint:gosec
func roundCD(d *cd.Decimal) (string, error) {
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
