package benchmarks

import (
	"math"
	"testing"

	cd "github.com/cockroachdb/apd/v3"
	gv "github.com/govalues/decimal"
	ss "github.com/shopspring/decimal"
)

func BenchmarkDecimal_Add(b *testing.B) {

	b.Run("govalues/decimal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.New(2, 0)
			y := gv.New(3, 0)
			_ = x.Add(y)
		}
	})

	b.Run("shopspring/decimal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := ss.New(2, 0)
			y := ss.New(3, 0)
			_ = x.Add(y)
		}
	})

	b.Run("cockroachdb/apd", func(b *testing.B) {
		cd.BaseContext.Precision = 19
		cd.BaseContext.Rounding = cd.RoundHalfEven
		for i := 0; i < b.N; i++ {
			x := cd.New(2, 0)
			y := cd.New(3, 0)
			z := cd.New(0, 0)
			cd.BaseContext.Add(z, x, y)
		}
	})
}

func BenchmarkDecimal_Mul(b *testing.B) {

	b.Run("govalues/decimal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.New(2, 0)
			y := gv.New(3, 0)
			_ = x.Mul(y)
		}
	})

	b.Run("shopspring/decimal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := ss.New(2, 0)
			y := ss.New(3, 0)
			_ = x.Mul(y)
		}
	})

	b.Run("cockroachdb/apd", func(b *testing.B) {
		cd.BaseContext.Precision = 19
		cd.BaseContext.Rounding = cd.RoundHalfEven
		for i := 0; i < b.N; i++ {
			x := cd.New(2, 0)
			y := cd.New(3, 0)
			z := cd.New(0, 0)
			cd.BaseContext.Mul(z, x, y)
		}
	})
}

func BenchmarkDecimal_FMA(b *testing.B) {

	b.Run("govalues/decimal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.New(2, 0)
			y := gv.New(3, 0)
			z := gv.New(4, 0)
			_ = x.FMA(y, z)
		}
	})
}

func BenchmarkDecimal_Pow(b *testing.B) {

	b.Run("govalues/decimal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.New(11, 1)
			_ = x.Pow(60)
		}
	})

	b.Run("shopspring/decimal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := ss.New(11, -1)
			y := ss.New(60, 0)
			_ = x.Pow(y).RoundBank(19)
		}
	})

	b.Run("cockroachdb/apd", func(b *testing.B) {
		cd.BaseContext.Precision = 38
		cd.BaseContext.Rounding = cd.RoundHalfEven
		for i := 0; i < b.N; i++ {
			x := cd.New(11, -1)
			y := cd.New(60, 0)
			z := cd.New(0, 0)
			cd.BaseContext.Pow(z, x, y)
			cd.BaseContext.Quantize(z, z, -19)
		}
	})
}

func BenchmarkDecimal_QuoFinite(b *testing.B) {

	b.Run("govalues/decimal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.New(2, 0)
			y := gv.New(4, 0)
			_ = x.Quo(y)
		}
	})

	b.Run("shopspring/decimal", func(b *testing.B) {
		ss.DivisionPrecision = 38
		for i := 0; i < b.N; i++ {
			x := ss.New(2, 0)
			y := ss.New(4, 0)
			_ = x.Div(y)
		}
	})

	b.Run("cockroachdb/apd", func(b *testing.B) {
		cd.BaseContext.Precision = 38
		cd.BaseContext.Rounding = cd.RoundHalfEven
		for i := 0; i < b.N; i++ {
			x := cd.New(2, 0)
			y := cd.New(4, 0)
			z := cd.New(0, 0)
			cd.BaseContext.Quo(z, x, y)
		}
	})
}

func BenchmarkDecimal_QuoInfinite(b *testing.B) {

	b.Run("govalues/decimal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.New(2, 0)
			y := gv.New(3, 0)
			_ = x.Quo(y) // implicitly calculates 38 digits and rounds to 19 digits
		}
	})

	b.Run("shopspring/decimal", func(b *testing.B) {
		ss.DivisionPrecision = 38
		for i := 0; i < b.N; i++ {
			x := ss.New(2, 0)
			y := ss.New(3, 0)
			_ = x.Div(y).RoundBank(19)
		}
	})

	b.Run("cockroachdb/apd", func(b *testing.B) {
		cd.BaseContext.Precision = 38
		cd.BaseContext.Rounding = cd.RoundHalfEven
		for i := 0; i < b.N; i++ {
			x := cd.New(2, 0)
			y := cd.New(3, 0)
			z := cd.New(0, 0)
			cd.BaseContext.Quo(z, x, y)
			cd.BaseContext.Quantize(z, z, -19)
		}
	})
}

func BenchmarkDecimal_DailyInterestCalculation(b *testing.B) {

	b.Run("govalues/decimal", func(b *testing.B) {
		interest := gv.New(1000000000, 9) // = 1.000000000
		balance := gv.New(1000000, 2)     // = 10000.00
		yearlyRate := gv.New(1, 1)        // = 0.10
		daysInYear := gv.New(365, 0)      // = 365
		dailyRate := gv.New(0, 0)         // = 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dailyRate = yearlyRate.Quo(daysInYear)
			_ = interest.Add(balance.Mul(dailyRate).Round(9))
		}
	})

	b.Run("shopspring/decimal", func(b *testing.B) {
		ss.DivisionPrecision = 38
		interest := ss.New(1000000000, -9) // = 1.000000000
		balance := ss.New(1000000, -2)     // = 10000.00
		yearlyRate := ss.New(1, -1)        // = 0.10
		daysInYear := ss.New(365, 0)       // = 365
		dailyRate := ss.New(0, 0)          // = 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dailyRate = yearlyRate.Div(daysInYear).RoundBank(19)
			_ = interest.Add(balance.Mul(dailyRate).RoundBank(9))
		}
	})

	b.Run("cockroachdb/apd", func(b *testing.B) {
		cd.BaseContext.Precision = 38
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
			cd.BaseContext.Quantize(dailyRate, dailyRate, -19)
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

func Decimal_FMA(f *testing.F) {
	f.Error("not implemented")
}
