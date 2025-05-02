package response

import "io"

type StatusCode int

const (
	StatusCodeSuccess             StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

func getStatusLine(statusCode StatusCode) []byte {

	switch statusCode {
	case StatusCodeSuccess:
		return []byte("HTTP/1.1 200 OK\r\n")
	case StatusCodeBadRequest:
		return []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusCodeInternalServerError:
		return []byte("HTTP/1.1 500 Internal Server Error\r\n")
	}
	return nil

}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write(getStatusLine(statusCode))
	return err
}
