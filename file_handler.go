package logger

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// FileHandler writes all logs to provided log file.
type FileHandler struct {
	fileLock   sync.Mutex
	file       File
	filesystem Filesystem
	path       string
	formatter  Formatter
}

// NewFileHandler creates FileHandler with
// main log file expected to be at "path" location.
func NewFileHandler(path string) *FileHandler {
	return &FileHandler{path: path, filesystem: NewDefaultFilesystem()}
}

// WithFilesystem allows to replace default Filesystem
// implementation with custom made.
func (f *FileHandler) WithFilesystem(filesystem Filesystem) {
	f.filesystem = filesystem
}

func (f *FileHandler) UseFormatter(formatter Formatter) {
	f.formatter = formatter
}

func (f *FileHandler) Handle(log Log) error {
	f.fileLock.Lock()
	defer f.fileLock.Unlock()

	var message string
	ofErr := f.openFile()

	if ofErr != nil {
		return ofErr
	}

	if f.formatter != nil {
		message = f.formatter.Format(log).FormattedMessage
	} else {
		message = log.Message
	}

	_, wErr := f.file.WriteString(fmt.Sprintf("%s\n", message))

	if wErr != nil {
		return fmt.Errorf("uanble to write log to file: %w", wErr)
	}

	return wErr
}

func (f *FileHandler) HandleBatch(logs []Log) error {
	f.fileLock.Lock()
	defer f.fileLock.Unlock()

	ofErr := f.openFile()

	if ofErr != nil {
		return ofErr
	}

	writer := bufio.NewWriter(f.file)

	for _, l := range logs {
		var message string
		if f.formatter != nil {
			message = f.formatter.Format(l).FormattedMessage
		} else {
			message = l.Message
		}

		_, bwErr := writer.WriteString(fmt.Sprintf("%s\n", message))

		if bwErr != nil {
			return fmt.Errorf("uanble to write log to buffer: %w", bwErr)
		}
	}

	fwErr := writer.Flush()

	if fwErr != nil {
		return fmt.Errorf("uanble to write log to file: %w", fwErr)
	}

	return nil
}

func (f *FileHandler) Close() error {
	f.fileLock.Lock()
	defer f.fileLock.Unlock()

	if f.file == nil {
		return nil
	}

	return f.file.Close()
}

func (f *FileHandler) openFile() error {
	if f.file != nil {
		return nil
	}

	var lfErr error
	f.file, lfErr = f.filesystem.OpenFile(f.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0644))

	if lfErr != nil {
		return fmt.Errorf("unable to open log file: %w", lfErr)
	}

	return nil
}
