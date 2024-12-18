package json_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/test"
)

type item struct {
	Foo string `json:"foo"`
}

var opt = test.EncodingOptions[item]{
	Mime:   "application/json",
	Raw:    `{"foo":"bar"}` + "\n",
	Parsed: item{Foo: "bar"},
}

func TestJSON(t *testing.T) {
	test.Encoding(t, opt)
}

func BenchmarkJSON(b *testing.B) {
	test.BenchmarkEncoding(b, opt)
}
