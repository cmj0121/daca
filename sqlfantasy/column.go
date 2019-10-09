/* Copyright (C) 2019-2019 cmj. All right reserved. */
package sqlfantasy

import (
	"reflect"
)

type Column struct {
	Name string       /* Column name */
	Kind reflect.Kind /* Kind of the column */
}
