package common

import "io"

// FormatConverter represents common converters between different formats
type FormatConverter interface {
	Convert(w io.Writer, r io.Reader) error
}
