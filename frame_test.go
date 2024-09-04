package gosocket

import (
	"bufio"
	"fmt"
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
			f.CreateHeader(true, OpcodeTextFrame, false, 100, false)
			c.Write(f.Header[:2])
			c.Close()
		}()

		time.Sleep(100 * time.Millisecond)
		var f = Frame{}
		var _, err = f.ParseHeader(bufio.NewReader(s))
		assert.NoError(t, err)
	})

}

func TestFrame(t *testing.T) {
	frame := []byte{0x01, 0x00}
	processWebSocketFrame(frame)
}

func processWebSocketFrame(frame []byte) {
	// 提取第一个字节，FIN 位在该字节的最高位
	firstByte := frame[0]
	fin := firstByte & 0x80
	opcode := firstByte & 0x0F

	if fin != 0 {
		fmt.Println("This is the final frame of the message. fin: ", fin)
	} else {
		fmt.Println("More frames are expected for this message. fin: ", fin)
	}

	fmt.Printf("Opcode: %d\n", opcode)
}
