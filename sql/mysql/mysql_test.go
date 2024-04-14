package mysql_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/govalues/decimal"
)

const (
	url           = "root:password@tcp(localhost:3306)/test"
	selectNull    = "SELECT null"
	dropTable     = "DROP TABLE IF EXISTS decimal_tests"
	createTable   = "CREATE TABLE decimal_tests (id INT AUTO_INCREMENT PRIMARY KEY, number DECIMAL(19,2))"
	insertDecimal = "INSERT INTO decimal_tests (number) VALUES (?)"
	selectDecimal = "SELECT number FROM decimal_tests WHERE id = ?"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("mysql", url)
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
		}
	})

	t.Run("decimal.NullDecimal", func(t *testing.T) {
		row := db.QueryRow(selectNull)
		got := decimal.NullDecimal{}
		err := row.Scan(&got)
		if err != nil {
			t.Errorf("QueryRow(%q) failed: %v", selectNull, err)
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
		{"0", "0.00"},
		{"0.0", "0.00"},
		{"0.00", "0.00"},
		{"0.000", "0.00"},
		{"0.000000000000000000", "0.00"},

		{"1", "1.00"},
		{"1.0", "1.00"},
		{"1.00", "1.00"},
		{"1.000", "1.00"},
		{"1.000000000000000000", "1.00"},

		{"-1", "-1.00"},
		{"-1.0", "-1.00"},
		{"-1.00", "-1.00"},
		{"-1.000", "-1.00"},
		{"-1.000000000000000000", "-1.00"},

		{"0.1", "0.10"},
		{"0.10", "0.10"},
		{"0.100", "0.10"},
		{"0.1000", "0.10"},
		{"0.1000000000000000000", "0.10"},

		{"-0.1", "-0.10"},
		{"-0.10", "-0.10"},
		{"-0.100", "-0.10"},
		{"-0.1000", "-0.10"},
		{"-0.1000000000000000000", "-0.10"},

		{"0.1", "0.10"},
		{"0.01", "0.01"},
		{"0.001", "0.00"},
		{"0.0001", "0.00"},
		{"0.0000000000000000001", "0.00"},

		{"-0.1", "-0.10"},
		{"-0.01", "-0.01"},
		{"-0.001", "-0.00"},
		{"-0.0001", "-0.00"},
		{"-0.0000000000000000001", "-0.00"},

		{"1", "1.00"},
		{"10", "10.00"},
		{"100", "100.00"},
		{"1000", "1000.00"},
		{"10000000000000000", "10000000000000000.00"},
		{"10000000000000000.00", "10000000000000000.00"},

		{"-1", "-1.00"},
		{"-10", "-10.00"},
		{"-100", "-100.00"},
		{"-1000", "-1000.00"},
		{"-10000000000000000", "-10000000000000000.00"},

		{"99999999999999999.99", "99999999999999999.99"},
		{"-99999999999999999.99", "-99999999999999999.99"},

		// Rounding

		{"0.005", "0.01"},
		{"0.015", "0.02"},
		{"0.025", "0.03"},
		{"0.035", "0.04"},

		{"-0.005", "-0.01"},
		{"-0.015", "-0.02"},
		{"-0.025", "-0.03"},
		{"-0.035", "-0.04"},

		{"9999999999999999.994", "9999999999999999.99"},
		{"9999999999999999.995", "10000000000000000.00"},
		{"9999999999999999.996", "10000000000000000.00"},

		{"-9999999999999999.994", "-9999999999999999.99"},
		{"-9999999999999999.995", "-10000000000000000.00"},
		{"-9999999999999999.996", "-10000000000000000.00"},

		{"2.718281828459045235", "2.72"},
		{"3.141592653589793238", "3.14"},
	}

	t.Run("decimal.Decimal", func(t *testing.T) {
		for _, tt := range tests {
			d, err := decimal.Parse(tt.d)
			if err != nil {
				t.Errorf("Parse(%q) failed: %v", tt.d, err)
				continue
			}
			result, err := db.Exec(insertDecimal, d)
			if err != nil {
				t.Errorf("Exec(%q, %v) failed: %v", insertDecimal, d, err)
				continue
			}
			rowID, err := result.LastInsertId()
			if err != nil {
				t.Errorf("LastInsertId() failed: %v", err)
				continue
			}
			row := db.QueryRow(selectDecimal, rowID)
			got := decimal.Decimal{}
			err = row.Scan(&got)
			if err != nil {
				t.Errorf("QueryRow(%q, %v) failed: %v", selectDecimal, rowID, err)
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
			result, err := db.Exec(insertDecimal, d)
			if err != nil {
				t.Errorf("Exec(%q, %v) failed: %v", insertDecimal, d, err)
				continue
			}
			rowID, err := result.LastInsertId()
			if err != nil {
				t.Errorf("LastInsertId() failed: %v", err)
				continue
			}
			row := db.QueryRow(selectDecimal, rowID)
			got := decimal.NullDecimal{}
			err = row.Scan(&got)
			if err != nil {
				t.Errorf("QueryRow(%q, %v) failed: %v", selectDecimal, rowID, err)
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
