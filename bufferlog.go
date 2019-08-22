package bufferlog

import (
	"io"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type BufLog struct {
	buf         []byte `json:"buf"`
	mux         sync.RWMutex
	exit        chan struct{}
	underlyFile io.WriteCloser `json:"underlyFile"`

	Len           int           `json:"Len"`
	FlushInterval time.Duration `json:"FlushInterval"`
}

//NewBufferLog implements return bufferlog filled with size, flush ticket and underly file
func NewBufferLog(bufferSize int, flushInterval time.Duration, exit chan struct{}, w io.WriteCloser) *BufLog {
	one := newBufferLog(bufferSize, flushInterval, w)
	one.exit = exit
	go one.flushIntervally()
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

func (b *BufLog) Write(bs []byte) (err error) {
	if b == nil {
		return ERR_EMPTY_REFENCE
	}
	b.mux.Lock()
	defer b.mux.Unlock()
	if len(bs)+len(b.buf) > b.Len {
		if err := b.flush(); err != nil {
			return errors.Wrap(err, "Write")
		}
	}
	b.buf = append(b.buf, bs...)
	return
}

var ERR_EMPTY_REFENCE = errors.New("empty pointer")

func (b *BufLog) Close() (err error) {
	if b == nil {
		return ERR_EMPTY_REFENCE
	}
	b.mux.Lock()
	defer b.mux.Unlock()
	if err = b.flush(); err != nil {
		return errors.Wrap(err, "flush")
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
	if err = b.flush(); err != nil {
		return errors.Wrap(err, "Flush")
	}
	return
}

func (b *BufLog) flush() (err error) {
	if b == nil {
		return ERR_EMPTY_REFENCE
	}
	if len(b.buf) > 0 {
		_, err = b.underlyFile.Write(b.buf[:len(b.buf)])
		if err != nil {
			return errors.Wrap(err, "flush")
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
