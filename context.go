/* Copyright (C) 2019-2019 cmj. All right reserved. */
package daca

import (
	"bytes"
	"encoding/json"
	"net/http"
	"text/template"
	"strings"
)

type Context struct {
	*Endpoint
	*Server
	writer  http.ResponseWriter
	request *http.Request
}

func (ctx *Context) Write(in []byte) {
	/* reply the Accept as possible */
	accept := strings.Split(ctx.request.Header.Get("Accept"), ",")
	if len(accept) > 1 && accept[0] != "*/*" {
		ctx.Header("Content-Type", accept[0])
	}

	/* Set HTML body */
	ctx.writer.Write(in)
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
	ctx.writer.Header().Set(key, value)
}

func (ctx *Context) Query(key string) (out string) {
	vars := ctx.request.URL.Query()
	out = vars.Get(key)
	return
}
