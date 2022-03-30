package common

import "io"

type FormatConverter interface {
	Convert(w io.Writer, r io.Reader) error
}
