/* Copyright (C) 2019-2019 cmj. All right reserved. */
package logger

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"sync"
	"text/template"
	"time"
)

const (
	/* The log level, smaller means important */
	FATAL    = 0
	CRIT     = 10
	CRITICAL = 10
	ERROR    = 20
	WARN     = 30
	WARNING  = 30
	INFO     = 40
	DEBUG    = 50
)

/* The logger instance which control and display the log message */
type Logger struct {
	lock      *sync.Mutex        /* ensure the atomic writes */
	writer    io.Writer          /* The final I/O writer */
	level     int                /* Set the display log level */
	tmpl      *template.Template /* format template */
	pc_offset int                /* the offset of the caller (default: 2) */
}

func New(writer io.Writer) (log *Logger) {
	log = &Logger{
		lock:      &sync.Mutex{},
		writer:    writer,
		pc_offset: 2,
	}

	log.Level(WARNING)
	log.Format("[{{.Now.Format \"2006-01-02 15:04:05\"}}] {{.File}} L#{{.Line}} - {{.Message}}\n")
	return
}

func DefaultLogger() (log *Logger) {
	log = default_logger
	return
}

func (log *Logger) Offset(offset int) (out *Logger) {
	log.pc_offset = offset
	out = log
	return
}

func (log *Logger) Level(lv int) (out *Logger) {
	log.level = lv
	out = log
	return
}

func (log *Logger) Format(tmpl string) (out *Logger) {
	log.tmpl = template.Must(template.New("logger").Parse(tmpl))
	out = log
	return
}

func (log *Logger) Logf(lv int, msg string, args ...interface{}) {
	log.lock.Lock()
	defer log.lock.Unlock()

	if lv <= log.level {
		msg := fmt.Sprintf(msg, args...)

		var level string

		switch {
		case FATAL >= lv:
			level = "FATAL"
		case CRIT >= lv:
			level = " CRIT"
		case ERROR >= lv:
			level = "ERROR"
		case WARN >= lv:
			level = " WARN"
		case INFO >= lv:
			level = " INFO"
		default:
			level = "DEBUG"
		}

		/* specified log template */
		_, file, line, _ := runtime.Caller(log.pc_offset)

		if log.tmpl != nil {
			data := struct {
				Level   string    /* Log level: str */
				Message string    /* main message */
				Now     time.Time /* timestamp for now */
				File    string    /* The source file*/
				Dir     string    /* The source file*/
				Line    int       /* The line number of source file */
			}{
				Level:   level,
				Message: msg,
				Now:     time.Now(),
				File:    filepath.Base(file),
				Dir:     filepath.Dir(file),
				Line:    line,
			}

			/* Write the executed template to writer */
			var buff bytes.Buffer
			if err := log.tmpl.Execute(&buff, data); err != nil {
				fmt.Println("Cannot display the log", err)
				return
			}

			msg = buff.String()
		}

		log.writer.Write([]byte(msg))
	}
}

func (log *Logger) Fatal(msg string, args ...interface{}) {
	log.Logf(FATAL, msg, args...)
	return
}

func (log *Logger) Crit(msg string, args ...interface{}) {
	log.Logf(CRIT, msg, args...)
	return
}

func (log *Logger) Warn(msg string, args ...interface{}) {
	log.Logf(WARN, msg, args...)
	return
}

func (log *Logger) Info(msg string, args ...interface{}) {
	log.Logf(INFO, msg, args...)
	return
}

func (log *Logger) Debug(msg string, args ...interface{}) {
	log.Logf(DEBUG, msg, args...)
	return
}
