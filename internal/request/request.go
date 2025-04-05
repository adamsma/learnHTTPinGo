package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	parseState  requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	INITIALIZED requestState = iota
	DONE
)

const crlf = "\r\n"
const BUFFER_SIZE = 8

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf := make([]byte, BUFFER_SIZE)
	readToIndex := 0

	req := &Request{
		parseState: INITIALIZED,
	}

	for req.parseState != DONE {

		if readToIndex >= len(buf) {
			newBuf := make([]byte, 2*len(buf), 2*cap(buf))
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])

		readToIndex += n

		if errors.Is(err, io.EOF) {
			req.parseState = DONE
			break
		}

		if err != nil {
			return nil, err
		}

		n, err = req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[n:])
		readToIndex -= n

	}

	return req, nil

}

func (r *Request) parse(data []byte) (int, error) {

	if r.parseState == INITIALIZED {
		n, rl, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *rl
		r.parseState = DONE

		return n, nil
	}

	if r.parseState == DONE {
		return 0, errors.New("error: tyring to read data in a done state")
	}

	return 0, fmt.Errorf("unknown state: %d", r.parseState)

}

func parseRequestLine(data []byte) (int, *RequestLine, error) {

	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, nil, nil
	}

	reqLineText := string(data[:idx])
	rl, err := requestLineFromString(reqLineText)
	if err != nil {
		return 0, nil, err
	}

	return idx, rl, nil

}

func requestLineFromString(str string) (*RequestLine, error) {

	parts := strings.Split(str, " ")
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

	versionParts := strings.Split(version, "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version = versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	rl := &RequestLine{
		HttpVersion:   version,
		RequestTarget: target,
		Method:        method,
	}

	return rl, nil

}
