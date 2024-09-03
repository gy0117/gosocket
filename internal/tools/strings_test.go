package tools

import (
	"math/rand"
	"testing"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano()) // 使用当前时间作为随机数生成器的种子
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func TestStringToBytes(t *testing.T) {
	t.Run("c1", func(t *testing.T) {
		s := "abc"
		b := StringToBytesStandard(s)
		t.Log("b1: ", b)
	})
	t.Run("c2", func(t *testing.T) {
		s := "abc"
		b2 := StringToBytesUnSafe(s)
		t.Log("b2: ", b2)
	})
	t.Run("c3", func(t *testing.T) {
		s := "abc"
		b3 := StringToBytesUnSafe2(s)
		t.Log("b3: ", b3)
	})
}

var fakeStrings = RandomString(100000)

func BenchmarkStringToBytesStandard(b *testing.B) {
	//b.Log("fakeStrings: ", fakeStrings)
	for i := 0; i < b.N; i++ {
		_ = StringToBytesStandard(fakeStrings)
	}
}

func BenchmarkStringToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = StringToBytesUnSafe(fakeStrings)
	}
}

func BenchmarkStringToBytesUnSafe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = StringToBytesUnSafe2(fakeStrings)
	}
}
