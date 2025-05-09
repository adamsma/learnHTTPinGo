package response

import (
	"errors"
	"fmt"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {

	if w.writeState != WriteStateBody {
		return 0, errors.New("cannot write headers in current write state")
	}

	chunkLen := fmt.Sprintf("%X\r\n", len(p))
	n1, err := w.conn.Write([]byte(chunkLen))
	if err != nil {
		return n1, err
	}

	n2, err := w.conn.Write(p)
	w.conn.Write([]byte("\r\n"))

	return n1 + n2 + 2, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {

	if w.writeState != WriteStateBody {
		return 0, errors.New("cannot write headers in current write state")
	}

	n, err := w.conn.Write([]byte("0\r\n"))

	w.writeState = WriteStateTrailers

	return n, err
}
