package mongo_test

import (
	"context"
	"log"
	"math"
	"os"
	"testing"

	"github.com/govalues/decimal"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	uri    = "mongodb://root:password@localhost:27017"
	dbName = "test"
)

var db *mongo.Database

func TestMain(m *testing.M) {
	ctx := context.TODO()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Connect(%q) failed: %v\n", uri, err)
	}
	defer func() {
		_ = client.Disconnect(ctx)
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Ping() failed: %v\n", err)
	}

	db = client.Database(dbName)

	os.Exit(m.Run())
}

func TestDecimal_insert(t *testing.T) {
	ctx := context.TODO()
	coll := db.Collection("decimal_insert")

	tests := []struct {
		Int32       int32
		WantInt32   string
		Int64       int64
		WantInt64   string
		Float64     float64
		WantFloat64 string
		String      string
		WantString  string
		Decimal     string
		WantDecimal string
	}{
		// Largest negative values
		{
			math.MinInt32, "-2147483648",
			math.MinInt64, "-9223372036854775808",
			-9.99999999999999e+18, "-9999999999999990000",
			"-9999999999999999999", "-9999999999999999999",
			"-9999999999999999999", "-9999999999999999999",
		},
		// Smallest negative values
		{
			-1, "-1",
			-1, "-1",
			-math.SmallestNonzeroFloat64, "0.0000000000000000000",
			"-0.0000000000000000001", "-0.0000000000000000001",
			"-0.0000000000000000001", "-0.0000000000000000001",
		},
		// Zero
		{
			0, "0",
			0, "0",
			0, "0",
			"0.0000000000000000000", "0.0000000000000000000",
			"0.0000000000000000000", "0.0000000000000000000",
		},
		// Pi
		{
			3, "3",
			3, "3",
			math.Pi, "3.141592653589793",
			decimal.Pi.String(), "3.141592653589793238",
			decimal.Pi.String(), "3.141592653589793238",
		},
		// Smallest positive values
		{
			1, "1",
			1, "1",
			math.SmallestNonzeroFloat64, "0.0000000000000000000",
			"0.0000000000000000001", "0.0000000000000000001",
			"0.0000000000000000001", "0.0000000000000000001",
		},
		// Largest positive values
		{
			math.MaxInt32, "2147483647",
			math.MaxInt64, "9223372036854775807",
			9.99999999999999e+18, "9999999999999990000",
			"9999999999999999999", "9999999999999999999",
			"9999999999999999999", "9999999999999999999",
		},
	}

	for _, tt := range tests {
		// Arrange
		in := struct {
			Int32   int32           `bson:"int32"`
			Int64   int64           `bson:"int64"`
			Float64 float64         `bson:"float64"`
			String  string          `bson:"string"`
			Decimal decimal.Decimal `bson:"decimal"`
		}{
			Int32:   tt.Int32,
			Int64:   tt.Int64,
			Float64: tt.Float64,
			String:  tt.String,
			Decimal: decimal.MustParse(tt.Decimal),
		}
		result, err := coll.InsertOne(ctx, in)
		if err != nil {
			t.Errorf("InsertOne(%#v) failed: %v", in, err)
			continue
		}
		recordID := result.InsertedID

		// Act
		var got struct {
			Int32   decimal.Decimal `bson:"int32"`
			Int64   decimal.Decimal `bson:"int64"`
			Float64 decimal.Decimal `bson:"float64"`
			String  decimal.Decimal `bson:"string"`
			Decimal decimal.Decimal `bson:"decimal"`
		}
		err = coll.FindOne(ctx, bson.M{"_id": recordID}).Decode(&got)
		if err != nil {
			t.Errorf("FindOne(%v) failed: %v", recordID, err)
			continue
		}

		// Assert
		if got.Int32 != decimal.MustParse(tt.WantInt32) {
			t.Errorf("Int32 = %v, want %v", got.Int32, tt.WantInt32)
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
		if got.Decimal != decimal.MustParse(tt.WantDecimal) {
			t.Errorf("Decimal = %v, want %v", got.Decimal, tt.WantDecimal)
		}
	}
}

func TestDecimal_selectNull(t *testing.T) {
	ctx := context.TODO()
	coll := db.Collection("decimal_null")

	result, err := coll.InsertOne(ctx, bson.M{"decimal": nil})
	if err != nil {
		t.Fatalf("InsertOne failed: %v", err)
	}
	recordID := result.InsertedID

	var got struct {
		Decimal decimal.Decimal `bson:"decimal"`
	}
	err = coll.FindOne(ctx, bson.M{"_id": recordID}).Decode(&got)
	if err != nil {
		t.Errorf("Decode() failed: %v", err)
		return
	}

	want := decimal.Decimal{}
	if got.Decimal != want {
		t.Errorf("Decimal = %v, want %v", got.Decimal, want)
	}
}

func TestNullDecimal_selectNull(t *testing.T) {
	ctx := context.TODO()
	coll := db.Collection("decimal_null")

	result, err := coll.InsertOne(ctx, bson.M{"value": nil})
	if err != nil {
		t.Fatalf("InsertOne failed: %v", err)
	}
	recordID := result.InsertedID

	var got struct {
		Decimal decimal.NullDecimal `bson:"decimal"`
	}
	err = coll.FindOne(ctx, bson.M{"_id": recordID}).Decode(&got)
	if err != nil {
		t.Errorf("Decode() failed: %v", err)
		return
	}

	want := decimal.NullDecimal{}
	if got.Decimal != want {
		t.Errorf("NullDecimal = %v, want %v", got.Decimal, want)
	}
}
