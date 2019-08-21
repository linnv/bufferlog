<a href="https://circleci.com/gh/linnv/bufferlog">
<img src="https://circleci.com/gh/linnv/bufferlog.svg?style=shield" alt="circleci">
</a>

👾 Bufferlog is a lib for improving log persistence efficiency

### Example

```
fileBuffer := "./demotestRaw.log"
under := &lumberjack.Logger{
	Filename:   fileRaw,
	MaxSize:    100, // megabytes
	MaxBackups: 3,
	LocalTime:  true,
	MaxAge:     28, // days
}
logger := NewBufferLog(3*1024, time.Second*10, under)
logger.Write([]byte("abc\n"))
```

### Performace

Bellow shows the benchmark result of writing log to file directly and writing by bufferlog

```
go test
writeCount 1000000 bufferDemo costs  311 millisecons actually 311.168419ms
writeCount 1000000 rawDemo costs  7521 millisecons actually 7.521849673s
PASS
ok      github.com/linnv/bufferlog      7.840s
# go  test -bench=^BenchmarkBufferLog-count 5 -benchmem
go test -bench Benchmark -run xx -count 5 -benchmem
goos: darwin
goarch: amd64
pkg: github.com/linnv/bufferlog
BenchmarkBufferLog/rawWriter-4               500           3625811 ns/op             151 B/op          1 allocs/op
BenchmarkBufferLog/rawWriter-4               200           8577105 ns/op             367 B/op          2 allocs/op
BenchmarkBufferLog/rawWriter-4               500           3560456 ns/op             155 B/op          1 allocs/op
BenchmarkBufferLog/rawWriter-4               500           3012957 ns/op             205 B/op          1 allocs/op
BenchmarkBufferLog/rawWriter-4               300           9162163 ns/op             252 B/op          2 allocs/op
BenchmarkBufferLog/bufferWriter-4          30000             45015 ns/op               3 B/op          0 allocs/op
BenchmarkBufferLog/bufferWriter-4          30000             45081 ns/op               3 B/op          0 allocs/op
BenchmarkBufferLog/bufferWriter-4          30000             71305 ns/op               3 B/op          0 allocs/op
BenchmarkBufferLog/bufferWriter-4          30000             62619 ns/op               3 B/op          0 allocs/op
BenchmarkBufferLog/bufferWriter-4          30000             46142 ns/op               3 B/op          0 allocs/op
PASS
ok      github.com/linnv/bufferlog      24.353s
```
