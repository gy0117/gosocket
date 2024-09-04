package gosocket

import (
	"bytes"
	"compress/flate"
	"io"
)

func compressMessage(message []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := flate.NewWriter(&buf, flate.BestCompression)
	if err != nil {
		return nil, err
	}

	if _, err = writer.Write(message); err != nil {
		return nil, err
	}
	_ = writer.Close()

	return buf.Bytes(), nil
}

func decompressMessage(message []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(message))

	var result bytes.Buffer
	if _, err := io.Copy(&result, reader); err != nil {
		return nil, err
	}
	_ = reader.Close()
	return result.Bytes(), nil
}
