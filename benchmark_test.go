package decimal_test

import (
	"encoding/binary"
	"io"
	"os"
	"testing"

	cd "github.com/cockroachdb/apd/v3"
	gv "github.com/govalues/decimal"
	ss "github.com/shopspring/decimal"
)

var (
	resultString  string
	resultError   error
	resultDecimal ss.Decimal
)

func BenchmarkDecimal_Add(b *testing.B) {
	tests := map[string]struct {
		xcoef  int64
		xscale int32
		ycoef  int64
		yscale int32
	}{
		"5+6":                           {5, 0, 6, 0},
		"5.000000+6.000000":             {5000000, 6, 6000000, 6},
		"5.000000000000+6.000000000000": {5000000000000, 12, 6000000000000, 12},
		"5.000000000000000000+6.000000000000000000": {5000000000000000000, 18, 6000000000000000000, 18},
	}
	for name, tt := range tests {
		b.Run(name, func(b *testing.B) {
			b.Run("mod=govalues", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					x := gv.MustNew(tt.xcoef, int(tt.xscale))
					y := gv.MustNew(tt.ycoef, int(tt.yscale))
					_, resultError = x.Add(y)
				}
			})

			b.Run("mod=cockroachdb", func(b *testing.B) {
				cd.BaseContext.Precision = 19
				cd.BaseContext.Rounding = cd.RoundHalfEven
				for i := 0; i < b.N; i++ {
					x := cd.New(tt.xcoef, -tt.xscale)
					y := cd.New(tt.ycoef, -tt.yscale)
					z := cd.New(0, 0)
					_, resultError = cd.BaseContext.Add(z, x, y)
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					x := ss.New(tt.xcoef, -tt.xscale)
					y := ss.New(tt.ycoef, -tt.yscale)
					resultDecimal = x.Add(y)
				}
			})
		})
	}
}

func BenchmarkDecimal_Mul(b *testing.B) {
	tests := map[string]struct {
		xcoef  int64
		xscale int32
		ycoef  int64
		yscale int32
	}{
		"2*3":                           {2, 0, 3, 0},
		"2.000000*3.000000":             {2000000, 6, 3000000, 6},
		"2.000000000000*3.000000000000": {2000000000000, 12, 3000000000000, 12},
		"2.000000000000000000*3.000000000000000000": {2000000000000000000, 18, 3000000000000000000, 18},
	}

	for name, tt := range tests {
		b.Run(name, func(b *testing.B) {
			b.Run("mod=govalues", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					x := gv.MustNew(tt.xcoef, int(tt.xscale))
					y := gv.MustNew(tt.ycoef, int(tt.yscale))
					_, resultError = x.Mul(y)
				}
			})

			b.Run("mod=cockroachdb", func(b *testing.B) {
				cd.BaseContext.Precision = 19
				cd.BaseContext.Rounding = cd.RoundHalfEven
				for i := 0; i < b.N; i++ {
					x := cd.New(tt.xcoef, -tt.xscale)
					y := cd.New(tt.ycoef, -tt.yscale)
					z := cd.New(0, 0)
					_, resultError = cd.BaseContext.Mul(z, x, y)
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					x := ss.New(tt.xcoef, -tt.xscale)
					y := ss.New(tt.ycoef, -tt.yscale)
					resultDecimal = x.Mul(y)
				}
			})
		})
	}
}

func BenchmarkDecimal_QuoExact(b *testing.B) {
	tests := map[string]struct {
		xcoef  int64
		xscale int32
		ycoef  int64
		yscale int32
	}{
		"2÷4":                           {2, 0, 4, 0},
		"2.000000÷4.000000":             {2000000, 6, 4000000, 6},
		"2.000000000000÷4.000000000000": {2000000000000, 12, 4000000000000, 12},
		"2.000000000000000000÷4.000000000000000000": {2000000000000000000, 18, 4000000000000000000, 18},
	}
	for name, tt := range tests {
		b.Run(name, func(b *testing.B) {
			b.Run("mod=govalues", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					x := gv.MustNew(tt.xcoef, int(tt.xscale))
					y := gv.MustNew(tt.ycoef, int(tt.yscale))
					_, resultError = x.Quo(y)
				}
			})

			b.Run("mod=cockroachdb", func(b *testing.B) {
				cd.BaseContext.Precision = 19
				cd.BaseContext.Rounding = cd.RoundHalfEven
				for i := 0; i < b.N; i++ {
					x := cd.New(tt.xcoef, -tt.xscale)
					y := cd.New(tt.ycoef, -tt.yscale)
					z := cd.New(0, 0)
					_, resultError = cd.BaseContext.Quo(z, x, y)
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				ss.DivisionPrecision = 19
				for i := 0; i < b.N; i++ {
					x := ss.New(tt.xcoef, -tt.xscale)
					y := ss.New(tt.ycoef, -tt.yscale)
					resultDecimal = x.Div(y)
				}
			})
		})
	}
}

func BenchmarkDecimal_QuoInfinite(b *testing.B) {
	tests := map[string]struct {
		xcoef  int64
		xscale int32
		ycoef  int64
		yscale int32
	}{
		"2÷3":                           {2, 0, 3, 0},
		"2.000000÷3.000000":             {2000000, 6, 3000000, 6},
		"2.000000000000÷3.000000000000": {2000000000000, 12, 3000000000000, 12},
		"2.000000000000000000÷3.000000000000000000": {2000000000000000000, 18, 3000000000000000000, 18},
	}
	for name, tt := range tests {
		b.Run(name, func(b *testing.B) {
			b.Run("mod=govalues", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					x := gv.MustNew(tt.xcoef, int(tt.xscale))
					y := gv.MustNew(tt.ycoef, int(tt.yscale))
					_, resultError = x.Quo(y) // implicitly calculates 38 digits and rounds to 19 digits
				}
			})

			b.Run("mod=cockroachdb", func(b *testing.B) {
				cd.BaseContext.Precision = 38
				cd.BaseContext.Rounding = cd.RoundHalfEven
				for i := 0; i < b.N; i++ {
					x := cd.New(tt.xcoef, -tt.xscale)
					y := cd.New(tt.ycoef, -tt.yscale)
					z := cd.New(0, 0)
					_, resultError = cd.BaseContext.Quo(z, x, y)
					_, resultError = cd.BaseContext.Quantize(z, z, -19)
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				ss.DivisionPrecision = 38
				for i := 0; i < b.N; i++ {
					x := ss.New(tt.xcoef, -tt.xscale)
					y := ss.New(tt.ycoef, -tt.yscale)
					resultDecimal = x.Div(y).RoundBank(19)
				}
			})
		})
	}
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
					_, resultError = x.Pow(int(tt.power)) // implicitly calculates 38 digits and rounds to 19 digits
				}
			})

			b.Run("mod=cockroachdb", func(b *testing.B) {
				cd.BaseContext.Precision = 19
				cd.BaseContext.Rounding = cd.RoundHalfEven
				for i := 0; i < b.N; i++ {
					x := cd.New(tt.coef, -tt.scale)
					y := cd.New(tt.power, 0)
					z := cd.New(0, 0)
					_, resultError = cd.BaseContext.Pow(z, x, y)
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				ss.DivisionPrecision = 19
				for i := 0; i < b.N; i++ {
					x := ss.New(tt.coef, -tt.scale)
					y := ss.New(tt.power, 0)
					resultDecimal = x.Pow(y).RoundBank(19)
				}
			})
		})
	}
}

func BenchmarkParse(b *testing.B) {
	tests := []string{
		"1",
		"12345.12345",
		"123456789.1234567890",
		"1234567890123456789.12345678901234567890",
	}
	for _, s := range tests {
		b.Run(s, func(b *testing.B) {
			b.Run("mod=govalues", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, resultError = gv.Parse(s)
				}
			})

			b.Run("mod=cockroachdb", func(b *testing.B) {
				d := cd.New(0, 0)
				for i := 0; i < b.N; i++ {
					_, _, resultError = d.SetString(s)
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, resultError = ss.NewFromString(s)
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
					b.Fatal(err)
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
					b.Fatal(err)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					resultString = d.Text('f')
				}
			})

			b.Run("mod=shopspring", func(b *testing.B) {
				d, err := ss.NewFromString(s)
				if err != nil {
					b.Fatal(err)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					resultString = d.String()
				}
			})
		})
	}
}

func readTelcoTests() ([]int64, error) {
	file, err := os.Open("expon180.1e6b")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data := make([]int64, 0, 1000000)
	buf := make([]byte, 8)
	for {
		_, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		num := binary.BigEndian.Uint64(buf)
		data = append(data, int64(num))
	}
	return data, nil
}

// BenchmarkDecimal_Telco implements computational part of "[Telco benchmark]"
// by Mike Cowlishaw.
// I/O part is not implemented.
//
// [Telco benchmark]: https://speleotrove.com/decimal/telco.html
func BenchmarkDecimal_Telco(b *testing.B) {
	tests, err := readTelcoTests()
	if err != nil {
		b.Fatal(err)
	}

	b.Run("mod=govalues", func(b *testing.B) {
		totalFinalPrice := gv.Zero
		totalBaseTax := gv.Zero
		totalDistTax := gv.Zero
		baseRate := gv.MustParse("0.0013")
		distRate := gv.MustParse("0.00894")
		baseTaxRate := gv.MustParse("0.0675")
		distTaxRate := gv.MustParse("0.0341")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var err error
			tt := tests[i%len(tests)]
			callType := tt & 0x01

			// Duration, Seconds
			duration := gv.MustNew(tt, 0)

			// Price
			var price gv.Decimal
			if callType == 0 {
				price, err = duration.Mul(baseRate)
			} else {
				price, err = duration.Mul(distRate)
			}
			if err != nil {
				b.Fatal(err)
			}
			price = price.Round(2)

			// Base Tax
			baseTax, err := price.Mul(baseTaxRate)
			if err != nil {
				b.Fatal(err)
			}
			baseTax = baseTax.Trunc(2)
			totalBaseTax, err = totalBaseTax.Add(baseTax)
			if err != nil {
				b.Fatal(err)
			}
			finalPrice, err := price.Add(baseTax)
			if err != nil {
				b.Fatal(err)
			}

			// Distance Tax
			if callType != 0 {
				distTax, err := price.Mul(distTaxRate)
				if err != nil {
					b.Fatal(err)
				}
				distTax = distTax.Trunc(2)
				totalDistTax, err = totalDistTax.Add(distTax)
				if err != nil {
					b.Fatal(err)
				}
				finalPrice, err = finalPrice.Add(distTax)
				if err != nil {
					b.Fatal(err)
				}
			}

			// Final Price
			totalFinalPrice, err = totalFinalPrice.Add(finalPrice)
			if err != nil {
				b.Fatal(err)
			}
			resultString = finalPrice.String()
		}
	})

	b.Run("mod=cockroachdb", func(b *testing.B) {
		cd.BaseContext.Precision = 19
		cd.BaseContext.Rounding = cd.RoundHalfEven
		totalFinalPrice := cd.New(0, 0)
		totalBaseTax := cd.New(0, 0)
		totalDistTax := cd.New(0, 0)
		baseRate := cd.New(13, -4)     // 0.0013
		distRate := cd.New(894, -5)    // 0.00894
		baseTaxRate := cd.New(675, -4) // 0.0675
		distTaxRate := cd.New(341, -4) // 0.0341
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var err error
			tt := tests[i%len(tests)]
			callType := tt & 0x01

			// Duration, Seconds
			duration := cd.New(tt, 0)

			// Price
			price := new(cd.Decimal)
			if callType == 0 {
				_, err = cd.BaseContext.Mul(price, duration, baseRate)
			} else {
				_, err = cd.BaseContext.Mul(price, duration, distRate)
			}
			if err != nil {
				b.Fatal(err)
			}
			cd.BaseContext.Rounding = cd.RoundHalfEven
			_, err = cd.BaseContext.Quantize(price, price, -2)
			if err != nil {
				b.Fatal(err)
			}

			// Base Tax
			baseTax := new(cd.Decimal)
			_, err = cd.BaseContext.Mul(baseTax, price, baseTaxRate)
			if err != nil {
				b.Fatal(err)
			}
			cd.BaseContext.Rounding = cd.RoundDown
			_, err = cd.BaseContext.Quantize(baseTax, baseTax, -2)
			if err != nil {
				b.Fatal(err)
			}
			_, err = cd.BaseContext.Add(totalBaseTax, totalBaseTax, baseTax)
			if err != nil {
				b.Fatal(err)
			}
			finalPrice := new(cd.Decimal)
			_, err = cd.BaseContext.Add(finalPrice, price, baseTax)
			if err != nil {
				b.Fatal(err)
			}

			// Distance Tax
			if callType != 0 {
				distTax := new(cd.Decimal)
				_, err = cd.BaseContext.Mul(distTax, price, distTaxRate)
				if err != nil {
					b.Fatal(err)
				}
				cd.BaseContext.Rounding = cd.RoundDown
				_, err = cd.BaseContext.Quantize(distTax, distTax, -2)
				if err != nil {
					b.Fatal(err)
				}
				_, err = cd.BaseContext.Add(totalDistTax, totalDistTax, distTax)
				if err != nil {
					b.Fatal(err)
				}
				_, err = cd.BaseContext.Add(finalPrice, finalPrice, distTax)
				if err != nil {
					b.Fatal(err)
				}
			}

			// Final Price
			_, err = cd.BaseContext.Add(totalFinalPrice, totalFinalPrice, finalPrice)
			if err != nil {
				b.Fatal(err)
			}
			resultString = finalPrice.String()
		}
	})

	b.Run("mod=shopspring", func(b *testing.B) {
		totalFinalPrice := ss.Zero
		totalBaseTax := ss.Zero
		totalDistTax := ss.Zero
		baseRate := ss.RequireFromString("0.0013")
		distRate := ss.RequireFromString("0.00894")
		baseTaxRate := ss.RequireFromString("0.0675")
		distTaxRate := ss.RequireFromString("0.0341")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tt := tests[i%len(tests)]
			callType := tt & 0x01

			// Duration, Seconds
			duration := ss.NewFromInt(tt)

			// Price
			var price ss.Decimal
			if callType == 0 {
				price = duration.Mul(baseRate)
			} else {
				price = duration.Mul(distRate)
			}
			price = price.RoundBank(2)

			// Base Tax
			baseTax := price.Mul(baseTaxRate)
			baseTax = baseTax.RoundDown(2)
			totalBaseTax = totalBaseTax.Add(baseTax)
			finalPrice := price.Add(baseTax)

			// Distance Tax
			if callType != 0 {
				distTax := price.Mul(distTaxRate)
				distTax = distTax.RoundDown(2)
				totalDistTax = totalDistTax.Add(distTax)
				finalPrice = finalPrice.Add(distTax)
			}

			// Final Price
			totalFinalPrice = totalFinalPrice.Add(finalPrice)
			resultString = finalPrice.String()
		}
	})
}
