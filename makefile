.PHONY: all test clean bench pprof

all: test bench

test:
	go test
clean:
	rm demo*
bench:
	# go  test -bench=^BenchmarkBufferLog$ -count 5 -benchmem
	go test -bench Benchmark -run xx -count 5 -benchmem
pprof:
	go test -bench Benchmark -run xx -count 5 -benchmem -memprofile memprofile.out -cpuprofile profile.out
