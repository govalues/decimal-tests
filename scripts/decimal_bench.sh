cd ..
go test -run=^$ -count=10 -bench . github.com/govalues/benchmarks > bench.txt
benchstat -col /mod bench.txt > benchstat.txt