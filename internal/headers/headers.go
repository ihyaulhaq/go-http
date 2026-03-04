package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

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

	value = strings.TrimSpace(value)

	if key != strings.TrimSpace(key) {
		return 0, false, fmt.Errorf("invalid header: whitespace not allowed around field name")
	}

	h[key] = value

	return consumed, false, nil

}
