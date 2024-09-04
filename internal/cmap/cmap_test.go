package cmap

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCmp(t *testing.T) {
	t.Run("c1", func(t *testing.T) {
		var cmp = New[string, int](10, 15)
		assert.Equal(t, cmp.segmentSize, 16)
	})

	t.Run("c2", func(t *testing.T) {
		var cmp = New[string, int](10, 15)
		k1, k2 := "kk1", "kk2"

		go func() {
			cmp.Put(k1, 100)
		}()

		go func() {
			cmp.Put(k2, 200)
		}()

		go func() {
			value, ok := cmp.Get(k1)
			if ok {
				assert.Equal(t, value, 100)
			}
		}()
		time.Sleep(time.Second)
	})
}
