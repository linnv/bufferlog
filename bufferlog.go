package bufferlog

import (
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type BufLog struct {
	buf         []byte
	mux         sync.RWMutex
	exit        chan struct{}
	underlyFile io.WriteCloser

	Len           int           `json:"Len"`
	FlushInterval time.Duration `json:"FlushInterval"`
}

//NewBufferLog implements return bufferlog filled with size, flush ticket and underly file
func NewBufferLog(bufferSize int, flushInterval time.Duration, exit chan struct{}, w io.WriteCloser) *BufLog {
	one := newBufferLog(bufferSize, flushInterval, w)
	one.exit = exit
	go func() {
		if err := one.flushIntervally(); err != nil {
			print(err.Error())
		}
	}()
	return one
}

func newBufferLog(bufferSize int, flushInterval time.Duration, w io.WriteCloser) *BufLog {
	if bufferSize < 1024 {
		bufferSize = 1024
	}
	one := &BufLog{
		Len:           bufferSize,
		FlushInterval: flushInterval,
		underlyFile:   w,
	}
	makeSlice := func(n int) []byte {
		defer func() {
			if err := recover(); err != nil {
				panic(err)
			}
		}()
		return make([]byte, 0, n)
	}

	one.buf = makeSlice(one.Len)
	return one
}

const maxSize = 1 << 23 //8M
func (b *BufLog) Write(bs []byte) (n int, err error) {
	if b == nil {
		return 0, ERR_EMPTY_REFENCE
	}
	b.mux.Lock()
	//@TODO remove defer
	defer b.mux.Unlock()

	if len(bs) > maxSize {
		return b.underlyFile.Write(bs)
	}

	if len(bs)+len(b.buf) > b.Len {
		if n, err = b.flush(); err != nil {
			err = errors.Wrap(err, "Write")
			return
		}
	}

	b.buf = append(b.buf, bs...)
	if len(bs) > b.Len {
		log.Printf("resize from %d to %d", b.Len, cap(bs))
		b.Len = cap(bs)
	}
	return len(bs), nil
}

var ERR_EMPTY_REFENCE = errors.New("empty pointer")

func (b *BufLog) Close() (err error) {
	if b == nil {
		return ERR_EMPTY_REFENCE
	}
	b.mux.Lock()
	defer b.mux.Unlock()
	if _, err = b.flush(); err != nil {
		return errors.Wrap(err, "flush")
	}
	if b.underlyFile == os.Stdout {
		return
	}
	if err = b.underlyFile.Close(); err != nil {
		return errors.Wrap(err, "underlyFile.Close")
	}
	return
}

func (b *BufLog) Flush() (err error) {
	if b == nil {
		return ERR_EMPTY_REFENCE
	}
	b.mux.Lock()
	defer b.mux.Unlock()
	if _, err = b.flush(); err != nil {
		return errors.Wrap(err, "Flush")
	}
	return
}

func (b *BufLog) flush() (n int, err error) {
	if b == nil {
		return 0, ERR_EMPTY_REFENCE
	}
	if len(b.buf) > 0 {
		n, err = b.underlyFile.Write(b.buf[:len(b.buf)])
		if err != nil {
			err = errors.Wrap(err, "flush")
			return
		}
		b.buf = b.buf[:0]
	}
	return
}

func (b *BufLog) flushIntervally() (err error) {
	ticker := time.NewTicker(b.FlushInterval)
	log.Printf("BufLog flushes buffer every %v\n", b.FlushInterval)
	for {
		select {
		case <-b.exit:
			log.Println("exit Buflog")
			if err = b.Close(); err != nil {
				return errors.Wrap(err, "flushIntervally on Close")
			}
			return
		case <-ticker.C:
			if err = b.Flush(); err != nil {
				return errors.Wrap(err, "flushIntervally on Flush")
			}
		}
	}
}
