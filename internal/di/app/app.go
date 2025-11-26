package app

import (
	"io"
)

type App interface {
	Serve(in io.Reader, out io.Writer) int
}
