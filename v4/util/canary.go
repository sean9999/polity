package util

import (
	"context"
	"io"
	"strings"
)

// Canary captures output
type Canary struct {
	ch     chan []byte
	closed bool
}

// Canary implements io.ReadWriteCloser
var _ io.ReadWriteCloser = (*Canary)(nil)

func NewCanary(ctx context.Context, size int) *Canary {
	if size < 1 {
		size = 65_536
	}
	ch := make(chan []byte, size)
	c := new(Canary{ch: ch})
	go func() {
		<-ctx.Done()
		_ = c.Close()
	}()
	return c
}

func (c *Canary) Close() error {
	if c.closed == true {
		return io.ErrClosedPipe
	}
	close(c.ch)
	c.closed = true
	return nil
}

func (c *Canary) Next() ([]byte, error) {
	data, ok := <-c.ch
	if !ok {
		return nil, io.EOF
	}
	return data, nil
}

func (c *Canary) WatchFor(_ context.Context, str string) bool {
	i := 0
	for data := range c.ch {
		i++
		src := string(data)
		if strings.Contains(src, str) {
			return true
		}
	}
	return false
}

func (c *Canary) Read(b []byte) (int, error) {
	data, ok := <-c.ch
	if !ok {
		c.ch = nil
		return 0, io.EOF
	}
	i := copy(b, data)
	if len(b) < len(data) {
		return len(b), io.ErrShortBuffer
	}
	return i, nil
}

func (c *Canary) Write(b []byte) (int, error) {
	if c.closed {
		return 0, io.ErrClosedPipe
	}
	data := make([]byte, len(b))
	i := copy(data, b)
	c.ch <- data
	return i, nil
}
