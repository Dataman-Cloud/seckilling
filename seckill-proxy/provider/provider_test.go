package provider

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"text/template"
)

type myProvider struct {
	baseProvider
}

func (p *myProvider) Foo() string {
	return "bar"
}

func TestConfigurationErrors(t *testing.T) {
	templateErrorFile, err := ioutil.TempFile("", "provider-configuration-error")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(templateErrorFile.Name())
	data := []byte("Not a valid template {{ Bar }}")
	err = ioutil.WriteFile(templateErrorFile.Name(), data, 0700)
	if err != nil {
		t.Fatal(err)
	}

	templateInvalidTOMLFile, err := ioutil.TempFile("", "provider-configuration-error")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(templateInvalidTOMLFile.Name())
	data = []byte(`Hello {{ .Name }}
{{ Foo }}`)
	err = ioutil.WriteFile(templateInvalidTOMLFile.Name(), data, 0700)
	if err != nil {
		t.Fatal(err)
	}

	invalids := []struct {
		provider        *myProvider
		defaultTemplate string
		expectedError   string
		funcMap         template.FuncMap
		templateObjects interface{}
	}{
		{
			provider: &myProvider{
				baseProvider{
					Filename: "/non/existent/template.tmpl",
				},
			},
			expectedError: "open /non/existent/template.tmpl: no such file or directory",
		},
		{
			provider:        &myProvider{},
			defaultTemplate: "non/existent/template.tmpl",
			expectedError:   "Asset non/existent/template.tmpl not found",
		},
		{
			provider: &myProvider{
				baseProvider{
					Filename: templateErrorFile.Name(),
				},
			},
			expectedError: `function "Bar" not defined`,
		},
		{
			provider: &myProvider{
				baseProvider{
					Filename: templateInvalidTOMLFile.Name(),
				},
			},
			expectedError: "Near line 1, key 'Hello': Near line 1: Expected key separator '=', but got '<' instead",
			funcMap: template.FuncMap{
				"Foo": func() string {
					return "bar"
				},
			},
			templateObjects: struct{ Name string }{Name: "bar"},
		},
	}

	for _, invalid := range invalids {
		configuration, err := invalid.provider.getConfiguration(invalid.defaultTemplate, invalid.funcMap, nil)
		if err == nil || !strings.Contains(err.Error(), invalid.expectedError) {
			t.Fatalf("should have generate an error with %q, got %v", invalid.expectedError, err)
		}
		if configuration != nil {
			t.Fatalf("shouldn't have return a configuration object : %v", configuration)
		}
	}
}

func TestGetConfiguration(t *testing.T) {
	templateFile, err := ioutil.TempFile("", "provider-configuration")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(templateFile.Name())
	data := []byte(`[backends]
  [backends.backend1]
    [backends.backend1.circuitbreaker]
      expression = "NetworkErrorRatio() > 0.5"
    [backends.backend1.servers.server1]
    url = "http://172.17.0.2:80"
    weight = 10
    [backends.backend1.servers.server2]
    url = "http://172.17.0.3:80"
    weight = 1

[frontends]
  [frontends.frontend1]
  backend = "backend1"
  passHostHeader = true
    [frontends.frontend11.routes.test_2]
    rule = "Path"
    value = "/test"`)
	err = ioutil.WriteFile(templateFile.Name(), data, 0700)
	if err != nil {
		t.Fatal(err)
	}

	provider := &myProvider{
		baseProvider{
			Filename: templateFile.Name(),
		},
	}
	configuration, err := provider.getConfiguration(templateFile.Name(), nil, nil)
	if err != nil {
		t.Fatalf("Shouldn't have error out, got %v", err)
	}
	if configuration == nil {
		t.Fatalf("Configuration should not be nil, but was")
	}
}

func TestReplace(t *testing.T) {
	cases := []struct {
		str      string
		expected string
	}{
		{
			str:      "",
			expected: "",
		},
		{
			str:      "foo",
			expected: "bar",
		},
		{
			str:      "foo foo",
			expected: "bar bar",
		},
		{
			str:      "somethingfoo",
			expected: "somethingbar",
		},
	}

	for _, c := range cases {
		actual := replace("foo", "bar", c.str)
		if actual != c.expected {
			t.Fatalf("expected %q, got %q, for %q", c.expected, actual, c.str)
		}
	}
}
