package logger

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type BasicFormatter struct {
	appName    string
	dateFormat string
}

func NewBasicFormatter(appName string, dateFormat string) *BasicFormatter {
	return &BasicFormatter{appName: appName, dateFormat: dateFormat}
}

func (f *BasicFormatter) Format(log Log) FormattedLog {
	return FormattedLog{
		Log:              log,
		FormattedMessage: f.createFormattedMessage(log),
	}
}

func (f *BasicFormatter) createFormattedMessage(log Log) string {
	ln, lnErr := log.Level.Name()

	if errors.Is(lnErr, ErrLevelNameMappingNotFound) {
		ln = "unknown"
	}

	res := &strings.Builder{}
	res.WriteString(f.appName)
	res.WriteString(" | ")
	res.WriteString(log.CreatedAt.Format(f.dateFormat))
	res.WriteString(" | ")
	res.WriteString(strings.ToUpper(ln.String()))
	res.WriteString(" | ")
	res.WriteString(f.applyDataOnMessage(log.Message, log.Data))

	if log.Data.Len() > 0 {
		res.WriteString(" | ")
		res.WriteString(f.createDataSection(log.Data))
	}

	return res.String()
}

func (f *BasicFormatter) applyDataOnMessage(message string, data Data) string {
	if data.Len() < 1 {
		return message
	}

	for key, val := range data {
		message = strings.ReplaceAll(message, fmt.Sprintf("{%s}", key), val)
	}

	return message
}

func (f *BasicFormatter) createDataSection(data Data) string {
	var keys []string
	res := &strings.Builder{}

	for key := range data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		res.WriteString(key)
		res.WriteString(":")
		res.WriteString(data[key])
		res.WriteString(" ")
	}

	text := res.String()

	return text[:len(text)-1]
}
