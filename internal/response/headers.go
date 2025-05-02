package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

func GetDefaultHeaders(contentLen int) headers.Headers {

	h := headers.NewHeaders()
	h.Set("content-length", strconv.Itoa(contentLen))
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")

	return h

}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	for key, val := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, val)))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	return err

}
