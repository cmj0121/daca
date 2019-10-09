/* Copyright (C) 2019-2019 cmj. All right reserved. */
package sqlfantasy

import (
	"fmt"
)

type Query struct {
	*Table

	limit  int /* limit of the query */
	offset int /* offset of the query */
}

func (q *Query) RawString(column string) (out string) {
	out = fmt.Sprintf(
		"SELECT %v FROM %s",
		column,
		q.Table.Name,
	)

	if q.offset > 0 {
		out = fmt.Sprintf(
			"%s OFFSET %d",
			out,
			q.limit,
		)
	}

	if q.limit > 0 {
		out = fmt.Sprintf(
			"%s LIMIT %d",
			out,
			q.limit,
		)
	}

	return
}

func (q *Query) Limit(limit int) (out *Query) {
	q.limit = limit
	out = q
	return
}

func (q *Query) Offset(offset int) (out *Query) {
	q.offset = offset
	out = q
	return
}
