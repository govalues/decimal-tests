package sqlite_test

import (
	"log"
	"math"
	"os"
	"testing"

	"github.com/govalues/decimal"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const url = ":memory:"

var db *sqlx.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sqlx.Connect("sqlite", url)
	if err != nil {
		log.Fatalf("Connect(%q) failed: %v\n", url, err)
	}
	defer db.Close()

	os.Exit(m.Run())
}

func TestDecimal_insert(t *testing.T) {
	const createTable = `
CREATE TABLE IF NOT EXISTS decimal_insert (
	int64   INTEGER,
	float64 REAL,
	string  TEXT
)`
	_, err := db.Exec(createTable)
	if err != nil {
		t.Fatalf("Exec(%q) failed: %v", createTable, err)
	}

	tests := []struct {
		Int64       int64
		WantInt64   string
		Float64     float64
		WantFloat64 string
		String      string
		WantString  string
	}{
		// Largest negative values
		{
			math.MinInt64, "-9223372036854775808",
			-9.99999999999999e+18, "-9999999999999990000",
			"-9999999999999999999", "-9999999999999999999",
		},
		// Smallest negative values
		{
			-1, "-1",
			-math.SmallestNonzeroFloat64, "0.0000000000000000000",
			"-0.0000000000000000001", "-0.0000000000000000001",
		},
		// Zero
		{
			0, "0",
			0, "0",
			"0.0000000000000000000", "0.0000000000000000000",
		},
		// Pi
		{
			3, "3",
			math.Pi, "3.141592653589793",
			decimal.Pi.String(), "3.141592653589793238",
		},
		// Smallest positive values
		{
			1, "1",
			math.SmallestNonzeroFloat64, "0.0000000000000000000",
			"0.0000000000000000001", "0.0000000000000000001",
		},
		// Largest positive values
		{
			math.MaxInt64, "9223372036854775807",
			9.99999999999999e+18, "9999999999999990000",
			"9999999999999999999", "9999999999999999999",
		},
	}

	for _, tt := range tests {
		var got struct {
			Int64   decimal.Decimal `db:"int64"`
			Float64 decimal.Decimal `db:"float64"`
			String  decimal.Decimal `db:"string"`
		}

		const insert = `
INSERT INTO decimal_insert (int64, float64, string)
VALUES ($1, $2, $3)
RETURNING int64, float64, string`

		err := db.QueryRowx(
			insert,
			tt.Int64,
			tt.Float64,
			tt.String,
		).StructScan(&got)
		if err != nil {
			t.Fatalf("StructScan failed: %v", err)
		}

		if got.Int64 != decimal.MustParse(tt.WantInt64) {
			t.Errorf("Int64 = %v, want %v", got.Int64, tt.WantInt64)
		}
		if got.Float64 != decimal.MustParse(tt.WantFloat64) {
			t.Errorf("Float64 = %v, want %v", got.Float64, tt.WantFloat64)
		}
		if got.String != decimal.MustParse(tt.WantString) {
			t.Errorf("String = %v, want %v", got.String, tt.WantString)
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
