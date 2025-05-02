package response

import (
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	STATUS_200 StatusCode = iota
	STATUS_400
	STATUS_500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {

	switch statusCode {
	case STATUS_200:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
	case STATUS_400:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
	case STATUS_500:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}
	}
	return nil

}

func GetDefaultHeaders(contentLen int) headers.Headers {

	h := headers.NewHeaders()
	h.Set("content-length", strconv.Itoa(contentLen))
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")

	return h

}
