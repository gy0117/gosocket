package gosocket

import (
	"fmt"
	"testing"
)

func TestNewBufferPool(t *testing.T) {
	bufferPool := NewBufferPool(64, 256)
	buffer := bufferPool.Get(128)
	fmt.Println(buffer)
}
