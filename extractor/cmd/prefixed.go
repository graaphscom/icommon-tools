package cmd

import "io"

func NewPrefixedWriter(prefix string, wrapped io.Writer) PrefixedWriter {
	return PrefixedWriter{prefix: []byte(prefix), wrapped: wrapped}
}

type PrefixedWriter struct {
	prefix  []byte
	wrapped io.Writer
}

func (pw PrefixedWriter) Write(p []byte) (int, error) {
	_, err := pw.wrapped.Write(append(pw.prefix, p...))
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
