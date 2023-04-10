package benchmarks

import (
	"testing"

	cd "github.com/cockroachdb/apd/v3"
	gv "github.com/govalues/decimal"
	ss "github.com/shopspring/decimal"
)

func BenchmarkDecimal_Add(b *testing.B) {

	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.New(2, 0)
			y := gv.New(3, 0)
			_ = x.Add(y)
		}
	})

	b.Run("mod=cockroachdb", func(b *testing.B) {
		cd.BaseContext.Precision = 19
		cd.BaseContext.Rounding = cd.RoundHalfEven
		for i := 0; i < b.N; i++ {
			x := cd.New(2, 0)
			y := cd.New(3, 0)
			z := cd.New(0, 0)
			cd.BaseContext.Add(z, x, y)
		}
	})

	b.Run("mod=shopspring", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := ss.New(2, 0)
			y := ss.New(3, 0)
			_ = x.Add(y)
		}
	})
}

func BenchmarkDecimal_Mul(b *testing.B) {

	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.New(2, 0)
			y := gv.New(3, 0)
			_ = x.Mul(y)
		}
	})

	b.Run("mod=cockroachdb", func(b *testing.B) {
		cd.BaseContext.Precision = 19
		cd.BaseContext.Rounding = cd.RoundHalfEven
		for i := 0; i < b.N; i++ {
			x := cd.New(2, 0)
			y := cd.New(3, 0)
			z := cd.New(0, 0)
			cd.BaseContext.Mul(z, x, y)
		}
	})

	b.Run("mod=shopspring", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := ss.New(2, 0)
			y := ss.New(3, 0)
			_ = x.Mul(y)
		}
	})
}

func BenchmarkDecimal_QuoFinite(b *testing.B) {

	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.New(2, 0)
			y := gv.New(4, 0)
			_ = x.Quo(y)
		}
	})

	b.Run("mod=cockroachdb", func(b *testing.B) {
		cd.BaseContext.Precision = 38
		cd.BaseContext.Rounding = cd.RoundHalfEven
		for i := 0; i < b.N; i++ {
			x := cd.New(2, 0)
			y := cd.New(4, 0)
			z := cd.New(0, 0)
			cd.BaseContext.Quo(z, x, y)
		}
	})

	b.Run("mod=shopspring", func(b *testing.B) {
		ss.DivisionPrecision = 38
		for i := 0; i < b.N; i++ {
			x := ss.New(2, 0)
			y := ss.New(4, 0)
			_ = x.Div(y)
		}
	})
}

func BenchmarkDecimal_QuoInfinite(b *testing.B) {

	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.New(2, 0)
			y := gv.New(3, 0)
			_ = x.Quo(y) // implicitly calculates 38 digits and rounds to 19 digits
		}
	})

	b.Run("mod=cockroachdb", func(b *testing.B) {
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

	b.Run("mod=shopspring", func(b *testing.B) {
		ss.DivisionPrecision = 38
		for i := 0; i < b.N; i++ {
			x := ss.New(2, 0)
			y := ss.New(3, 0)
			_ = x.Div(y).RoundBank(19)
		}
	})
}

func BenchmarkDecimal_Pow(b *testing.B) {

	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.New(11, 1)
			_ = x.Pow(60)
		}
	})

	b.Run("mod=cockroachdb", func(b *testing.B) {
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

	b.Run("mod=shopspring", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := ss.New(11, -1)
			y := ss.New(60, 0)
			_ = x.Pow(y).RoundBank(19)
		}
	})
}

func BenchmarkParse(b *testing.B) {

	tests := []string{
		"123456789.1234567890",
		"123.456",
		"1",
	}

	for _, str := range tests {

		b.Run(str, func(b *testing.B) {

			b.Run("mod=govalues", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = gv.Parse(str)
				}
			})

			b.Run("mod=cockroachdb", func(b *testing.B) {
				d := cd.New(0, 0)
				for i := 0; i < b.N; i++ {
					d.SetString(str)
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = ss.NewFromString(str)
				}
			})
		})
	}
}

func BenchmarkDecimal_String(b *testing.B) {

	tests := []string{
		"123456789.1234567890",
		"123.456",
		"1",
	}

	for _, str := range tests {

		b.Run(str, func(b *testing.B) {

			b.Run("mod=govalues", func(b *testing.B) {
				d, err := gv.Parse(str)
				if err != nil {
					panic(err)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = d.String()
				}
			})

			b.Run("mod=cockroachdb", func(b *testing.B) {
				d := cd.New(0, 0)
				d, _, err := d.SetString(str)
				if err != nil {
					panic(err)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = d.Text('f')
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				d, err := ss.NewFromString(str)
				if err != nil {
					panic(err)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = d.String()
				}
			})
		})
	}
}

func BenchmarkDecimal_DailyInterest(b *testing.B) {

	b.Run("mod=govalues", func(b *testing.B) {
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

	b.Run("mod=cockroachdb", func(b *testing.B) {
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

	b.Run("mod=shopspring", func(b *testing.B) {
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
}
