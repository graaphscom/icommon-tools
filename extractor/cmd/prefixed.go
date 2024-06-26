package cmd

import (
	"io"
	"sync"
)

func NewPrefixedWriter(prefix string, wrapped io.Writer) PrefixedWriter {
	return PrefixedWriter{
		prefix:        []byte(prefix),
		wrapped:       wrapped,
		lock:          &sync.Mutex{},
		atStartOfLine: true,
	}
}

type PrefixedWriter struct {
	prefix        []byte
	wrapped       io.Writer
	lock          *sync.Mutex
	atStartOfLine bool
}

func (pw PrefixedWriter) Write(p []byte) (int, error) {
	pw.lock.Lock()
	defer pw.lock.Unlock()

	var toWrite []byte

	for _, c := range p {
		if pw.atStartOfLine {
			toWrite = append(toWrite, pw.prefix...)
		}

		toWrite = append(toWrite, c)

		pw.atStartOfLine = c == '\n' || c == '\r'
	}

	_, err := pw.wrapped.Write(toWrite)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func NewPrefixedError(prefix string, wrapped error) error {
	if wrapped == nil {
		return nil
	}

	return PrefixedError{prefix, wrapped}
}

type PrefixedError struct {
	prefix  string
	wrapped error
}

func (pe PrefixedError) Error() string {
	return pe.prefix + pe.wrapped.Error()
}
