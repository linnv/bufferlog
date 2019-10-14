package bufferlog

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func TestNewBufferLog(t *testing.T) {
	exit := make(chan struct{})
	fileBuffer := "./demoBuffer.log"
	under := &lumberjack.Logger{
		Filename:   fileBuffer,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		LocalTime:  true,
		MaxAge:     28, // days
	}
	logger := NewBufferLog(3*1024, time.Second*1, exit, under)
	bsWrite := []byte("abcd\n")
	if n, err := logger.Write(bsWrite); err != nil {
		errStr := "wrote " + strconv.Itoa(n) + " bytes want " + strconv.Itoa(len(bsWrite)) + " bytes, err:" + err.Error()
		print(errStr)
	}
	close(exit)
	time.Sleep(time.Second * 2)
}

func BenchmarkBufferLog(b *testing.B) {
	b.Run("rawWriter", func(b *testing.B) {
		filename := "./demobenraw.log"
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			under := &lumberjack.Logger{
				Filename:   filename,
				MaxSize:    100, // megabytes
				MaxBackups: 3,
				LocalTime:  true,
				MaxAge:     28, // days
			}
			for pb.Next() {
				for i := 0; i < 1024; i++ {
					bsWrite := []byte("abcd\n")
					if n, err := under.Write(bsWrite); err != nil {
						errStr := "wrote " + strconv.Itoa(n) + " bytes want " + strconv.Itoa(len(bsWrite)) + " bytes, err:" + err.Error()
						print(errStr)
					}
				}
			}
			under.Close()
		})
		if err := os.Remove(filename); err != nil {
			b.Fatal(err)
		}
	})
	b.Run("bufferWriter", func(b *testing.B) {
		filename := "./demobufferben.log"
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			under := &lumberjack.Logger{
				Filename:   filename,
				MaxSize:    100, // megabytes
				MaxBackups: 3,
				LocalTime:  true,
				MaxAge:     28, // days
			}
			logger := newBufferLog(3*1024, time.Second*10, under)
			for pb.Next() {
				for i := 0; i < 1024; i++ {
					bsWrite := []byte("abcd\n")
					if n, err := logger.Write(bsWrite); err != nil {
						errStr := "wrote " + strconv.Itoa(n) + " bytes want " + strconv.Itoa(len(bsWrite)) + " bytes, err:" + err.Error()
						print(errStr)
					}
				}
			}
			logger.Flush()
			under.Close()
		})
		if err := os.Remove(filename); err != nil {
			b.Fatal(err)
		}
	})

}

func Test_newBufferLog(t *testing.T) {
	fileRaw := "./demotest.log"
	fileBuffer := "./demotestRaw.log"
	under := &lumberjack.Logger{
		Filename:   fileRaw,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		LocalTime:  true,
		MaxAge:     28, // days
	}
	logger := newBufferLog(3*1024, time.Second*10, under)
	underRaw := &lumberjack.Logger{
		Filename:   fileBuffer,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		LocalTime:  true,
		MaxAge:     28, // days
	}

	routineDo := func(count int, f func()) {
		for i := 0; i < count; i++ {
			f()
		}
	}

	var wg sync.WaitGroup
	const jobsPerRoutine = 1000
	const writeCount = jobsPerRoutine * 300

	start := time.Now()
	for i := 0; i < writeCount; i++ {
		if i%jobsPerRoutine == 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				routineDo(jobsPerRoutine, func() {
					bsWrite := []byte("abcd\n")
					if n, err := logger.Write(bsWrite); err != nil {
						errStr := "wrote " + strconv.Itoa(n) + " bytes want " + strconv.Itoa(len(bsWrite)) + " bytes, err:" + err.Error()
						print(errStr)
					}
				})
			}()
		}
	}
	wg.Wait()
	logger.Flush()
	fmt.Printf("writeCount %d bufferDemo costs  %d millisecons actually %v\n", writeCount, time.Since(start).Nanoseconds()/1000000, time.Since(start))

	start = time.Now()
	for i := 0; i < writeCount; i++ {
		if i%jobsPerRoutine == 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				routineDo(jobsPerRoutine, func() {
					bsWrite := []byte("abcd\n")
					if n, err := underRaw.Write(bsWrite); err != nil {
						errStr := "wrote " + strconv.Itoa(n) + " bytes want " + strconv.Itoa(len(bsWrite)) + " bytes, err:" + err.Error()
						print(errStr)
					}
				})
			}()
		}
	}
	wg.Wait()
	fmt.Printf("writeCount %d rawDemo costs  %d millisecons actually %v\n", writeCount, time.Since(start).Nanoseconds()/1000000, time.Since(start))

	under.Close()
	underRaw.Close()
	fileinfoRaw, err := os.Stat(fileRaw)
	if err != nil {
		t.Fatal(err)
	}
	fileinfoBuffer, err := os.Stat(fileBuffer)
	if err != nil {
		t.Fatal(err)
	}
	if fileinfoRaw.Size() != fileinfoBuffer.Size() {
		t.Fatalf("file size not equal:raw [%s] and buffer [%s]", fileRaw, fileBuffer)
	}

}
