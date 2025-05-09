package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Set(key, value string) {
	h[strings.ToLower(key)] = value
}

func (h Headers) Get(key string) (value string, exists bool) {
	value, exists = h[strings.ToLower(key)]
	return
}

func (h Headers) Delete(key string) {
	delete(h, strings.ToLower(key))
}

const crlf = "\r\n"
const validCharRegexp = "^[a-z0-9!#$%&'*+.^_`|~-]*$"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	if strings.HasPrefix(string(data), crlf) {
		return 2, true, nil
	}

	n, header, value, err := parseHeaderLine(data)
	if err != nil {
		return 0, false, err
	}

	// need more data
	if n == 0 {
		return 0, false, nil
	}

	if curVal, exists := h[header]; exists {
		value = curVal + ", " + value
	}

	h.Set(header, value)

	return n, false, nil
}

func parseHeaderLine(data []byte) (n int, key string, value string, err error) {

	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, "", "", nil
	}

	headerLineText := string(data[:idx])
	key, value, found := strings.Cut(headerLineText, ":")
	if !found {
		return 0, "", "", fmt.Errorf("malformed header: %s", headerLineText)
	}

	if key != strings.TrimRight(key, " ") {
		return 0, "", "", fmt.Errorf("malformed header: %s", headerLineText)
	}

	key = strings.ToLower(strings.TrimSpace(key))
	value = strings.TrimSpace(value)

	if matched, _ := regexp.MatchString(validCharRegexp, key); !matched {
		err = fmt.Errorf("header field name contains invalid character(s): %s", key)
		return 0, "", "", err
	}

	return idx + len(crlf), key, value, nil

}
