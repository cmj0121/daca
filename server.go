/* Copyright (C) 2019-2019 cmj. All right reserved. */
package daca

import (
	"github.com/cmj0121/daca/logger"
	"net/http"
	"strings"
)

type Server struct {
	*logger.Logger /* logger sub-system */
	bind         string         /* bind address */
	default_page Handler        /* default page */
	endpoints    map[string][]*Endpoint
}

func NewServer(bind string) (srv *Server) {
	srv = &Server{
		Logger:    logger.DefaultLogger(),
		bind:      bind,
		endpoints: make(map[string][]*Endpoint, 0),
	}

	return
}

func (srv *Server) Run() {
	srv.Info("Run on %v", srv.bind)
	srv.Fatal("%v", http.ListenAndServe(srv.bind, srv))
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Server:  srv,
		writer:  w,
		request: r,
	}

	method := strings.ToUpper(ctx.request.Method)
	url := ctx.request.URL.Path

	for _, endpoint := range srv.endpoints[method] {
		if endpoint.Match(url) {
			/* Save the endpoint to context */
			ctx.Endpoint = endpoint
			srv.Info("[%s] %s", method, url)
			endpoint.Handler(ctx)
			return
		}
	}

	/* default page */
	if srv.default_page != nil {
		srv.Info("[%s] %s as default page", method, url)
		srv.default_page(ctx)
		return
	}
	/* 404 Not Found */
	http.NotFound(w, r)
}

func (srv *Server) DefaultPage(fn Handler) (out *Server) {
	srv.default_page = fn
	out = srv
	return
}

func (srv *Server) Route(url string, fn Handler, methods ...string) (endpoint *Endpoint) {
	endpoint = NewEndpoint(url, fn, methods...)

	for _, method := range methods {
		method = strings.ToUpper(method)

		if _, ok := srv.endpoints[method]; ok == false {
			srv.endpoints[method] = make([]*Endpoint, 0)
		}

		srv.endpoints[method] = append(srv.endpoints[method], endpoint)
		srv.Fatal("Add %v", endpoint)
	}

	return
}
