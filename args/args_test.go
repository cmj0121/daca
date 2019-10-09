package args

import (
	"testing"
)

func TestCreateParser(t *testing.T) {
	p := Parser("simple parser").Version("1.0.0")

	p.String()
}

func TestCreateFlip(t *testing.T) {
	p := Parser("simple parser").Version("1.0.0")

	p.Flip("enable").Default(false).Help("Enable the option")
	p.Flip("show").Help("Show")
}

func TestCreateFlag(t *testing.T) {
	p := Parser("simple parser").Version("1.0.0")

	p.Flag("enable").Default("yes").Help("Enable the option")
	p.Flag("show").Help("Show")
}

func TestCreateParameter(t *testing.T) {
	p := Parser("simple parser").Version("1.0.0")

	p.Parameter("action").Default("yes").Choice([]string{"yes", "no"})
}

/* failure detection */
func TestDuplicateFlip(t *testing.T) {
	p := Parser("simple parser").Version("1.0.0")
	defer func() {
		if r := recover(); r == nil {
			t.Error("Does not detect duplicate Flip name")
		}
	}()

	p.Flip("interface")
	p.Flip("interface")
}

func TestDuplicateFlag(t *testing.T) {
	p := Parser("simple parser").Version("1.0.0")
	defer func() {
		if r := recover(); r == nil {
			t.Error("Does not detect duplicate Flag name")
		}
	}()

	p.Flag("interface")
	p.Flag("interface")
}
