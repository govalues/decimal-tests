# Benchmarks

This repository contains benchmarks for the decimal arithmetic library [govalues/decimal].
Additionally, the results of [govalues/decimal] are compared to [cockroachdb/apd] and [shopspring/decimal].

## Getting started

Clone the repository:

```bash
git clone https://github.com/govalues/decimal-tests.git
```

Install the necessary dependencies:

```bash
go install golang.org/x/perf/cmd/benchstat
```

## Running Benchmarks

To measure CPU usage, run the following command:

```bash
go test -count=30 -timeout=120m -bench . github.com/govalues/decimal-tests > results.txt
benchstat -filter ".unit:ns/op" -col /mod results.txt
```

To measure RAM usage, run the following command:

```bash
go test -count=6 -timeout=30m -benchmem -bench . github.com/govalues/decimal-tests > results.txt
benchstat -filter ".unit:B/op" -col /mod results.txt
```

[govalues/decimal]: https://github.com/govalues/decimal
[shopspring/decimal]: https://github.com/shopspring/decimal
[cockroachdb/apd]: https://github.com/cockroachdb/apd
