package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

const timeFormat = "15:04:05.000"

type EventLogger struct {
	w io.Writer
}

func New(w io.Writer) (*EventLogger, error) {
	if w == nil {
		return nil, errors.New("nil values in EventLogger constructor")
	}

	return &EventLogger{w: w}, nil
}

func (l *EventLogger) Log(t time.Time, msg string) {
	_, err := fmt.Fprintf(l.w, "[%s] %s\n", t.Format(timeFormat), msg)
	if err != nil {
		log.Fatal("failed to write log")
	}
}
