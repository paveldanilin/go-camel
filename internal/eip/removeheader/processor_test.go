package removeheader

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestRemoveHeaderProcessor(t *testing.T) {
	p := NewProcessor("", "", "xyz", "qwerty")

	e := exchange.NewExchange(nil)
	e.Message().SetHeader("xyz", 1)
	e.Message().SetHeader("z", "y")
	e.Message().SetHeader("qwerty", true)

	p.Process(e)

	if e.Message().HasHeader("xyz") != false {
		t.Fatalf("TestRemoveHeaderProcessor() = removed message header 'xyz' is still present")
	}

	if e.Message().HasHeader("qwerty") != false {
		t.Fatalf("TestRemoveHeaderProcessor() = removed message header 'qwerty' is still present")
	}

	if e.Message().HasHeader("z") != true {
		t.Fatalf("TestRemoveHeaderProcessor() = missing message header 'z'")
	}
}
