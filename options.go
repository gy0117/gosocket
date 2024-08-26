package gows

import (
	"bufio"
	"github.com/gy/gows/internal"
)

var defaultServerOptions = &ServerOptions{}

// ServerOptions server configurations
type ServerOptions struct {
	bufReaderPool *internal.Pool[*bufio.Reader]
}

type Config struct {
}
