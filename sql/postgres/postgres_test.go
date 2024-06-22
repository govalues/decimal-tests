package postgres_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/govalues/decimal"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	url           = "postgres://root:password@localhost:5432/test"
	selectNull    = "SELECT null"
	dropTable     = "DROP TABLE IF EXISTS decimal_tests"
	createTable   = "CREATE TABLE decimal_tests (num1 DECIMAL, num2 DECIMAL(19, 2), num3 DECIMAL(19, 0))"
	insertDecimal = "INSERT INTO decimal_tests (num1, num2, num3) VALUES ($1, $2, $3) RETURNING num1, num2, num3"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("pgx", url)
	if err != nil {
		log.Fatalf("Open(%q) failed: %v\n", url, err)
	}
	defer db.Close()
	_, err = db.Exec(dropTable)
	if err != nil {
		log.Fatalf("Exec(%q) failed: %v\n", dropTable, err)
	}
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("Exec(%q) failed: %v\n", createTable, err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestDecimal_selectNull(t *testing.T) {
	t.Run("decimal.Decimal", func(t *testing.T) {
		row := db.QueryRow(selectNull)
		got := decimal.Decimal{}
		err := row.Scan(&got)
		if err == nil {
			t.Errorf("QueryRow(%q) did not fail", selectNull)
			return
		}
	})

	t.Run("decimal.NullDecimal", func(t *testing.T) {
		got := decimal.NullDecimal{}
		row := db.QueryRow(selectNull)
		err := row.Scan(&got)
		if err != nil {
			t.Errorf("QueryRow(%q) failed: %v", selectNull, err)
			return
		}
		want := decimal.NullDecimal{}
		if got != want {
			t.Errorf("Scan() = %v, want %v", got, want)
		}
	})
}

func TestDecimal_insert(t *testing.T) {
	tests := []struct {
		d, want1, want2, want3 string
	}{
		{"0", "0", "0.00", "0"},
		{"0.0", "0.0", "0.00", "0"},
		{"0.00", "0.00", "0.00", "0"},
		{"0.000", "0.000", "0.00", "0"},
		{"0.000000000000000000", "0.000000000000000000", "0.00", "0"},

		{"1", "1", "1.00", "1"},
		{"1.0", "1.0", "1.00", "1"},
		{"1.00", "1.00", "1.00", "1"},
		{"1.000", "1.000", "1.00", "1"},
		{"1.000000000000000000", "1.000000000000000000", "1.00", "1"},

		{"-1", "-1", "-1.00", "-1"},
		{"-1.0", "-1.0", "-1.00", "-1"},
		{"-1.00", "-1.00", "-1.00", "-1"},
		{"-1.000", "-1.000", "-1.00", "-1"},
		{"-1.000000000000000000", "-1.000000000000000000", "-1.00", "-1"},

		{"0.1", "0.1", "0.10", "0"},
		{"0.10", "0.10", "0.10", "0"},
		{"0.100", "0.100", "0.10", "0"},
		{"0.1000", "0.1000", "0.10", "0"},
		{"0.1000000000000000000", "0.1000000000000000000", "0.10", "0"},

		{"-0.1", "-0.1", "-0.10", "0"},
		{"-0.10", "-0.10", "-0.10", "0"},
		{"-0.100", "-0.100", "-0.10", "0"},
		{"-0.1000", "-0.1000", "-0.10", "0"},
		{"-0.1000000000000000000", "-0.1000000000000000000", "-0.10", "0"},

		{"0.1", "0.1", "0.10", "0"},
		{"0.01", "0.01", "0.01", "0"},
		{"0.001", "0.001", "0.00", "0"},
		{"0.0001", "0.0001", "0.00", "0"},
		{"0.0000000000000000001", "0.0000000000000000001", "0.00", "0"},

		{"-0.1", "-0.1", "-0.10", "0"},
		{"-0.01", "-0.01", "-0.01", "0"},
		{"-0.001", "-0.001", "0.00", "0"},
		{"-0.0001", "-0.0001", "0.00", "0"},
		{"-0.0000000000000000001", "-0.0000000000000000001", "0.00", "0"},

		{"1", "1", "1.00", "1"},
		{"10", "10", "10.00", "10"},
		{"100", "100", "100.00", "100"},
		{"1000", "1000", "1000.00", "1000"},
		{"10000000000000000", "10000000000000000", "10000000000000000.00", "10000000000000000"},

		{"-1", "-1", "-1.00", "-1"},
		{"-10", "-10", "-10.00", "-10"},
		{"-100", "-100", "-100.00", "-100"},
		{"-1000", "-1000", "-1000.00", "-1000"},
		{"-10000000000000000", "-10000000000000000", "-10000000000000000.00", "-10000000000000000"},

		{"0.005", "0.005", "0.01", "0"},
		{"0.015", "0.015", "0.02", "0"},
		{"0.025", "0.025", "0.03", "0"},
		{"0.035", "0.035", "0.04", "0"},

		{"-0.005", "-0.005", "-0.01", "0"},
		{"-0.015", "-0.015", "-0.02", "0"},
		{"-0.025", "-0.025", "-0.03", "0"},
		{"-0.035", "-0.035", "-0.04", "0"},

		{"9999999999999999.994", "9999999999999999.994", "9999999999999999.99", "10000000000000000"},
		{"9999999999999999.995", "9999999999999999.995", "10000000000000000.00", "10000000000000000"},
		{"9999999999999999.996", "9999999999999999.996", "10000000000000000.00", "10000000000000000"},

		{"-9999999999999999.994", "-9999999999999999.994", "-9999999999999999.99", "-10000000000000000"},
		{"-9999999999999999.995", "-9999999999999999.995", "-10000000000000000.00", "-10000000000000000"},
		{"-9999999999999999.996", "-9999999999999999.996", "-10000000000000000.00", "-10000000000000000"},

		{"2.718281828459045235", "2.718281828459045235", "2.72", "3"},
		{"3.141592653589793238", "3.141592653589793238", "3.14", "3"},
	}

	t.Run("decimal.Decimal", func(t *testing.T) {
		for _, tt := range tests {
			d, err := decimal.Parse(tt.d)
			if err != nil {
				t.Errorf("Parse(%q) failed: %v", tt.d, err)
				continue
			}

			row := db.QueryRow(insertDecimal, d, d, d)
			var got1, got2, got3 decimal.Decimal
			err = row.Scan(&got1, &got2, &got3)
			if err != nil {
				t.Errorf("QueryRow(%q, %v) failed: %v", insertDecimal, d, err)
				continue
			}

			want1, err := decimal.Parse(tt.want1)
			if err != nil {
				t.Errorf("Parse(%q) failed: %v", tt.want1, err)
				continue
			}
			if got1 != want1 {
				t.Errorf("Scan(&got1) = %v, want %v", got1, want1)
				continue
			}

			want2, err := decimal.Parse(tt.want2)
			if err != nil {
				t.Errorf("Parse(%q) failed: %v", tt.want2, err)
				continue
			}
			if got2 != want2 {
				t.Errorf("Scan(&got2) = %v, want %v", got2, want2)
				continue
			}

			want3, err := decimal.Parse(tt.want3)
			if err != nil {
				t.Errorf("Parse(%q) failed: %v", tt.want3, err)
				continue
			}
			if got3 != want3 {
				t.Errorf("Scan(&got3) = %v, want %v", got3, want3)
				continue
			}
		}
	})

	t.Run("decimal.NullDecimal", func(t *testing.T) {
		for _, tt := range tests {
			d := decimal.NullDecimal{}
			err := d.Scan(tt.d)
			if err != nil {
				t.Errorf("Scan(%v) failed: %v", tt.d, err)
				continue
			}

			row := db.QueryRow(insertDecimal, d, d, d)
			var got1, got2, got3 decimal.NullDecimal
			err = row.Scan(&got1, &got2, &got3)
			if err != nil {
				t.Errorf("QueryRow(%q, %v) failed: %v", insertDecimal, d, err)
				continue
			}

			want1 := decimal.NullDecimal{}
			err = want1.Scan(tt.want1)
			if err != nil {
				t.Errorf("Scan(%q) failed: %v", tt.want1, err)
				continue
			}
			if got1 != want1 {
				t.Errorf("Scan(&got1) = %v, want %v", got1, want1)
				continue
			}

			want2 := decimal.NullDecimal{}
			err = want2.Scan(tt.want2)
			if err != nil {
				t.Errorf("Scan(%q) failed: %v", tt.want2, err)
				continue
			}
			if got2 != want2 {
				t.Errorf("Scan(&got2) = %v, want %v", got2, want2)
				continue
			}

			want3 := decimal.NullDecimal{}
			err = want3.Scan(tt.want3)
			if err != nil {
				t.Errorf("Scan(%q) failed: %v", tt.want3, err)
				continue
			}
			if got3 != want3 {
				t.Errorf("Scan(&got3) = %v, want %v", got3, want3)
				continue
			}
		}
	})
}
