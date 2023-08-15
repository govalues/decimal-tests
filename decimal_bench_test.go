package benchmarks

import (
	"strconv"
	"testing"

	cd "github.com/cockroachdb/apd/v3"
	gv "github.com/govalues/decimal"
	ss "github.com/shopspring/decimal"
)

var (
	resultString  string
	resultFloat64 float64
)

func BenchmarkDecimal_Add(b *testing.B) {
	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := gv.MustNew(2, 0)
			y := gv.MustNew(3, 0)
			_, _ = x.Add(y)
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
			x := gv.MustNew(2, 0)
			y := gv.MustNew(3, 0)
			_, _ = x.Mul(y)
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
			x := gv.MustNew(2, 0)
			y := gv.MustNew(4, 0)
			_, _ = x.Quo(y)
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
			x := gv.MustNew(2, 0)
			y := gv.MustNew(3, 0)
			_, _ = x.Quo(y) // implicitly calculates 38 digits and rounds to 19 digits
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
			x := gv.MustNew(11, 1)
			_, _ = x.Pow(60)
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
					resultString = d.String()
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
					resultString = d.Text('f')
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				d, err := ss.NewFromString(str)
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
	tests := []float64{
		123456789.12345678,
		123.456,
		1,
	}

	for _, f := range tests {
		str := strconv.FormatFloat(f, 'f', -1, 64)
		b.Run(str, func(b *testing.B) {
			b.Run("mod=govalues", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = gv.NewFromFloat64(f)
				}
			})

			b.Run("mod=cockroachdb", func(b *testing.B) {
				d := cd.New(0, 0)
				for i := 0; i < b.N; i++ {
					d.SetFloat64(f)
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = ss.NewFromFloat(f)
				}
			})
		})
	}
}

func BenchmarkDecimal_Float64(b *testing.B) {
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
					resultFloat64, _ = d.Float64()
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
					resultFloat64, _ = d.Float64()
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				d, err := ss.NewFromString(str)
				if err != nil {
					panic(err)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					resultFloat64, _ = d.Float64()
				}
			})
		})
	}
}
