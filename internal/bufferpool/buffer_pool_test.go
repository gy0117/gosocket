package bufferpool

import (
	"fmt"
	"testing"
)

func TestNewBufferPool(t *testing.T) {
	bufferPool := NewBufferPools(64, 256)
	buffer := bufferPool.Get(128)
	fmt.Println(buffer)
}
