package bufferlog

import "io"

type BufferLogger interface {
	io.WriteCloser
	Flush() error
}
