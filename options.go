package gosocket

import (
	"bufio"
)

var defaultServerOptions = &ServerOptions{}

// ServerOptions server configurations
type ServerOptions struct {
	bufReaderPool *Pool[*bufio.Reader]
}
