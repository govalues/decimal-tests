package benchmarks

import (
	"fmt"
	"math"
	"testing"

	cd "github.com/cockroachdb/apd"
	gv "github.com/govalues/decimal"
	ss "github.com/shopspring/decimal"
)

func BenchmarkDecimal_Add(b *testing.B) {

	b.Run("govalues/decimal", func(b *testing.B) {
		x := gv.New(2, 0)
		y := gv.New(3, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			x.Add(y)
		}
	})

	b.Run("shopspring/decimal", func(b *testing.B) {
		x := ss.New(2, 0)
		y := ss.New(3, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			x.Add(y)
		}
	})

	b.Run("cockroachdb/apd", func(b *testing.B) {
		cd.BaseContext.Precision = 19
		cd.BaseContext.Rounding = cd.RoundHalfEven
		x := cd.New(2, 0)
		y := cd.New(3, 0)
		z := cd.New(0, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cd.BaseContext.Add(z, x, y)
		}
	})
}

func BenchmarkDecimal_Mul(b *testing.B) {

	b.Run("govalues/decimal", func(b *testing.B) {
		x := gv.New(2, 0)
		y := gv.New(3, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			x.Mul(y)
		}
	})

	b.Run("shopspring/decimal", func(b *testing.B) {
		x := ss.New(2, 0)
		y := ss.New(3, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			x.Mul(y)
		}
	})

	b.Run("cockroachdb/apd", func(b *testing.B) {
		cd.BaseContext.Precision = 19
		cd.BaseContext.Rounding = cd.RoundHalfEven
		x := cd.New(2, 0)
		y := cd.New(3, 0)
		z := cd.New(0, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cd.BaseContext.Mul(z, x, y)
		}
	})
}

func BenchmarkDecimal_Pow(b *testing.B) {

	b.Run("govalues/decimal", func(b *testing.B) {
		x := gv.New(11, 1)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			x.Pow(60)
		}
	})

	b.Run("shopspring/decimal", func(b *testing.B) {
		x := ss.New(11, -1)
		y := ss.New(60, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			x.Pow(y)
		}
	})

	b.Run("cockroachdb/apd", func(b *testing.B) {
		cd.BaseContext.Precision = 19
		cd.BaseContext.Rounding = cd.RoundHalfEven
		x := cd.New(11, -1)
		y := cd.New(60, 0)
		z := cd.New(0, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cd.BaseContext.Pow(z, x, y)
		}
	})
}

func BenchmarkDecimal_QuoFinite(b *testing.B) {

	b.Run("govalues/decimal", func(b *testing.B) {
		x := gv.New(2, 0)
		y := gv.New(4, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			x.Quo(y)
		}
	})

	b.Run("shopspring/decimal", func(b *testing.B) {
		ss.DivisionPrecision = 38
		x := ss.New(2, 0)
		y := ss.New(4, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			x.Div(y).RoundBank(19)
		}
	})

	b.Run("cockroachdb/apd", func(b *testing.B) {
		cd.BaseContext.Precision = 38
		cd.BaseContext.Rounding = cd.RoundHalfEven
		x := cd.New(2, 0)
		y := cd.New(4, 0)
		z := cd.New(0, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cd.BaseContext.Quo(z, x, y)
		}
	})
}

func BenchmarkDecimal_QuoInfinite(b *testing.B) {

	b.Run("govalues/decimal", func(b *testing.B) {
		x := gv.New(2, 0)
		y := gv.New(3, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			x.Quo(y)
		}
	})

	b.Run("shopspring/decimal", func(b *testing.B) {
		ss.DivisionPrecision = 38
		x := ss.New(2, 0)
		y := ss.New(3, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			x.Div(y).RoundBank(19)
		}
	})

	b.Run("cockroachdb/apd", func(b *testing.B) {
		cd.BaseContext.Precision = 38
		cd.BaseContext.Rounding = cd.RoundHalfEven
		x := cd.New(2, 0)
		y := cd.New(3, 0)
		z := cd.New(0, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cd.BaseContext.Quo(z, x, y)
		}
	})
}

func BenchmarkDecimal_DailyInterestCalculation(b *testing.B) {

	b.Run("govalues/decimal", func(b *testing.B) {
		interest := gv.New(1000000000, 9) // = 1.000000000
		balance := gv.New(1000000, 2)     // = 10000.00
		yearlyRate := gv.New(1, 1)        // = 0.10
		daysInYear := gv.New(365, 0)      // = 365
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dailyRate := yearlyRate.Quo(daysInYear)
			interest.Add(balance.Mul(dailyRate).Round(9))
		}
	})

	b.Run("shopspring/decimal", func(b *testing.B) {
		interest := ss.New(1000000000, -9) // = 1.000000000
		balance := ss.New(1000000, -2)     // = 10000.00
		yearlyRate := ss.New(1, -1)        // = 0.10
		daysInYear := ss.New(365, 0)       // = 365
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dailyRate := yearlyRate.Div(daysInYear)
			interest.Add(balance.Mul(dailyRate).RoundBank(9))
		}
	})

	b.Run("cockroachdb/apd", func(b *testing.B) {
		cd.BaseContext.Precision = 19
		cd.BaseContext.Rounding = cd.RoundHalfEven
		interest := cd.New(1000000000, -9) // = 1.000000000
		balance := cd.New(1000000, -2)     // = 10000.00
		yearlyRate := cd.New(1, -1)        // = 0.10
		daysInYear := cd.New(365, 0)       // = 365
		dailyRate := cd.New(0, 0)
		factor := cd.New(0, 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cd.BaseContext.Quo(dailyRate, yearlyRate, daysInYear)
			cd.BaseContext.Mul(factor, balance, dailyRate)
			cd.BaseContext.Quantize(factor, factor, -9)
			cd.BaseContext.Add(factor, factor, interest)
		}
	})
}

/**********************************************************
* Fuzzing
**********************************************************/

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

func quo_govalues(xcoef int64, xscale int, ycoef int64, yscale int) (num string, scale int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	x := gv.New(xcoef, xscale)
	y := gv.New(ycoef, yscale)
	z := x.Quo(y).Reduce()
	return z.String(), z.Scale(), nil
}

func div_shopspring(xcoef int64, xscale int, ycoef int64, yscale int, zscale int) string {
	x := ss.New(xcoef, int32(-xscale))
	y := ss.New(ycoef, int32(-yscale))
	z := x.Div(y).RoundBank(int32(zscale))
	return z.String()
}

func quo_cockroachdb(xcoef int64, xscale int, ycoef int64, yscale int, zscale int) string {
	x := cd.New(xcoef, int32(-xscale))
	y := cd.New(ycoef, int32(-yscale))
	z := cd.New(0, 0)
	cd.BaseContext.Quo(z, x, y)
	cd.BaseContext.Quantize(z, z, int32(-zscale))
	z.Reduce(z)
	if z.IsZero() {
		z.Abs(z)
	}
	return z.Text('f')
}

func FuzzDecimal_Quo(f *testing.F) {

	ss.DivisionPrecision = 38
	cd.BaseContext.Precision = 38
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, x := range corpus {
		for _, y := range corpus {
			f.Add(x.coef, x.scale, y.coef, y.scale)
		}
	}

	f.Fuzz(func(t *testing.T, xcoef int64, xscale int, ycoef int64, yscale int) {
		// GoValues
		wantGV, scale, err := quo_govalues(xcoef, xscale, ycoef, yscale)
		if err != nil { // Error is expected if integer part of a result has more than 19 digits
			t.Skip()
			return
		}
		// ShopSpring
		gotSS := div_shopspring(xcoef, xscale, ycoef, yscale, scale)
		if wantGV != gotSS {
			t.Errorf("div_shopspring(%v, %v, %v, %v, %v) = %v, want %v", xcoef, xscale, ycoef, yscale, scale, gotSS, wantGV)
		}
		// Cockroach DB
		gotCD := quo_cockroachdb(xcoef, xscale, ycoef, yscale, scale)
		if wantGV != gotCD {
			t.Errorf("quo_cockroachdb(%v, %v, %v, %v, %v) = %v, want %v", xcoef, xscale, ycoef, yscale, scale, gotCD, wantGV)
		}
	})
}

func mul_govalues(xcoef int64, xscale int, ycoef int64, yscale int) (num string, scale int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	x := gv.New(xcoef, xscale)
	y := gv.New(ycoef, yscale)
	z := x.Mul(y).Reduce()
	return z.String(), z.Scale(), nil
}

func mul_shopspring(xcoef int64, xscale int, ycoef int64, yscale int, zscale int) string {
	x := ss.New(xcoef, int32(-xscale))
	y := ss.New(ycoef, int32(-yscale))
	z := x.Mul(y).RoundBank(int32(zscale))
	return z.String()
}

func mul_cockroachdb(xcoef int64, xscale int, ycoef int64, yscale int, zscale int) string {
	x := cd.New(xcoef, int32(-xscale))
	y := cd.New(ycoef, int32(-yscale))
	z := cd.New(0, 0)
	cd.BaseContext.Mul(z, x, y)
	cd.BaseContext.Quantize(z, z, int32(-zscale))
	z.Reduce(z)
	if z.IsZero() {
		z.Abs(z)
	}
	return z.Text('f')
}

func FuzzDecimal_Mul(f *testing.F) {

	ss.DivisionPrecision = 19
	cd.BaseContext.Precision = 38
	cd.BaseContext.Rounding = cd.RoundHalfEven

	for _, x := range corpus {
		for _, y := range corpus {
			f.Add(x.coef, x.scale, y.coef, y.scale)
		}
	}

	f.Fuzz(func(t *testing.T, xcoef int64, xscale int, ycoef int64, yscale int) {
		// GoValues
		wantGV, scale, err := mul_govalues(xcoef, xscale, ycoef, yscale)
		if err != nil { // Error is expected if integer part of a result has more than 19 digits
			t.Skip()
			return
		}
		// ShopSpring
		gotSS := mul_shopspring(xcoef, xscale, ycoef, yscale, scale)
		if wantGV != gotSS {
			t.Errorf("mul_shopspring(%v, %v, %v, %v, %v) = %v, want %v", xcoef, xscale, ycoef, yscale, scale, gotSS, wantGV)
		}
		// Cockroach DB
		gotCD := mul_cockroachdb(xcoef, xscale, ycoef, yscale, scale)
		if wantGV != gotCD {
			t.Errorf("mul_cockroachdb(%v, %v, %v, %v, %v) = %v, want %v", xcoef, xscale, ycoef, yscale, scale, gotCD, wantGV)
		}
	})
}
