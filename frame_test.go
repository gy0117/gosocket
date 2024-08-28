package gosocket

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestParseHeader(t *testing.T) {
	t.Run("c1", func(t *testing.T) {
		s, c := net.Pipe()
		_ = c.Close()

		var f = Frame{}
		var _, err = f.ParseHeader(bufio.NewReader(s))
		assert.Error(t, err)
	})

	t.Run("c2", func(t *testing.T) {
		s, c := net.Pipe()
		go func() {
			f := Frame{}
			f.CreateHeader(true, OpcodeTextFrame, false, 100)
			c.Write(f.Header[:2])
			c.Close()
		}()

		time.Sleep(100 * time.Millisecond)
		var f = Frame{}
		var _, err = f.ParseHeader(bufio.NewReader(s))
		assert.NoError(t, err)
	})

}
