package internal

import (
	"github.com/dolthub/maphash"
	"sync"
)

// ConcurrentMap 并发安全map
type ConcurrentMap[K comparable, V any] struct {
	segments    []*Segment[K, V]
	segmentSize int
	mapHash     maphash.Hasher[K]
}

type Segment[K comparable, V any] struct {
	items        map[K]V
	sync.RWMutex // 读写锁
}

const defaultSegmentNum = 32

func New[K comparable, V any](segmentNum, segmentCap int) ConcurrentMap[K, V] {
	if segmentNum <= 0 {
		segmentNum = defaultSegmentNum
	}

	segmentNum = CeilPow2(segmentNum)

	cmp := ConcurrentMap[K, V]{
		segments:    make([]*Segment[K, V], segmentNum),
		segmentSize: segmentNum,
		mapHash:     maphash.NewHasher[K](),
	}
	for i := 0; i < cmp.segmentSize; i++ {
		cmp.segments[i] = &Segment[K, V]{items: make(map[K]V, segmentCap)}
	}
	return cmp
}

func (cmp *ConcurrentMap[K, V]) Get(key K) (value V, ok bool) {
	sg := cmp.getSegment(key)
	sg.RLock()
	defer sg.RUnlock()
	value, ok = sg.items[key]
	return
}

func (cmp *ConcurrentMap[K, V]) Put(key K, value V) {
	sg := cmp.getSegment(key)
	sg.Lock()
	defer sg.Unlock()
	sg.items[key] = value
}

func (cmp *ConcurrentMap[K, V]) Delete(key K) {
	sg := cmp.getSegment(key)
	sg.Lock()
	defer sg.Unlock()
	delete(sg.items, key)
}

func (cmp *ConcurrentMap[K, V]) getSegment(key K) *Segment[K, V] {
	hashCode := cmp.mapHash.Hash(key)
	index := hashCode & uint64(cmp.segmentSize-1)
	sg := cmp.segments[index]
	return sg
}
