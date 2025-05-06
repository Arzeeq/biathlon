package logger_test

import (
	"biathlon/internal/logger"
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	time := time.Date(0, 1, 1, 9, 5, 59, 867000000, time.UTC)
	msg := "my message"

	var logBuf bytes.Buffer
	l, err := logger.New(&logBuf)
	require.NoError(t, err)

	l.Log(time, msg)

	require.Equal(t, "[09:05:59.867] my message\n", logBuf.String())
}
