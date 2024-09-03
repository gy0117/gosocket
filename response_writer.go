package gosocket

import (
	"bytes"
	"github.com/gy/gosocket/pkg/bufferpool"
	"net"
)

type ResponseWriter struct {
	buf *bytes.Buffer
}

func NewResponseWriter() *ResponseWriter {
	rw := &ResponseWriter{
		buf: bufferpool.Pools.Get(1024),
	}
	rw.buf.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	rw.buf.WriteString("Upgrade: websocket\r\n")
	rw.buf.WriteString("Connection: Upgrade\r\n")
	return rw
}

func (w *ResponseWriter) Close() {
	bufferpool.Pools.Put(w.buf)
	w.buf = nil
}

func (w *ResponseWriter) AddHeader(key string, value string) {
	w.buf.WriteString(key)
	w.buf.WriteString(": ")
	w.buf.WriteString(value)
	w.buf.WriteString("\r\n")
}

// TODO 需要设置超时
func (w *ResponseWriter) Write(conn net.Conn) error {
	w.buf.WriteString("\r\n")
	if _, err := w.buf.WriteTo(conn); err != nil {
		return err
	}
	return nil
}
