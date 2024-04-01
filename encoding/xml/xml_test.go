package xml_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/test"
)

type item struct {
	Foo string `xml:"foo"`
}

var opt = test.EncodingOptions[item]{
	Mime:   "application/xml",
	Raw:    "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<item><foo>bar</foo></item>",
	Parsed: item{Foo: "bar"},
}

func TestXML(t *testing.T) {
	test.Encoding(t, opt)
}

func BenchmarkXML(b *testing.B) {
	test.BenchmarkEncoding(b, opt)
}
