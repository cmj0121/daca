/* Copyright (C) 2019-2019 cmj. All right reserved. */
package daca

import (
	"bytes"
	"encoding/json"
	"net/http"
	"html/template"
)

type Context struct {
	*Endpoint
	*Server
	W  http.ResponseWriter
	R *http.Request
}

func (ctx *Context) Write(in []byte) {
	/* Set HTML body */
	ctx.W.Write(in)
}

func (ctx *Context) WriteJson(in interface{}) {
	/* Set MIME as application/json */
	ctx.Header("Content-Type", "application/json")
	/* Write the JSON payload */
	data, _ := json.Marshal(in)
	ctx.Write(data)
}

func (ctx *Context) WriteTemplate(tmpl string, in interface{}) {
	var buff bytes.Buffer

	t := template.Must(template.New("tmpl").Parse(tmpl))
	t.Execute(&buff, in)
	ctx.Write(buff.Bytes())
}

func (ctx *Context) WriteTemplateFile(path string, in interface{}) {
	var buff bytes.Buffer

	t := template.Must(template.ParseFiles(path))
	t.Execute(&buff, in)
	ctx.Write(buff.Bytes())
}

func (ctx *Context) Header(key, value string) {
	/* Set HTML header */
	ctx.W.Header().Set(key, value)
}

func (ctx *Context) Query(key string) (out string) {
	vars := ctx.R.URL.Query()
	out = vars.Get(key)
	return
}
