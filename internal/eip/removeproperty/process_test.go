package removeproperty

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestRemoveProperty(t *testing.T) {
	p := NewProcessor("test", "test", "x", "y", "z")

	e := exchange.NewExchange(nil)
	e.SetProperty("x", 1)
	e.SetProperty("y", 2)
	e.SetProperty("z", 3)
	e.SetProperty("q", "Hello")

	p.Process(e)

	if e.HasProperty("x") != false {
		t.Fatalf("TestRemoveProperty() = removed exchange property 'x' is still present")
	}

	if e.HasProperty("y") != false {
		t.Fatalf("TestRemoveProperty() = removed exchange property 'y' is still present")
	}

	if e.HasProperty("q") != true {
		t.Fatalf("TestRemoveProperty() = missing exchange property 'q'")
	}
}
