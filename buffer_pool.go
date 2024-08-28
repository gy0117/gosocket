package gosocket

import (
	"bytes"
	"math"
	"sync"
)

var bufferPool = new(BufferPool)

func init() {
	bufferPool = NewBufferPool(64, 128*1024)
}

// BufferPool 缓冲池，创建不同大小的缓冲池，并使用sync.Pool管理
type BufferPool struct {
	pools      map[int]*sync.Pool
	start, end int
}

func NewBufferPool(start, end int) *BufferPool {
	start = nextPowerOfTwo(start)
	end = nextPowerOfTwo(end)

	bp := &BufferPool{
		pools: make(map[int]*sync.Pool),
		start: start,
		end:   end,
	}

	for i := start; i <= end; i *= 2 {
		bp.pools[i] = newPool(i)
	}

	return bp
}

func (bp *BufferPool) Put(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	if pool, ok := bp.pools[buf.Cap()]; ok {
		pool.Put(buf)
	}
}

func (bp *BufferPool) Get(n int) *bytes.Buffer {
	size := max(bp.start, nextPowerOfTwo(n))
	if pool, ok := bp.pools[size]; ok {
		buf := pool.Get().(*bytes.Buffer)

		buf.Reset()
		return buf
	}
	return bytes.NewBuffer(make([]byte, 0, n))
}

func newPool(cap int) *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, cap))
		},
	}
}

// 向上取整，2的幂次，比如：n = 10，返回16，10 位于 2^3 ~ 2^4之间
func nextPowerOfTwo(n int) int {
	if n <= 0 {
		return 1
	}
	// 使用 math.Ceil 和 math.Log2 计算最接近的 2 的幂
	// math.Log2：计算n的对数；math.Ceil：对对数结果进行向上取整，得到最小的整数k，使得2^k >= n
	// math.Pow：计算2的k次幂
	return int(math.Pow(2, math.Ceil(math.Log2(float64(n)))))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
