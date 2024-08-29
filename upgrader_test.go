package gosocket

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"testing"
)

var _ http.ResponseWriter = (*httpWriter)(nil)

type httpWriter struct {
	conn net.Conn
	brw  *bufio.ReadWriter
}

func newHttpWriter() *httpWriter {
	server, _ := net.Pipe()
	var r = bytes.NewBuffer(nil)
	var w = bytes.NewBuffer(nil)
	var brw = bufio.NewReadWriter(bufio.NewReader(r), bufio.NewWriter(w))

	return &httpWriter{
		conn: server,
		brw:  brw,
	}
}

func (h *httpWriter) Header() http.Header {
	return http.Header{}
}

func (h *httpWriter) Write(bytes []byte) (int, error) {
	return 0, nil
}

func (h *httpWriter) WriteHeader(statusCode int) {

}

func (h *httpWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return h.conn, h.brw, nil
}

func TestHijack(t *testing.T) {
	var upgrade = NewUpgrade(&ServerOptions{
		BufReaderPool: NewPool(func() *bufio.Reader {
			return bufio.NewReaderSize(nil, 1024)
		}),
	})
	h := http.Header{}
	h.Set("Connection", "Upgrade")
	h.Set("Upgrade", "websocket")
	h.Set("Sec-WebSocket-Version", "13")

	var r = &http.Request{
		Header: h,
		Method: http.MethodGet,
	}
	_, err := upgrade.Upgrade(newHttpWriter(), r)
	assert.NoError(t, err)

	h1 := http.Header{}
	h1.Set("Connection", "upgrade")
	h1.Set("Upgrade", "websockt")
	h1.Set("Sec-WebSocket-Version", "13")

	var r1 = &http.Request{
		Header: h1,
		Method: http.MethodGet,
	}
	_, err1 := upgrade.Upgrade(newHttpWriter(), r1)
	assert.Error(t, err1)
}
