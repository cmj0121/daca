package args

import (
	"errors"
	"fmt"
	"strings"
)

type Flag struct {
	name     string
	shortcut byte
	parser   *argparser

	value         interface{}
	default_value interface{}
	choices       []string
	required      bool
	validator     func(interface{}) error
	callback      func(p *argparser, opt option) error
	help          string
}

/* getter */
func (opt *Flag) Get() (out interface{}) {
	out = opt.value
	if out == nil {
		out = opt.default_value
	}
	return
}

func (opt *Flag) Set(in []string, idx *int) (err error) {
	var value string

	switch {
	case in[*idx][:2] == "--" && in[*idx][2:len(opt.name)+2] == opt.name:
		if len(in[*idx]) > len(opt.name)+2 && in[*idx][len(opt.name)+3] == '=' {
			/* The format : --key=value */
			value = in[*idx][len(opt.name)+3:]
		}
	case in[*idx][0] == '-' && in[*idx][1] == opt.shortcut:
		value = in[*idx][2:]
	}

	if value == "" {
		/* value is the next args */
		if *idx+1 == len(in) {
			err = errors.New(fmt.Sprintf("options -%c, --%s need value", opt.shortcut, opt.name))
			return
		}

		*idx++
		value = in[*idx]
	}

	if opt.validator != nil {
		if err = opt.validator(value); err != nil {
			return
		}
	}

	/* choices validator */
	if err = ChoiceValidator(value, opt.choices); err != nil {
		return
	}

	if opt.callback != nil {
		err = opt.callback(opt.parser, opt)
	}

	opt.value = value
	return
}

func (opt *Flag) String(format string) (out string) {
	if opt.shortcut == 0 {
		out = fmt.Sprintf("--%s", opt.name)
	} else {
		out = fmt.Sprintf("-%c, --%s", opt.shortcut, opt.name)
	}

	out = fmt.Sprintf(format, out, opt.help)

	if opt.default_value != nil {
		out = fmt.Sprintf("%s (default: %v)", out, opt.default_value)
	}

	if opt.choices != nil {
		out = fmt.Sprintf("%s %v", out, opt.choices)
	}
	return
}

func (opt *Flag) Hint() (out string) {
	if opt.shortcut == 0 {
		out = fmt.Sprintf("--%s %s", opt.name, strings.ToUpper(opt.name))
	} else {
		out = fmt.Sprintf("-%c %s", opt.shortcut, strings.ToUpper(opt.name))
	}

	if opt.required == false {
		out = fmt.Sprintf("[%s]", out)
	}
	return
}

func (opt *Flag) Check() (err error) {
	if opt.required && opt.value == nil {
		err = errors.New(fmt.Sprintf("options '%s' is required", opt.name))
	}
	return
}

/* setter */
func (opt *Flag) Shortcut(in byte) (out option) {
	opt.shortcut = in
	out = opt.parser.Shortcut(in, opt)
	return
}

func (opt *Flag) Default(in interface{}) (out option) {
	opt.default_value = in.(string)
	out = opt
	return
}

func (opt *Flag) Choice(in []string) (out option) {
	opt.choices = in
	out = opt
	return
}

func (opt *Flag) Required(in bool) (out option) {
	opt.required = in
	out = opt
	return
}

func (opt *Flag) Validator(in func(interface{}) error) (out option) {
	opt.validator = in
	out = opt
	return
}

func (opt *Flag) Callback(in func(p *argparser, opt option) error) (out option) {
	opt.callback = in
	out = opt
	return
}

func (opt *Flag) Help(in string) (out option) {
	opt.help = in
	out = opt
	return
}
