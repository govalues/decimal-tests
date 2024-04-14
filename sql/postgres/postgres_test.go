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
	createTable   = "CREATE TABLE decimal_tests (number DECIMAL)"
	insertDecimal = "INSERT INTO decimal_tests (number) VALUES ($1) RETURNING number"
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
		got := decimal.Decimal{}
		row := db.QueryRow(selectNull)
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
		d, want string
	}{
		{"0", "0"},
		{"0.0", "0.0"},
		{"0.00", "0.00"},
		{"0.000", "0.000"},
		{"0.000000000000000000", "0.000000000000000000"},

		{"1", "1"},
		{"1.0", "1.0"},
		{"1.00", "1.00"},
		{"1.000", "1.000"},
		{"1.000000000000000000", "1.000000000000000000"},

		{"-1", "-1"},
		{"-1.0", "-1.0"},
		{"-1.00", "-1.00"},
		{"-1.000", "-1.000"},
		{"-1.000000000000000000", "-1.000000000000000000"},

		{"0.1", "0.1"},
		{"0.10", "0.10"},
		{"0.100", "0.100"},
		{"0.1000", "0.1000"},
		{"0.1000000000000000000", "0.1000000000000000000"},

		{"-0.1", "-0.1"},
		{"-0.10", "-0.10"},
		{"-0.100", "-0.100"},
		{"-0.1000", "-0.1000"},
		{"-0.1000000000000000000", "-0.1000000000000000000"},

		{"0.1", "0.1"},
		{"0.01", "0.01"},
		{"0.001", "0.001"},
		{"0.0001", "0.0001"},
		{"0.0000000000000000001", "0.0000000000000000001"},

		{"-0.1", "-0.1"},
		{"-0.01", "-0.01"},
		{"-0.001", "-0.001"},
		{"-0.0001", "-0.0001"},
		{"-0.0000000000000000001", "-0.0000000000000000001"},

		{"1", "1"},
		{"10", "10"},
		{"100", "100"},
		{"1000", "1000"},
		{"1000000000000000000", "1000000000000000000"},

		{"-1", "-1"},
		{"-10", "-10"},
		{"-100", "-100"},
		{"-1000", "-1000"},
		{"-1000000000000000000", "-1000000000000000000"},

		{"9999999999999999999", "9999999999999999999"},
		{"-9999999999999999999", "-9999999999999999999"},

		{"2.718281828459045235", "2.718281828459045235"},
		{"3.141592653589793238", "3.141592653589793238"},
	}

	t.Run("decimal.Decimal", func(t *testing.T) {
		for _, tt := range tests {
			d, err := decimal.Parse(tt.d)
			if err != nil {
				t.Errorf("Parse(%q) failed: %v", tt.d, err)
				continue
			}
			row := db.QueryRow(insertDecimal, d)
			got := decimal.Decimal{}
			err = row.Scan(&got)
			if err != nil {
				t.Errorf("QueryRow(%q, %v) failed: %v", insertDecimal, d, err)
				continue
			}
			want, err := decimal.Parse(tt.want)
			if err != nil {
				t.Errorf("Parse(%q) failed: %v", tt.want, err)
				continue
			}
			if got != want {
				t.Errorf("Scan() = %v, want %v", got, want)
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
			row := db.QueryRow(insertDecimal, d)
			got := decimal.NullDecimal{}
			err = row.Scan(&got)
			if err != nil {
				t.Errorf("QueryRow(%q, %v) failed: %v", insertDecimal, d, err)
				continue
			}
			want := decimal.NullDecimal{}
			err = want.Scan(tt.want)
			if err != nil {
				t.Errorf("Scan(%q) failed: %v", tt.want, err)
				continue
			}
			if got != want {
				t.Errorf("Scan() = %v, want %v", got, want)
				continue
			}
		}
	})
}