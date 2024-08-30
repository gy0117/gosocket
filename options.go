package gosocket

import (
	"bufio"
	"net/http"
)

const (
	defaultReaderBufSize       = 1024 * 4
	defaultMaxWritePayloadSize = 1024 * 1024
	defaultMaxReadPayloadSize  = 1024 * 1024
)

// Config 连接配置与WsConn绑定
type Config struct {
	ReaderBufSize       int  // 读缓冲区大小
	MaxReadPayloadSize  int  // 最大的读取消息长度
	MaxWritePayloadSize int  // 最大的写入消息长度
	OpenUTF8Check       bool // 是否打开utf-8编码检查
}

// ServerOptions server configurations
// 下面的属性都有默认实现
type ServerOptions struct {
	readerBufPool *Pool[*bufio.Reader] // 读缓冲区

	ReaderBufSize       int                                            // 读缓冲区大小
	MaxReadPayloadSize  int                                            // 最大的读取消息长度
	MaxWritePayloadSize int                                            // 最大的写入消息长度
	OpenUTF8Check       bool                                           // 是否打开utf-8编码检查
	NewSessionMap       func() SessionManager                          // 创建sm，管理当前连接的信息
	PreSessionHandle    func(r *http.Request, sm SessionManager) error // 预处理请求
}

func (so *ServerOptions) CreateConfig() *Config {
	return &Config{
		ReaderBufSize:       so.ReaderBufSize,
		MaxReadPayloadSize:  so.MaxReadPayloadSize,
		MaxWritePayloadSize: so.MaxWritePayloadSize,
		OpenUTF8Check:       so.OpenUTF8Check,
	}
}

// SessionManager 连接管理接口
type SessionManager interface {
	Get(key string) (value any, ok bool)
	Put(key string, value any)
	Delete(key string)
}
