# decimal

This repository contains tests and benchmarks for the decimal arithmetic
library [govalues/decimal].

## Getting Started

Clone the repository:

```bash
git clone https://github.com/govalues/decimal-tests.git
cd decimal-tests
```

Install the necessary dependencies:

```bash
go install golang.org/x/perf/cmd/benchstat@latest
go install github.com/go-task/task/v3/cmd/task@latest
```

## Running Tests

| Command      | Description                                                              |
| ------------ | ------------------------------------------------------------------------ |
| `task fuzz`  | Check the correctness against [cockroachdb/apd] and [shopspring/decimal] |
| `task bench` | CPU and memory usage                                                     |

[govalues/decimal]: https://github.com/govalues/decimal
[shopspring/decimal]: https://github.com/shopspring/decimal
[cockroachdb/apd]: https://github.com/cockroachdb/apd
