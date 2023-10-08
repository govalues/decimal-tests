package decimal_tests

import (
	"testing"

	cd "github.com/cockroachdb/apd/v3"
	gv "github.com/govalues/decimal"
	ss "github.com/shopspring/decimal"
)

var (
	resultString string
	resultFloat  float64
	resultErr    error
	resultSS     ss.Decimal
)

func BenchmarkDecimal_Add(b *testing.B) {
	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.MustNew(2, 0)
			y := gv.MustNew(3, 0)
			_, resultErr = x.Add(y)
		}
	})

	b.Run("mod=cockroachdb", func(b *testing.B) {
		cd.BaseContext.Precision = 19
		cd.BaseContext.Rounding = cd.RoundHalfEven
		for i := 0; i < b.N; i++ {
			x := cd.New(2, 0)
			y := cd.New(3, 0)
			z := cd.New(0, 0)
			_, resultErr = cd.BaseContext.Add(z, x, y)
		}
	})

	b.Run("mod=shopspring", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := ss.New(2, 0)
			y := ss.New(3, 0)
			resultSS = x.Add(y)
		}
	})
}

func BenchmarkDecimal_Mul(b *testing.B) {
	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.MustNew(2, 0)
			y := gv.MustNew(3, 0)
			_, resultErr = x.Mul(y)
		}
	})

	b.Run("mod=cockroachdb", func(b *testing.B) {
		cd.BaseContext.Precision = 19
		cd.BaseContext.Rounding = cd.RoundHalfEven
		for i := 0; i < b.N; i++ {
			x := cd.New(2, 0)
			y := cd.New(3, 0)
			z := cd.New(0, 0)
			_, resultErr = cd.BaseContext.Mul(z, x, y)
		}
	})

	b.Run("mod=shopspring", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := ss.New(2, 0)
			y := ss.New(3, 0)
			resultSS = x.Mul(y)
		}
	})
}

func BenchmarkDecimal_QuoFinite(b *testing.B) {
	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.MustNew(2, 0)
			y := gv.MustNew(4, 0)
			_, resultErr = x.Quo(y)
		}
	})

	b.Run("mod=cockroachdb", func(b *testing.B) {
		cd.BaseContext.Precision = 19
		cd.BaseContext.Rounding = cd.RoundHalfEven
		for i := 0; i < b.N; i++ {
			x := cd.New(2, 0)
			y := cd.New(4, 0)
			z := cd.New(0, 0)
			_, resultErr = cd.BaseContext.Quo(z, x, y)
		}
	})

	b.Run("mod=shopspring", func(b *testing.B) {
		ss.DivisionPrecision = 19
		for i := 0; i < b.N; i++ {
			x := ss.New(2, 0)
			y := ss.New(4, 0)
			resultSS = x.Div(y)
		}
	})
}

func BenchmarkDecimal_QuoInfinite(b *testing.B) {
	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.MustNew(2, 0)
			y := gv.MustNew(3, 0)
			_, resultErr = x.Quo(y) // implicitly calculates 38 digits and rounds to 19 digits
		}
	})

	b.Run("mod=cockroachdb", func(b *testing.B) {
		cd.BaseContext.Precision = 38
		cd.BaseContext.Rounding = cd.RoundHalfEven
		for i := 0; i < b.N; i++ {
			x := cd.New(2, 0)
			y := cd.New(3, 0)
			z := cd.New(0, 0)
			_, resultErr = cd.BaseContext.Quo(z, x, y)
			_, resultErr = cd.BaseContext.Quantize(z, z, -19)
		}
	})

	b.Run("mod=shopspring", func(b *testing.B) {
		ss.DivisionPrecision = 38
		for i := 0; i < b.N; i++ {
			x := ss.New(2, 0)
			y := ss.New(3, 0)
			resultSS = x.Div(y).RoundBank(19)
		}
	})
}

func BenchmarkDecimal_Pow(b *testing.B) {
	tests := map[string]struct {
		coef  int64
		scale int32
		power int64
	}{
		"1.1^60":     {11, 1, 60},
		"1.01^600":   {101, 2, 600},
		"1.001^6000": {1001, 3, 6000},
	}

	for name, tt := range tests {
		b.Run(name, func(b *testing.B) {
			b.Run("mod=govalues", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					x := gv.MustNew(tt.coef, int(tt.scale))
					_, resultErr = x.Pow(int(tt.power)) // implicitly calculates 38 digits and rounds to 19 digits
				}
			})

			b.Run("mod=cockroachdb", func(b *testing.B) {
				cd.BaseContext.Precision = 19
				cd.BaseContext.Rounding = cd.RoundHalfEven
				for i := 0; i < b.N; i++ {
					x := cd.New(tt.coef, -tt.scale)
					y := cd.New(tt.power, 0)
					z := cd.New(0, 0)
					_, resultErr = cd.BaseContext.Pow(z, x, y)
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				ss.DivisionPrecision = 19
				for i := 0; i < b.N; i++ {
					x := ss.New(tt.coef, -tt.scale)
					y := ss.New(tt.power, 0)
					resultSS = x.Pow(y).RoundBank(19)
				}
			})
		})
	}
}

func BenchmarkParse(b *testing.B) {
	tests := []string{
		"1",
		"123.456",
		"123456789.1234567890",
	}

	for _, s := range tests {
		b.Run(s, func(b *testing.B) {
			b.Run("mod=govalues", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, resultErr = gv.Parse(s)
				}
			})

			b.Run("mod=cockroachdb", func(b *testing.B) {
				d := cd.New(0, 0)
				for i := 0; i < b.N; i++ {
					_, _, resultErr = d.SetString(s)
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, resultErr = ss.NewFromString(s)
				}
			})
		})
	}
}

func BenchmarkDecimal_String(b *testing.B) {
	tests := []string{
		"1",
		"123.456",
		"123456789.1234567890",
	}

	for _, s := range tests {
		b.Run(s, func(b *testing.B) {
			b.Run("mod=govalues", func(b *testing.B) {
				d, err := gv.Parse(s)
				if err != nil {
					panic(err)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					resultString = d.String()
				}
			})

			b.Run("mod=cockroachdb", func(b *testing.B) {
				d := cd.New(0, 0)
				d, _, err := d.SetString(s)
				if err != nil {
					panic(err)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					resultString = d.Text('f')
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				d, err := ss.NewFromString(s)
				if err != nil {
					panic(err)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					resultString = d.String()
				}
			})
		})
	}
}

func BenchmarkNewFromFloat64(b *testing.B) {
	tests := map[string]float64{
		"1":                  1,
		"123.456":            123.456,
		"123456789.12345678": 123456789.12345678,
	}

	for name, f := range tests {
		b.Run(name, func(b *testing.B) {
			b.Run("mod=govalues", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, resultErr = gv.NewFromFloat64(f)
				}
			})

			b.Run("mod=cockroachdb", func(b *testing.B) {
				d := cd.New(0, 0)
				for i := 0; i < b.N; i++ {
					_, resultErr = d.SetFloat64(f)
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					resultSS = ss.NewFromFloat(f)
				}
			})
		})
	}
}

func BenchmarkDecimal_Float64(b *testing.B) {
	tests := []string{
		"1",
		"123.456",
		"123456789.1234567890",
	}

	for _, s := range tests {
		b.Run(s, func(b *testing.B) {
			b.Run("mod=govalues", func(b *testing.B) {
				d, err := gv.Parse(s)
				if err != nil {
					panic(err)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					resultFloat, _ = d.Float64()
				}
			})

			b.Run("mod=cockroachdb", func(b *testing.B) {
				d := cd.New(0, 0)
				d, _, err := d.SetString(s)
				if err != nil {
					panic(err)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					resultFloat, resultErr = d.Float64()
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				d, err := ss.NewFromString(s)
				if err != nil {
					panic(err)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					resultFloat, _ = d.Float64()
				}
			})
		})
	}
}
