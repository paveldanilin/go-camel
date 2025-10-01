package test

import (
	uri2 "github.com/paveldanilin/go-camel/camel"
	"testing"
)

func TestParseCamelStyle(t *testing.T) {
	uri, err := uri2.Parse("timer:foo?period=1000", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedComponent := "timer"
	if uri.Component() != expectedComponent {
		t.Errorf("Expected component '%s', but got '%s'", expectedComponent, uri.Component())
	}

	expectedPath := "foo"
	if uri.Path() != expectedPath {
		t.Errorf("Expected path '%s', but got '%s'", expectedPath, uri.Path())
	}

	expectedPeriod := 1000
	if uri.MustParamInt("period") != expectedPeriod {
		t.Errorf("Expected period '%d', but got '%d'", expectedPeriod, uri.MustParamInt("period"))
	}
}

func TestParseCamelStyleMultiParams(t *testing.T) {
	uri, err := uri2.Parse("kafka:topicA?brokers=localhost:9092&acks=all", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedComponent := "kafka"
	if uri.Component() != expectedComponent {
		t.Errorf("Expected component '%s', but got '%s'", expectedComponent, uri.Component())
	}

	expectedPath := "topicA"
	if uri.Path() != expectedPath {
		t.Errorf("Expected path '%s', but got '%s'", expectedPath, uri.Path())
	}

	expectedParams := []string{"brokers", "acks"}
	if !uri.HasParams(expectedParams...) {
		t.Errorf("Expected params not found: %v, provided: %v", expectedParams, uri.Params())
	}
}

func TestParseCamelStyleFile(t *testing.T) {
	uri, err := uri2.Parse("file:/var/log?recursive=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedComponent := "file"
	if uri.Component() != expectedComponent {
		t.Errorf("Expected component '%s', but got '%s'", expectedComponent, uri.Component())
	}

	expectedPath := "/var/log"
	if uri.Path() != expectedPath {
		t.Errorf("Expected path '%s', but got '%s'", expectedPath, uri.Path())
	}

	expectedRecursive := true
	if uri.MustParamBool("recursive") != expectedRecursive {
		t.Errorf("Expected recursive %v, but got %v", expectedRecursive, uri.MustParamBool("recursive"))
	}
}

func TestParseRegularURL(t *testing.T) {
	uri, err := uri2.Parse("http://john:smith@222.111.222.111:8080/a/b?x=1#frag", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedComponent := "http"
	if uri.Component() != expectedComponent {
		t.Errorf("Expected component '%s', but got '%s'", expectedComponent, uri.Component())
	}

	expectedPath := "/a/b"
	if uri.Path() != expectedPath {
		t.Errorf("Expected path '%s', but got '%s'", expectedPath, uri.Path())
	}

	expectedHost := "222.111.222.111"
	if uri.Host() != expectedHost {
		t.Errorf("Expected host '%s', but got '%s", expectedHost, uri.Host())
	}

	expectedPort := "8080"
	if uri.Host() != expectedHost {
		t.Errorf("Expected port '%s', but got '%s", expectedPort, uri.Port())
	}

	expectedUser := "john"
	if uri.Username() != expectedUser {
		t.Errorf("Expected username '%s', but got '%s", expectedUser, uri.Username())
	}

	expectedPassword := "smith"
	if uri.Password() != expectedPassword {
		t.Errorf("Expected password '%s', but got '%s", expectedPassword, uri.Password())
	}

	expectedParam := "1"
	if uri.MustParam("x") != expectedParam {
		t.Errorf("Expected param x '%s', but got '%s", expectedParam, uri.MustParam("x"))
	}
}
