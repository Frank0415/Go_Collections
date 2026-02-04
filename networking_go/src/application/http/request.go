package http

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
)

type Request struct {
	Method  string
	URI     string
	Version string
	Headers map[string]string
	Body    []byte
}

type ParseStatus int

const (
	StatusIncomplete ParseStatus = iota
	StatusComplete
	StatusError
)

type ParseResult struct {
	Status    ParseStatus
	Request   *Request
	Remaining []byte
	Consumed  int
	Err       error
}

func ParseInput(buf []byte) ParseResult {
	headerEnd := bytes.Index(buf, []byte("\r\n\r\n"))
	if headerEnd == -1 {
		return ParseResult{
			Status:    StatusIncomplete,
			Remaining: buf,
		}
	}

	headerLen := headerEnd + 4
	headerPart := buf[:headerEnd]

	req, contentLength, err := parseHeaders(headerPart)
	if err != nil {
		return ParseResult{
			Status: StatusError,
			Err:    err,
		}
	}

	totalLen := headerLen + contentLength
	if len(buf) < totalLen {
		return ParseResult{
			Status:    StatusIncomplete,
			Remaining: buf,
			Request:   req,
		}
	}

	if contentLength > 0 {
		req.Body = make([]byte, contentLength)
		copy(req.Body, buf[headerLen:totalLen])
	}

	remaining := buf[totalLen:]
	if len(remaining) == 0 {
		remaining = nil
	}

	return ParseResult{
		Status:    StatusComplete,
		Request:   req,
		Remaining: remaining,
		Consumed:  totalLen,
	}
}

func parseHeaders(headerPart []byte) (*Request, int, error) {
	lines := bytes.Split(headerPart, []byte("\r\n"))
	if len(lines) == 0 {
		return nil, 0, errors.New("empty request")
	}

	parts := bytes.SplitN(lines[0], []byte(" "), 3)
	if len(parts) != 3 {
		return nil, 0, fmt.Errorf("invalid request line: %s", lines[0])
	}

	req := &Request{
		Method:  string(parts[0]),
		URI:     string(parts[1]),
		Version: string(parts[2]),
		Headers: make(map[string]string),
	}

	contentLength := 0
	for _, line := range lines[1:] {
		if len(line) == 0 {
			continue
		}
		colonIdx := bytes.IndexByte(line, ':')
		if colonIdx == -1 {
			continue
		}
		key := string(line[:colonIdx])
		value := string(bytes.TrimSpace(line[colonIdx+1:]))
		req.Headers[key] = value

		if key == "Content-Length" {
			if cl, err := strconv.Atoi(value); err == nil && cl >= 0 {
				contentLength = cl
			}
		}
	}

	return req, contentLength, nil
}

func ShrinkBuffer(buf []byte, maxCap int) []byte {
	if cap(buf) <= maxCap {
		return buf
	}
	if len(buf) < cap(buf)/4 {
		newBuf := make([]byte, len(buf))
		copy(newBuf, buf)
		return newBuf
	}
	return buf
}

func LogInput(request *Request, logging bool) {
	if logging {
		log.Println(request.Method, " request.URI: ", request.URI)
		for i := range request.Headers {
			log.Println(i, ": ", request.Headers[i])
		}
		log.Println("Body: ", string(request.Body))
	}
}
