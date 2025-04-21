package request

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	RequestLine    RequestLine
	parseState     requestState
	Headers        headers.Headers
	Body           []byte
	bodyLengthRead int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	INITIALIZED requestState = iota
	PARSING_HEADERS
	PARSING_BODY
	DONE
)

const crlf = "\r\n"
const BUFFER_SIZE = 8

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf := make([]byte, BUFFER_SIZE)
	readToIndex := 0

	req := &Request{
		parseState: INITIALIZED,
		Headers:    headers.NewHeaders(),
		Body:       make([]byte, 0),
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
			if req.parseState != DONE {
				return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", req.parseState, n)
			}
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

	totalBytesParsed := 0
	for r.parseState != DONE {

		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return n, err
		}

		totalBytesParsed += n

		// need more data
		if n == 0 {
			break
		}

	}

	return totalBytesParsed, nil

}

func (r *Request) parseSingle(data []byte) (int, error) {

	switch r.parseState {
	case INITIALIZED:

		n, rl, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		// need more data
		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *rl
		r.parseState = PARSING_HEADERS

		return n, nil

	case PARSING_HEADERS:

		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if done {
			r.parseState = PARSING_BODY
			return n, nil
		}

		return n, nil

	case PARSING_BODY:

		// if we are parsing the body, we expect to have a Content-Length header
		contentLengthVal, exists := r.Headers.Get("Content-Length")

		if !exists {
			// if there is no Content-Length header, we assume the body is empty
			r.parseState = DONE
			return len(data), nil
		}

		contentLength, err := strconv.Atoi(contentLengthVal)
		if err != nil {
			return 0, fmt.Errorf("invalid Content-Length header: %s", contentLengthVal)
		}

		if contentLength == 0 {
			r.parseState = DONE
			return len(data), nil
		}

		r.Body = append(r.Body, data...)
		r.bodyLengthRead += len(data)

		if r.bodyLengthRead > contentLength {
			return len(r.Body), fmt.Errorf("error: body length %d exceeds Content-Length %d", len(r.Body), contentLength)
		}

		if r.bodyLengthRead == contentLength {
			r.parseState = DONE
		}

		return len(data), nil

	case DONE:
		return 0, errors.New("error: tyring to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state: %d", r.parseState)
	}

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

	return idx + 2, rl, nil

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
