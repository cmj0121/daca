package main

import (
	"github.com/cmj0121/daca"
)

func foo(ctx *daca.Context) {
	key := ctx.Get("key")
	ctx.Write([]byte(key))
}

func main() {
	srv := daca.NewServer(":8888")
	srv.Route("/<key:path>", foo, "GET")
	srv.Run()
}
