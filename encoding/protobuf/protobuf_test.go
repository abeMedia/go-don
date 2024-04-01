package protobuf_test

import (
	"errors"
	"testing"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/encoding"
	"github.com/abemedia/go-don/encoding/protobuf/testdata"
	"github.com/abemedia/go-don/internal/test"
	"github.com/abemedia/go-don/pkg/httptest"
)

//go:generate protoc testdata/test.proto --go_out=. --go_opt=paths=source_relative

var opt = test.EncodingOptions[*testdata.Item]{
	Mime:   "application/protobuf",
	Raw:    "\n\x03bar",
	Parsed: &testdata.Item{Foo: "bar"},
}

func TestProtobuf(t *testing.T) {
	test.Encoding(t, opt)
}

func TestProtobufError(t *testing.T) {
	ctx := httptest.NewRequest("", "", "", nil)
	v := "test"
	mime := []byte("application/protobuf")

	t.Run("Decode", func(t *testing.T) {
		dec := encoding.GetDecoder(mime)
		err := dec(ctx, &v)
		if !errors.Is(err, don.ErrUnsupportedMediaType) {
			t.Fatal("should fail")
		}
	})

	t.Run("Encode", func(t *testing.T) {
		enc := encoding.GetEncoder(mime)
		err := enc(ctx, &v)
		if !errors.Is(err, don.ErrNotAcceptable) {
			t.Fatal("should fail")
		}
	})
}

func BenchmarkProtobuf(b *testing.B) {
	test.BenchmarkEncoding(b, opt)
}
