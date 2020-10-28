package logger

import (
	"io"
	"os"
)

type (
	File interface {
		io.StringWriter
		io.Writer
		io.Closer
	}
	Filesystem interface {
		OpenFile(name string, flag int, perm os.FileMode) (File, error)
	}
	DefaultFilesystem struct {
	}
)

func NewDefaultFilesystem() *DefaultFilesystem {
	return &DefaultFilesystem{}
}

func (fs *DefaultFilesystem) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return os.OpenFile(name, flag, perm)
}
