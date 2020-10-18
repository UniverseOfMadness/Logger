package logger

import (
	"fmt"
	"io"
	"strings"
)

type StringWriterHandler struct {
	writer    io.StringWriter
	formatter Formatter
}

func NewStringWriterHandler(writer io.StringWriter) *StringWriterHandler {
	return &StringWriterHandler{writer: writer}
}

func (h *StringWriterHandler) UseFormatter(formatter Formatter) *StringWriterHandler {
	h.formatter = formatter

	return h
}

func (h *StringWriterHandler) Handle(log Log) error {
	var message string

	if h.formatter != nil {
		message = h.formatter.Format(log).FormattedMessage
	} else {
		message = log.Message
	}

	_, err := h.writer.WriteString(fmt.Sprintf("%s\n", message))

	if err != nil {
		return fmt.Errorf("StringWriterHandler - error occurred while handling log: %w", err)
	}

	return nil
}

func (h *StringWriterHandler) HandleBatch(logs []Log) error {
	var buff []string

	for _, log := range logs {
		var message string

		if h.formatter != nil {
			message = h.formatter.Format(log).FormattedMessage
		} else {
			message = log.Message
		}

		buff = append(buff, message)
	}

	_, err := h.writer.WriteString(fmt.Sprintf("%s\n", strings.Join(buff, "\n")))

	if err != nil {
		return fmt.Errorf("StringWriterHandler - error occurred while handling logs: %w", err)
	}

	return nil
}
