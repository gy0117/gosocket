package tools

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/gy/gosocket/internal/types"
)

// CeilPow2 2^k >= x，返回最小的2^k
func CeilPow2(n int) int {
	x := 1
	for x < n {
		x *= 2
	}
	return x
}

// 按照 WebSocket 协议要求，与固定的 GUID 拼接
var guid = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

// GetSecWebSocketAccept base64 websocket key
func GetSecWebSocketAccept(key string) string {
	b := []byte(key)
	b = append(b, guid...)

	hash := sha1.New()
	hash.Write(b)
	hashed := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashed)
}

func GetSecWebSocketExtensions() string {
	return types.PermessageDeflate
}
