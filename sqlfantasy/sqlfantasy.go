/* Copyright (C) 2019-2019 cmj. All right reserved. */
package sqlfantasy

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
)

const (
	TAG_TABLE  = "table"
	TAG_COLUMN = "column"
)

type SF struct {
	Conn *sql.DB /* Lower-end DB connection */
}

func Open(driver, conn string) (out *SF, err error) {
	out = &SF{}

	out.Conn, err = sql.Open(driver, conn)
	return
}

func MustOpen(driver, conn string) (out *SF) {
	out, err := Open(driver, conn)
	if err != nil {
		/* raise panic */
		panic(err)
	}
	return
}

func (sf *SF) Close() {
	/* close the connection */
	sf.Conn.Close()
}

func ParseRows(rows *sql.Rows, table *Table) (out interface{}, err error) {
	objlist := reflect.New(reflect.SliceOf(reflect.TypeOf(table)))

	params := make([]string, len(table.Fields))
	ref_params := []interface{}{}

	for idx := 0; idx < len(params); idx++ {
		ref_params = append(ref_params, &params[idx])
	}

	for rows.Next() {
		if err = rows.Scan(ref_params...); err != nil {
			/* scan failure, return */
			return
		}

		tmp := reflect.New(reflect.TypeOf(table.source))
		for idx, column := range table.Columns() {
			/* Set field */
			field := tmp.Elem().FieldByName(column)
			field.Set(reflect.ValueOf(params[idx]))
		}
		reflect.Append(objlist, tmp)
	}

	out = objlist.Interface()
	return
}
