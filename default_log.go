package bufferlog

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//Buffer is just for debug, you'd better new Buflog by yourself to control Flush on exiting by the exit channel
var Buffer *BufLog

func init() {
	exit := make(chan struct{})
	flushInterval := time.Millisecond * 100
	Buffer = newBufferLog(1<<10, flushInterval, os.Stdout)
	Buffer.exit = exit
	go Buffer.flushIntervally()

	go func() {
		sigChan := make(chan os.Signal, 2)
		signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-sigChan
		close(exit)
		time.Sleep(flushInterval) //make sure Buffer has exited, or invoke Close() directly
		log.Printf("Buffer: receive exit signal \n")
	}()
}

func BufferDemo() {
	Buffer.Write([]byte("abcd\n"))
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}()

	<-sigChan
	log.Printf("BufferDemo: receive exit signal \n")
	time.Sleep(time.Second * 2) //make sure Buffer has exited, or invoke Close() directly
}
