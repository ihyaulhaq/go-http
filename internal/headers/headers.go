package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func isTokenChar(c byte) bool {
	return c == '!' || c == '#' || c == '$' || c == '%' ||
		c == '&' || c == '\'' || c == '*' || c == '+' ||
		c == '-' || c == '.' || c == '^' || c == '_' ||
		c == '`' || c == '|' || c == '~' ||
		(c >= '0' && c <= '9') ||
		(c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z')
}

func isValidToken(s string) bool {
	if len(s) == 0 {
		return false
	}

	for i := 0; i < len(s); i++ {
		if !isTokenChar(s[i]) {
			return false
		}
	}

	return true
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	str := string(data)

	idx := strings.Index(str, crlf)
	// not done
	if idx == -1 {
		return 0, false, nil
	}
	// header end
	if idx == 0 {
		return len(crlf), true, nil
	}

	line := str[:idx]
	consumed := idx + len(crlf)

	key, value, colon := strings.Cut(line, ":")
	if !colon {
		return 0, false, fmt.Errorf("headers format invalid")
	}

	if key != strings.TrimSpace(key) {
		return 0, false, fmt.Errorf("invalid header: whitespace not allowed around field name")
	}

	key = strings.ToLower(key)
	value = strings.TrimSpace(value)

	if !isValidToken(key) {
		return 0, false, fmt.Errorf("Invalid Headers: contain invalid char")
	}

	if val, ok := h[key]; ok {
		h[key] = val + "," + value
	} else {
		h[key] = value
	}

	return consumed, false, nil

}
