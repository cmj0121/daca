/* Copyright (C) 2019-2019 cmj. All right reserved. */
package logger

import (
	"bytes"
	"fmt"
	"testing"
)

var buffer []byte

type DummyWriter struct {
}

func (w DummyWriter) Write(in []byte) (n int, err error) {
	/* end of the log always be the NEWLINE */
	buffer = buffer[:0]
	buffer = append(buffer, in[:len(in)-1]...)
	return
}

func TestLogger(t *testing.T) {
	logger := New(DummyWriter{})
	logger.Format("{{.Message}}\n")

	cases := []struct {
		Level   int
		Message string
	}{
		{FATAL, "Fatal"},
		{CRIT, "Critical"},
	}

	for _, c := range cases {
		logger.Logf(c.Level, c.Message)
		if bytes.Compare(buffer, []byte(c.Message)) != 0 {
			t.Error(fmt.Sprintf("%v <> %v", string(c.Message), string(buffer)))
		}
	}
}

func TestLoggerFormat(t *testing.T) {
	logger := New(DummyWriter{})
	logger.Offset(1)

	cases := []struct {
		Format  string
		Message string
		Answer  string
	}{
		{"XXX\n", "Hello World", "XXX"},
		{"{{.Message}}\n", "Hello World", "Hello World"},
		{"{{.File}} L#{{.Line}} - {{.Message}}\n", "Hello World", "logger_test.go L#58 - Hello World"},
	}

	for _, c := range cases {
		logger.Format(c.Format)
		logger.Logf(CRIT, c.Message)
		if bytes.Compare(buffer, []byte(c.Answer)) != 0 {
			t.Error(fmt.Sprintf("%v <> %v", string(c.Answer), string(buffer)))
		}
	}
}

func BenchmarkLogger(b *testing.B) {
	logger := New(DummyWriter{})

	for i := 0; i < 1000000; i++ {
		logger.Logf(CRIT, "performance test : sequency")
	}
}
