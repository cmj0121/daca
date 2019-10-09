package args

import (
	"errors"
	"fmt"
)

type Flip struct {
	name     string
	shortcut byte
	parser   *argparser

	value         interface{}
	default_value interface{}
	callback      func(p *argparser, opt option) error
	required      bool
	help          string
}

/* getter */
func (opt *Flip) Get() (out interface{}) {
	out = opt.value
	if out == nil {
		out = opt.default_value
	}
	return
}

func (opt *Flip) Set(in []string, idx *int) (err error) {
	switch {
	case in[*idx][:2] == "--":
		if value := in[*idx][2+len(opt.name):]; value != "" {
			err = errors.New(fmt.Sprintf("flip '--%s' cannot add extra value", opt.name))
			return
		}
	case in[*idx][0] == '-' && in[*idx][1] == opt.shortcut:
		if value := in[*idx][2:]; value != "" {
			err = errors.New(fmt.Sprintf("shortcut '-%c' cannot add extra value", opt.shortcut))
			return
		}
	}

	if opt.value == nil {
		/* NOTE - Flip only store the true/false */
		opt.value = true
	} else {
		opt.value = !opt.value.(bool)
	}

	if opt.callback != nil {
		err = opt.callback(opt.parser, opt)
	}

	return
}

func (opt *Flip) String(format string) (out string) {
	if opt.shortcut == 0 {
		out = fmt.Sprintf("--%s", opt.name)
	} else {
		out = fmt.Sprintf("-%c, --%s", opt.shortcut, opt.name)
	}

	out = fmt.Sprintf(format, out, opt.help)
	if opt.default_value != nil {
		out = fmt.Sprintf("%s (default: %v)", out, opt.default_value)
	}
	return
}

func (opt *Flip) Hint() (out string) {
	if opt.shortcut == 0 {
		out = fmt.Sprintf("--%s", opt.name)
	} else {
		out = fmt.Sprintf("-%c", opt.shortcut)
	}

	if opt.required == false {
		out = fmt.Sprintf("[%s]", out)
	}
	return
}

func (opt *Flip) Check() (err error) {
	if opt.required && opt.value == nil {
		err = errors.New(fmt.Sprintf("options '%s' is required", opt.name))
	}
	return
}

/* setter */
func (opt *Flip) Shortcut(in byte) (out option) {
	opt.shortcut = in
	out = opt.parser.Shortcut(in, opt)
	return
}

func (opt *Flip) Default(in interface{}) (out option) {
	opt.default_value = in.(bool)
	out = opt
	return
}

func (opt *Flip) Choice(in []string) (out option) {
	panic("flip does NOT support Choice")
	out = opt
	return
}

func (opt *Flip) Required(in bool) (out option) {
	opt.required = in
	out = opt
	return
}

func (opt *Flip) Validator(in func(interface{}) error) (out option) {
	out = opt
	return
}

func (opt *Flip) Callback(in func(p *argparser, opt option) error) (out option) {
	opt.callback = in
	out = opt
	return
}

func (opt *Flip) Help(in string) (out option) {
	opt.help = in
	out = opt
	return
}
