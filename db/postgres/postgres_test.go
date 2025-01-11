package postgres_test

import (
	"log"
	"math"
	"os"
	"testing"

	"github.com/govalues/decimal"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const url = "postgres://root:password@localhost:5432/test"

var db *sqlx.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sqlx.Connect("pgx", url)
	if err != nil {
		log.Fatalf("Connect(%q) failed: %v\n", url, err)
	}
	defer db.Close()

	os.Exit(m.Run())
}

func TestDecimal_insert(t *testing.T) {
	const createTable = `
CREATE TABLE IF NOT EXISTS decimal_insert (
	int16 SMALLINT,
	int32 INTEGER,
	int64 BIGINT,
	float32 REAL,
	float64 DOUBLE PRECISION,
	string TEXT,
	decimal DECIMAL, 
	decimal19_2 DECIMAL(19, 2), 
	decimal19_0 DECIMAL(19, 0)
)`
	_, err := db.Exec(createTable)
	if err != nil {
		t.Fatalf("Exec(%q) failed: %v", createTable, err)
	}

	tests := []struct {
		Int16           int16
		WantInt16       string
		Int32           int32
		WantInt32       string
		Int64           int64
		WantInt64       string
		Float32         float32
		WantFloat32     string
		Float64         float64
		WantFloat64     string
		String          string
		WantString      string
		Decimal         string
		WantDecimal     string
		Decimal19_2     string
		WantDecimal19_2 string
		Decimal19_0     string
		WantDecimal19_0 string
	}{
		// Largest negative values
		{
			math.MinInt16, "-32768",
			math.MinInt32, "-2147483648",
			math.MinInt64, "-9223372036854775808",
			-9.9999999e+18, "-9999999980506448000",
			-9.99999999999999e+18, "-9999999999999990000",
			"-9999999999999999999", "-9999999999999999999",
			"-9999999999999999999", "-9999999999999999999",
			"-99999999999999999.99", "-99999999999999999.99",
			"-9999999999999999999", "-9999999999999999999",
		},
		// Smallest negative values
		{
			-1, "-1",
			-1, "-1",
			-1, "-1",
			-math.SmallestNonzeroFloat32, "0.0000000000000000000",
			-math.SmallestNonzeroFloat64, "0.0000000000000000000",
			"-0.0000000000000000001", "-0.0000000000000000001",
			"-0.0000000000000000001", "-0.0000000000000000001",
			"-0.0000000000000000001", "0.00",
			"-0.0000000000000000001", "0",
		},
		// Zero
		{
			0, "0",
			0, "0",
			0, "0",
			0, "0",
			0, "0",
			"0.0000000000000000000", "0.0000000000000000000",
			"0.0000000000000000000", "0.0000000000000000000",
			"0.0000000000000000000", "0.00",
			"0.0000000000000000000", "0",
		},
		// Pi
		{
			3, "3",
			3, "3",
			3, "3",
			math.Pi, "3.1415927410125732",
			math.Pi, "3.141592653589793",
			decimal.Pi.String(), "3.141592653589793238",
			decimal.Pi.String(), "3.141592653589793238",
			decimal.Pi.String(), "3.14",
			decimal.Pi.String(), "3",
		},
		// Smallest positive values
		{
			1, "1",
			1, "1",
			1, "1",
			math.SmallestNonzeroFloat32, "0.0000000000000000000",
			math.SmallestNonzeroFloat64, "0.0000000000000000000",
			"0.0000000000000000001", "0.0000000000000000001",
			"0.0000000000000000001", "0.0000000000000000001",
			"0.0000000000000000001", "0.00",
			"0.0000000000000000001", "0",
		},
		// Largest positive values
		{
			math.MaxInt16, "32767",
			math.MaxInt32, "2147483647",
			math.MaxInt64, "9223372036854775807",
			9.999999e+18, "9999998880994820000",
			9.99999999999999e+18, "9999999999999990000",
			"9999999999999999999", "9999999999999999999",
			"9999999999999999999", "9999999999999999999",
			"99999999999999999.99", "99999999999999999.99",
			"9999999999999999999", "9999999999999999999",
		},
	}

	for _, tt := range tests {
		var got struct {
			SmallInt    decimal.Decimal `db:"int16"`
			Int         decimal.Decimal `db:"int32"`
			BigInt      decimal.Decimal `db:"int64"`
			Real        decimal.Decimal `db:"float32"`
			Double      decimal.Decimal `db:"float64"`
			String      decimal.Decimal `db:"string"`
			Decimal     decimal.Decimal `db:"decimal"`
			Decimal19_2 decimal.Decimal `db:"decimal19_2"`
			Decimal19_0 decimal.Decimal `db:"decimal19_0"`
		}

		const insert = `
INSERT INTO decimal_insert (int16, int32, int64, float32, float64, string, decimal, decimal19_2, decimal19_0)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING int16, int32, int64, float32, float64, string, decimal, decimal19_2, decimal19_0`

		err := db.QueryRowx(
			insert,
			tt.Int16,
			tt.Int32,
			tt.Int64,
			tt.Float32,
			tt.Float64,
			tt.String,
			decimal.MustParse(tt.Decimal),
			decimal.MustParse(tt.Decimal19_2),
			decimal.MustParse(tt.Decimal19_0),
		).StructScan(&got)
		if err != nil {
			t.Fatalf("StructScan failed: %v", err)
		}

		if got.SmallInt != decimal.MustParse(tt.WantInt16) {
			t.Errorf("Int16 = %v, want %v", got.SmallInt, tt.WantInt16)
		}
		if got.Int != decimal.MustParse(tt.WantInt32) {
			t.Errorf("Int32 = %v, want %v", got.Int, tt.WantInt32)
		}
		if got.BigInt != decimal.MustParse(tt.WantInt64) {
			t.Errorf("Int64 = %v, want %v", got.BigInt, tt.WantInt64)
		}
		if got.Real != decimal.MustParse(tt.WantFloat32) {
			t.Errorf("Float32 = %v, want %v", got.Real, tt.WantFloat32)
		}
		if got.Double != decimal.MustParse(tt.WantFloat64) {
			t.Errorf("Float64 = %v, want %v", got.Double, tt.WantFloat64)
		}
		if got.String != decimal.MustParse(tt.WantString) {
			t.Errorf("String = %v, want %v", got.String, tt.WantString)
		}
		if got.Decimal != decimal.MustParse(tt.WantDecimal) {
			t.Errorf("Decimal = %v, want %v", got.Decimal, tt.WantDecimal)
		}
		if got.Decimal19_2 != decimal.MustParse(tt.WantDecimal19_2) {
			t.Errorf("Decimal19_2 = %v, want %v", got.Decimal19_2, tt.WantDecimal19_2)
		}
		if got.Decimal19_0 != decimal.MustParse(tt.WantDecimal19_0) {
			t.Errorf("Decimal19_0 = %v, want %v", got.Decimal19_0, tt.WantDecimal19_0)
		}
	}
}

func TestDecimal_selectNull(t *testing.T) {
	const selectNull = "SELECT null"
	var got decimal.Decimal
	err := db.QueryRowx(selectNull).Scan(&got)
	if err == nil {
		t.Errorf("QueryRowx(%q) did not fail, got %v", selectNull, got)
	}
}

func TestNullDecimal_selectNull(t *testing.T) {
	const selectNull = "SELECT null"
	var got decimal.NullDecimal
	err := db.QueryRowx(selectNull).Scan(&got)
	if err != nil {
		t.Errorf("QueryRowx(%q) failed: %v", selectNull, err)
		return
	}
	want := decimal.NullDecimal{}
	if got != want {
		t.Errorf("Scan() = %v, want %v", got, want)
	}
}
