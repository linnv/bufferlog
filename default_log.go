package bufferlog

import (
	"log"
	"os"
	"os/signal"
	"strconv"
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
	go func() {
		if err := Buffer.flushIntervally(); err != nil {
			print(err.Error())
		}
	}()

	go func() {
		sigChan := make(chan os.Signal, 2)
		signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-sigChan
		close(exit)
		time.Sleep(flushInterval) //make sure Buffer has exited, or invoke Close() directly
		Buffer.Flush()
		log.Printf("Buffer: receive exit signal \n")
	}()
}

func BufferDemo() {
	bsWrite := []byte("abcd\n")
	if n, err := Buffer.Write(bsWrite); err != nil {
		errStr := "wrote " + strconv.Itoa(n) + " bytes want " + strconv.Itoa(len(bsWrite)) + " bytes, err:" + err.Error()
		print(errStr)
	}

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		if err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM); err != nil {
			print(err.Error())
		}
	}()

	<-sigChan
	time.Sleep(time.Second * 2) //make sure Buffer has exited, or invoke Close() directly
}
