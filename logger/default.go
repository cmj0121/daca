/* Copyright (C) 2019-2019 cmj. All right reserved. */
package logger

var default_logger *Logger

/* shortcut routine */
func Fatal(msg string, args ...interface{}) {
	default_logger.Logf(FATAL, msg, args...)
	return
}

func Crit(msg string, args ...interface{}) {
	default_logger.Logf(CRIT, msg, args...)
	return
}

func Error(msg string, args ...interface{}) {
	default_logger.Logf(ERROR, msg, args...)
	return
}

func Warn(msg string, args ...interface{}) {
	default_logger.Logf(WARN, msg, args...)
	return
}

func Info(msg string, args ...interface{}) {
	default_logger.Logf(INFO, msg, args...)
	return
}

func Debug(msg string, args ...interface{}) {
	default_logger.Logf(DEBUG, msg, args...)
	return
}
