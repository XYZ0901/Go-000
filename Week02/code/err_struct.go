package code

import (
	"fmt"
	"io"
)

type errWriter struct {
	io.Writer
	err error
}

func (e *errWriter) Write(buf []byte) (int, error) {
	if e.err != nil {
		return 0, e.err
	}

	var n int
	n, e.err = e.Writer.Write(buf)
	return n, nil
}

type Header struct {
	Key, Value string
}

type Status struct {
	Code   int
	Reason string
}

func WriteResponse(w io.Writer, st Status, headers []Header, body io.Reader) error {
	ew := &errWriter{Writer: w}
	_, _ = fmt.Fprintf(ew, "HTTP/1.1 %d %s\r\n", st.Code, st.Reason)
	for _, h := range headers {
		_, _ = fmt.Fprintf(ew, "%s: %s\r\n", h.Key, h.Value)
	}
	_, _ = fmt.Fprint(ew, "\r\n")
	_, _ = io.Copy(ew, body)
	return ew.err
}
