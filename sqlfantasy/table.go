/* Copyright (C) 2019-2019 cmj. All right reserved. */
package sqlfantasy

import (
	"fmt"
	"reflect"
)

type Table struct {
	source interface{}       /* The source object */
	Name   string            /* The table name */
	Fields map[string]Column /* Column of the record */
}

func ParseTable(in interface{}) (out *Table) {
	var obj reflect.Type

	switch typ := reflect.TypeOf(in); typ.Kind() {
	case reflect.Ptr:
		obj = typ.Elem()
	case reflect.Struct:
		obj = typ
	default:
		panic(fmt.Errorf("%v cannot be table, need pointer of struct", in))
	}

	out = &Table{
		source: in,
		Fields: make(map[string]Column, 0),
	}

	for idx := 0; idx < obj.NumField(); idx++ {
		field := obj.Field(idx)

		if v, ok := field.Tag.Lookup(TAG_TABLE); ok {
			if out.Name != "" {
				/* duplicated table name declaim */
				panic(fmt.Errorf("%v has duplicate %s tag", in, TAG_TABLE))
			}
			out.Name = v
		}

		if v, ok := field.Tag.Lookup(TAG_COLUMN); ok {
			out.Fields[field.Name] = Column{
				Name: v,                 /* table column name */
				Kind: field.Type.Kind(), /* Kind of the column */
			}
		}
	}

	return
}

func (t *Table) Query() (out *Query) {
	out = &Query{
		Table: t,
	}
	return
}

func (t *Table) Columns() (out []string) {
	for _, column := range t.Fields {
		out = append(out, fmt.Sprintf("`%s`.%s", t.Name, column.Name))
	}
	return
}
