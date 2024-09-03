package gosocket

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateFrame(t *testing.T) {
	t.Run("c1", func(t *testing.T) {
		wsConn := WsConn{
			server: false,
		}
		frame, err := wsConn.createFrame(OpcodeTextFrame, []byte("abc"))
		t.Log("frame", string(frame.Bytes()))
		assert.NoError(t, err)
	})
}
