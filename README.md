# Benchmarks

This repository contains benchmarks for the decimal arithmetic library [govalues/decimal].
Additionally, the results of [govalues/decimal] are compared to [cockroachdb/apd] and [shopspring/decimal].

## Getting started

Clone the repository:

```bash
git clone https://github.com/govalues/decimal-benchmarks.git
```

Install the necessary dependencies:

```bash
go install golang.org/x/perf/cmd/benchstat
go mod download
```

## Running Benchmarks

To run all benchmarks, simply run the following command:

```bash
go test -count=30 -timeout=60m -bench . github.com/govalues/benchmarks > results.txt
```

To print a summary of the benchmark results, including statistics such as mean,
standard deviation, and confidence intervals, execute the following command:

```bash
benchstat -col /mod results.txt
```

[govalues/decimal]: https://github.com/govalues/decimal
[shopspring/decimal]: https://github.com/govalues/decimal
[cockroachdb/apd]: https://github.com/cockroachdb/apd
