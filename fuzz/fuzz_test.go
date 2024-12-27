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

func FuzzSum(f *testing.F) {
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
		gotGV, ok := sumGV(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := sumCD(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if err != nil {
			t.Errorf("sumCD(%v, %v, %v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, fcoef, fscale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("sumGV(%v, %v, %v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, fcoef, fscale, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := sumSS(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if err != nil {
			t.Errorf("sumSS(%v, %v, %v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, fcoef, fscale, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("sumGV(%v, %v, %v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, fcoef, fscale, gotGV, wantSS)
		}
	})
}

func FuzzProd(f *testing.F) {
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
		gotGV, ok := prodGV(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := prodCD(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if err != nil {
			t.Errorf("prodCD(%v, %v, %v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, fcoef, fscale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("prodGV(%v, %v, %v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, fcoef, fscale, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := prodSS(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if err != nil {
			t.Errorf("prodSS(%v, %v, %v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, fcoef, fscale, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("prodGV(%v, %v, %v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, fcoef, fscale, gotGV, wantSS)
		}
	})
}

func FuzzMean(f *testing.F) {
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
		gotGV, ok := meanGV(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := meanCD(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if err != nil {
			t.Errorf("meanCD(%v, %v, %v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, fcoef, fscale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("meanGV(%v, %v, %v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, fcoef, fscale, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := meanSS(dcoef, dscale, ecoef, escale, fcoef, fscale)
		if err != nil {
			t.Errorf("meanSS(%v, %v, %v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, fcoef, fscale, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("meanGV(%v, %v, %v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, fcoef, fscale, gotGV, wantSS)
		}
	})
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

func FuzzDecimal_PowInt(f *testing.F) {
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
		gotGV, ok := powIntGV(dcoef, dscale, power)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach Db
		wantCD, err := powIntCD(dcoef, dscale, power)
		if err != nil {
			t.Errorf("powIntCD(%v, %v, %v) failed: %v", dcoef, dscale, power, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("powIntGV(%v, %v, %v) = %v, want %v", dcoef, dscale, power, gotGV, wantCD)
			return
		}
		// ShopSpring
		if dcoef == 0 {
			t.Skip()
			return
		}
		wantSS, err := powIntSS(dcoef, dscale, power)
		if err != nil {
			t.Errorf("powIntSS(%v, %v, %v) failed: %v", dcoef, dscale, power, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("powIntGV(%v, %v, %v) = %v, want %v", dcoef, dscale, power, gotGV, wantSS)
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

func FuzzDecimal_Log(f *testing.F) {
	ss.DivisionPrecision = 100
	ss.PowPrecisionNegativeExponent = 100
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		f.Add(d.coef, d.scale)
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int) {
		// GoValues
		gotGV, ok := logGV(dcoef, dscale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := logCD(dcoef, dscale)
		if err != nil {
			t.Errorf("logCD(%v, %v) failed: %v", dcoef, dscale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("logGV(%v, %v) = %v, want %v", dcoef, dscale, gotGV, wantCD)
			return
		}
		// ShopSpring
		wantSS, err := logSS(dcoef, dscale)
		if err != nil {
			t.Errorf("logSS(%v, %v) failed: %v", dcoef, dscale, err)
			return
		}
		if gotGV != wantSS {
			t.Errorf("logGV(%v, %v) = %v, want %v", dcoef, dscale, gotGV, wantSS)
		}
	})
}

func FuzzDecimal_Log2(f *testing.F) {
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		f.Add(d.coef, d.scale)
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int) {
		// GoValues
		gotGV, ok := log2GV(dcoef, dscale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := log2CD(dcoef, dscale)
		if err != nil {
			t.Errorf("log2CD(%v, %v) failed: %v", dcoef, dscale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("log2GV(%v, %v) = %v, want %v", dcoef, dscale, gotGV, wantCD)
			return
		}
		// ShopSpring
		// There is no log2 function.
	})
}

func FuzzDecimal_Log10(f *testing.F) {
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		f.Add(d.coef, d.scale)
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int) {
		// GoValues
		gotGV, ok := log10GV(dcoef, dscale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := log10CD(dcoef, dscale)
		if err != nil {
			t.Errorf("log10CD(%v, %v) failed: %v", dcoef, dscale, err)
			return
		}
		if gotGV != wantCD {
			t.Errorf("log10GV(%v, %v) = %v, want %v", dcoef, dscale, gotGV, wantCD)
			return
		}
		// ShopSpring
		// There is no log10 function.
	})
}

func FuzzDecimal_Pow(f *testing.F) {
	cd.BaseContext.Precision = 100
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, d := range corpus {
		for _, e := range corpus {
			f.Add(d.coef, d.scale, e.coef, e.scale)
		}
	}

	f.Fuzz(func(t *testing.T, dcoef int64, dscale int, ecoef int64, escale int) {
		if dcoef == 0 && ecoef == 0 {
			t.Skip()
			return
		}
		// GoValues
		gotGV, ok := powGV(dcoef, dscale, ecoef, escale)
		if !ok {
			t.Skip()
			return
		}
		// Cockroach DB
		wantCD, err := powCD(dcoef, dscale, ecoef, escale)
		if err != nil {
			if err.Error() == "exponent out of range" {
				t.Skip()
			} else {
				t.Errorf("powCD(%v, %v, %v, %v) failed: %v", dcoef, dscale, ecoef, escale, err)
			}
			return
		}
		if gotGV != wantCD {
			t.Errorf("powGV(%v, %v, %v, %v) = %v, want %v", dcoef, dscale, ecoef, escale, gotGV, wantCD)
			return
		}
		// ShopSpring
		// Unfortunately, ShopSpring just hungs in many cases.
		// For example, 1.000000000000000001^92233720368547758.07
	})
}

func powGV(dcoef int64, dscale int, ecoef int64, escale int) (string, bool) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", false
	}
	e, err := gv.New(ecoef, escale)
	if err != nil {
		return "", false
	}
	f, err := d.Pow(e)
	if err != nil {
		return "", false
	}
	return f.Trim(0).String(), true
}

func powCD(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Pow(f, d, e)
	if err != nil {
		return "", err
	}
	return roundCD(f)
}

func log2GV(dcoef int64, dscale int) (string, bool) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", false
	}
	f, err := d.Log2()
	if err != nil {
		return "", false
	}
	return f.Trim(0).String(), true
}

func log2CD(dcoef int64, dscale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Ln(f, d)
	if err != nil {
		return "", err
	}
	e := cd.New(2, 0)
	_, err = cd.BaseContext.Ln(e, e)
	if err != nil {
		return "", err
	}
	_, err = cd.BaseContext.Quo(f, f, e)
	if err != nil {
		return "", err
	}
	return roundCD(f)
}

func log10GV(dcoef int64, dscale int) (string, bool) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", false
	}
	f, err := d.Log10()
	if err != nil {
		return "", false
	}
	return f.Trim(0).String(), true
}

func log10CD(dcoef int64, dscale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Log10(f, d)
	if err != nil {
		return "", err
	}
	return roundCD(f)
}

func logGV(dcoef int64, dscale int) (string, bool) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", false
	}
	f, err := d.Log()
	if err != nil {
		return "", false
	}
	return f.Trim(0).String(), true
}

func logCD(dcoef int64, dscale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Ln(f, d)
	if err != nil {
		return "", err
	}
	return roundCD(f)
}

func logSS(dcoef int64, dscale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e, err := d.Ln(100)
	if err != nil {
		return "", err
	}
	return roundSS(e)
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

func expCD(dcoef int64, dscale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Exp(f, d)
	if err != nil {
		return "", err
	}
	return roundCD(f)
}

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

func sqrtCD(dcoef int64, dscale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	f := cd.New(0, 0)
	_, err := cd.BaseContext.Sqrt(f, d)
	if err != nil {
		return "", err
	}
	return roundCD(f)
}

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

func divSS(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Div(e)
	return roundSS(f)
}

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

func mulSS(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Mul(e)
	return roundSS(f)
}

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

func prodGV(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, bool) {
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
	g, err := gv.Prod(d, e, f)
	if err != nil {
		return "", false
	}
	return g.Trim(0).String(), true
}

func prodCD(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(fcoef, int32(-fscale))
	g := cd.New(0, 0)
	_, err := cd.BaseContext.Mul(g, d, e)
	if err != nil {
		return "", err
	}
	_, err = cd.BaseContext.Mul(g, g, f)
	if err != nil {
		return "", err
	}
	return roundCD(g)
}

func prodSS(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := ss.New(fcoef, int32(-fscale))
	g := d.Mul(e).Mul(f)
	return roundSS(g)
}

func addSS(dcoef int64, dscale int, ecoef int64, escale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := d.Add(e)
	return roundSS(f)
}

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

func sumGV(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, bool) {
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
	g, err := gv.Sum(d, e, f)
	if err != nil {
		return "", false
	}
	return g.Trim(0).String(), true
}

func sumCD(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(fcoef, int32(-fscale))
	g := cd.New(0, 0)
	_, err := cd.BaseContext.Add(g, d, e)
	if err != nil {
		return "", err
	}
	_, err = cd.BaseContext.Add(g, g, f)
	if err != nil {
		return "", err
	}
	return roundCD(g)
}

func sumSS(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := ss.New(fcoef, int32(-fscale))
	g := ss.Sum(d, e, f)
	return roundSS(g)
}

func meanGV(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, bool) {
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
	g, err := gv.Mean(d, e, f)
	if err != nil {
		return "", false
	}
	return g.Trim(0).String(), true
}

func meanCD(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, error) {
	d := cd.New(dcoef, int32(-dscale))
	e := cd.New(ecoef, int32(-escale))
	f := cd.New(fcoef, int32(-fscale))
	g := cd.New(0, 0)
	_, err := cd.BaseContext.Add(g, d, e)
	if err != nil {
		return "", err
	}
	_, err = cd.BaseContext.Add(g, g, f)
	if err != nil {
		return "", err
	}
	_, err = cd.BaseContext.Quo(g, g, cd.New(3, 0))
	if err != nil {
		return "", err
	}
	return roundCD(g)
}

func meanSS(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := ss.New(fcoef, int32(-fscale))
	g := ss.Sum(d, e, f)
	g = g.Div(ss.New(3, 0))
	return roundSS(g)
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

func addQuoSS(dcoef int64, dscale int, ecoef int64, escale int, fcoef int64, fscale int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e := ss.New(ecoef, int32(-escale))
	f := ss.New(fcoef, int32(-fscale))
	g := d.Add(e.Div(f))
	return roundSS(g)
}

func powIntGV(dcoef int64, dscale int, power int) (string, bool) {
	d, err := gv.New(dcoef, dscale)
	if err != nil {
		return "", false
	}
	f, err := d.PowInt(power)
	if err != nil {
		return "", false
	}
	return f.Trim(0).String(), true
}

func powIntCD(dcoef int64, dscale int, power int) (string, error) {
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

func powIntSS(dcoef int64, dscale int, power int) (string, error) {
	d := ss.New(dcoef, int32(-dscale))
	e, err := d.PowInt32(int32(power))
	if err != nil {
		return "", err
	}
	return roundSS(e)
}

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
