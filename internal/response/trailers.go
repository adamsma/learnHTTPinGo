package response

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
)

func (w *Writer) WriteTrailers(headers headers.Headers) error {

	if w.writeState != WriteStateTrailers {
		return errors.New("cannot write trailers in current write state")
	}

	defer func() { w.writeState = WriteStateBody }()
	for k, v := range headers {
		_, err := w.conn.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.conn.Write([]byte("\r\n"))
	return err

}
