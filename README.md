# Bufferlog

<a href="https://circleci.com/gh/linnv/bufferlog">
<img src="https://circleci.com/gh/linnv/bufferlog.svg?style=shield" alt="circleci">
</a>

👾 Bufferlog is a lib for improving log persistence efficiency, actually it's just an io.WriteCloser with buffer-enabled.

### Example

```
sigChan := make(chan os.Signal, 2)
exit := make(chan struct{})
fileBuffer := "./demoBuffer.log"
under := &lumberjack.Logger{
	Filename:   fileBuffer,
	MaxSize:    100, // megabytes
	MaxBackups: 3,
	LocalTime:  true,
	MaxAge:     28, // days
}
logger := NewBufferLog(3*1024, time.Second*2, exit, under)
logger.Write([]byte("abc\n"))
signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGSTOP)
log.Print("use c-c to exit: \n")
<-sigChan
close(exit)
time.Sleep(time.Second * 3) //make sure logger has exited, or invoke Close() directly
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
