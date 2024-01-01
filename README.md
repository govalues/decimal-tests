# Benchmarks

This repository contains tests and benchmarks for the decimal arithmetic
library [govalues/decimal].

## Getting Started

Clone the repository:

```bash
git clone https://github.com/govalues/decimal-tests.git
cd decimal-tests
go mod download
```

## Running Benchmarks

Install the necessary dependencies:

```bash
go install golang.org/x/perf/cmd/benchstat
```

To measure CPU usage, run the following command:

```bash
go test -count=30 -timeout=120m -bench . github.com/govalues/decimal-tests > benchstat.txt
benchstat -filter ".unit:ns/op" -col /mod benchstat.txt
```

To measure RAM usage, run the following command:

```bash
go test -count=6 -timeout=30m -benchmem -bench . github.com/govalues/decimal-tests > benchstat.txt
benchstat -filter ".unit:B/op" -col /mod benchstat.txt
```

## Running Fuzz Tests

To compare the correctness of the decimal arithmetic library [govalues/decimal]
to [cockroachdb/apd] and [shopspring/decimal], run the following commands:

```bash
go test -fuzztime 20s -fuzz ^FuzzDecimalAdd$ github.com/govalues/decimal-tests
go test -fuzztime 20s -fuzz ^FuzzDecimalMul$ github.com/govalues/decimal-tests
go test -fuzztime 20s -fuzz ^FuzzDecimalQuo$ github.com/govalues/decimal-tests
go test -fuzztime 40s -fuzz ^FuzzDecimalPow$ github.com/govalues/decimal-tests
```

## Running Mutation Tests

Clone the `decimal` repository next to the `decimal-tests` repository:

```bash
git clone https://github.com/govalues/decimal.git
```

To execute mutation tests, run the following command:

```bash
go test -tags=mutation -timeout=120m > ooze.txt
```

[govalues/decimal]: https://github.com/govalues/decimal
[shopspring/decimal]: https://github.com/shopspring/decimal
[cockroachdb/apd]: https://github.com/cockroachdb/apd
