/* Copyright (C) 2019-2019 cmj. All right reserved. */
package daca

import (
	"fmt"
	"regexp"
	"strings"
)

type Endpoint struct {
	Handler /* the handler of this endpoints */

	url        string            /* raw URL */
	methods    []string          /* support HTTP methods */
	pattern    *regexp.Regexp    /* URL pattern */
	named_maps map[int]string    /* map of argument - index: name */
	named_args map[string]string /* map of process result - name: value */
}

func NewEndpoint(url string, fn Handler, methods ...string) (out *Endpoint) {
	out = &Endpoint{
		Handler:    fn,
		url:        url,
		methods:    methods,
		named_maps: make(map[int]string, 0),
		named_args: make(map[string]string),
	}

	/* fetch all sub-pattern in URL */
	re := regexp.MustCompile(`<(.*?)(:.*?)?>`)
	for idx, pattern := range re.FindAllStringSubmatch(url, -1) {
		args_name := pattern[1]
		args_type := pattern[2]

		if _, ok := out.named_args[args_name]; ok {
			panic(fmt.Sprintf("duplicated named argument %s in URL - %v", args_name, url))
			return
		}

		out.named_maps[idx] = args_name

		/* replace the old URL to new URL */
		sub_re := regexp.MustCompile(string(pattern[0]))
		switch args_type {
		case "", ":path": /* all passible file path */
			url = sub_re.ReplaceAllString(url, `([\w\.\-]*)`)
		case ":int":
			url = sub_re.ReplaceAllString(url, `(\d+)`)
		case ":str", ":string":
			url = sub_re.ReplaceAllString(url, `([\w\-]+)`)
		default:
			panic(fmt.Sprintf("Not implement type - %v", args_type))
		}
	}

	out.Regexp(fmt.Sprintf("^%s$", url))
	return
}

func (e *Endpoint) Match(in string) (out bool) {
	out = e.pattern.MatchString(in)

	/* process the match result and get the named argument */
	if out && len(e.named_maps) > 0 {
		for idx, value := range e.pattern.FindAllStringSubmatch(in, -1) {
			e.named_args[e.named_maps[idx]] = value[1]
		}
	}

	return
}

func (e *Endpoint) String() (out string) {
	out = fmt.Sprintf("[%s] %-20s %s", strings.Join(e.methods, ", "), e.url, e.pattern)
	return
}

func (e *Endpoint) Regexp(in string) (out *Endpoint) {
	e.pattern = regexp.MustCompile(in)
	out = e
	return
}

func (e *Endpoint) Get(in string) (out string) {
	out = e.named_args[in]
	return
}
