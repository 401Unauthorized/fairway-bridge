package Garmin_R10

import "bytes"

// jsonSplit is a bufio.SplitFunc that extracts complete JSON objects by counting curly braces.
func jsonSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip any leading whitespace.
	data = bytes.TrimLeft(data, " \t\r\n")
	if len(data) == 0 {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}
	// Look for a complete JSON object.
	var inString bool
	var escape bool
	var depth int
	start := -1
	for i, b := range data {
		if start == -1 && b == '{' {
			start = i
		}
		if start == -1 {
			continue
		}
		if b == '"' && !escape {
			inString = !inString
		}
		if inString {
			if b == '\\' && !escape {
				escape = true
			} else {
				escape = false
			}
			continue
		}
		if b == '{' {
			depth++
		} else if b == '}' {
			depth--
			if depth == 0 {
				// We have a complete JSON object.
				return i + 1, data[start : i+1], nil
			}
		}
	}
	// If we're at EOF, we may return what we have.
	if atEOF && start != -1 {
		return len(data), data[start:], nil
	}
	// Request more data.
	return 0, nil, nil
}
