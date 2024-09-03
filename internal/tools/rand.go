package tools

import (
	"crypto/rand"
	"io"
)

func GenerateMaskingKey() ([]byte, error) {
	// 创建一个4字节的slice用于存储masking key
	key := make([]byte, 4)

	// 使用crypto/rand包生成4字节的随机值
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}
