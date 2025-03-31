package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	reqText := strings.Split(string(data), "\r\n")
	reqLine := reqText[0]

	parts := strings.Split(reqLine, " ")
	if len(parts) != 3 {
		return nil, errors.New("malformed request line in request")
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]

	switch method {
	case "GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE":
	default:
		return nil, fmt.Errorf("unrecognized request method: %s", method)
	}

	if version != "HTTP/1.1" {
		return nil, errors.New("version not supported - HTTP/1.1 only")
	}
	version = strings.Split(version, "/")[1]

	rl := RequestLine{
		HttpVersion:   version,
		RequestTarget: target,
		Method:        method,
	}

	return &Request{RequestLine: rl}, nil

}
