package response

import (
	"httpfromtcp/internal/headers"
	"io"

	"errors"
)

type WriteState int

const (
	WriteStateStatusLine = iota
	WriteStateHeaders
	WriteStateBody
	WriteStateClosed
)

type Writer struct {
	conn       io.Writer
	writeState WriteState
}

func NewWriter(w io.Writer) Writer {
	return Writer{conn: w, writeState: WriteStateStatusLine}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	if w.writeState != WriteStateStatusLine {
		return errors.New("cannot write status line in current write state")
	}

	err := WriteStatusLine(w.conn, statusCode)
	w.writeState = WriteStateHeaders

	return err

}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	if w.writeState != WriteStateHeaders {
		return errors.New("cannot write headers in current write state")
	}

	err := WriteHeaders(w.conn, headers)
	w.writeState = WriteStateBody

	return err

}

func (w *Writer) WriteBody(p []byte) (int, error) {

	if w.writeState != WriteStateBody {
		return 0, errors.New("cannot write headers in current write state")
	}

	n, err := w.conn.Write(p)
	w.writeState = WriteStateClosed

	return n, err

}
