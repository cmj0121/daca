/* Copyright (C) 2019-2019 cmj. All right reserved. */
package logger

import (
	"os"
)

func init() {
	default_logger = New(os.Stderr)
	return
}
