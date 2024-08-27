package gows

import (
	"bufio"
	"github.com/gy/gosocket/internal"
)

var defaultServerOptions = &ServerOptions{}

// ServerOptions server configurations
type ServerOptions struct {
	bufReaderPool *internal.Pool[*bufio.Reader]
}
