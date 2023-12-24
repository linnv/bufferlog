// Package main provides ...
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/linnv/bufferlog"
)

type DiscardCloser struct {
	io.Writer
}

func (DiscardCloser) Close() error { return nil }
func main() {
	fmt.Println()
	exit := make(chan struct{})
	var Discard io.WriteCloser = DiscardCloser{ioutil.Discard}
	logfdBuf := bufferlog.NewBufferLog(3*1024, time.Millisecond*100, exit, Discard)
	size := 100
	fmt.Printf("size: %dMB\n", size)
	forsize := (1 << 20) * size
	bs := make([]byte, forsize)
	for i := 0; i < int(forsize); i++ {
		bs[i] = '1'
	}
	incMem := func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		for {
			select {
			case <-ticker.C:
				oneSize := 1 << 20 * (rand.Int63n(99) + 1)
				logfdBuf.Write(bs[:oneSize])
			}
		}
	}
	for i := 0; i < 100; i++ {
		go incMem()
	}
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	log.Print("use c-c to exit: \n")
	gotSignal := <-sigChan
	log.Printf("test receive sginal %v \n", gotSignal)
	os.Exit(0)

}
