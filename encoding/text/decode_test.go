package text_test

import (
	"errors"
	"io"
	"reflect"
	"strconv"
	"testing"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/encoding"
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/google/go-cmp/cmp"
	"github.com/valyala/fasthttp"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		in   string
		want any
	}{
		{"  \n", ""},
		{"test\n", "test"},
		{"test\n", []byte("test")},
		{"5\n", int(5)},
		{"5\n", int8(5)},
		{"5\n", int16(5)},
		{"5\n", int32(5)},
		{"5\n", int64(5)},
		{"5\n", uint(5)},
		{"5\n", uint8(5)},
		{"5\n", uint16(5)},
		{"5\n", uint32(5)},
		{"5\n", uint64(5)},
		{"5.1\n", float32(5.1)},
		{"5.1\n", float64(5.1)},
		{"true\n", true},
		{"test\n", unmarshaler{S: "test"}},
		{"test\n", unmarshaler{S: "test"}}, // Test cached unmarshaler.
		{"test\n", &unmarshaler{S: "test"}},
	}

	dec := encoding.GetDecoder([]byte("text/plain"))
	if dec == nil {
		t.Fatal("decoder not found")
	}

	for _, test := range tests {
		ctx := httptest.NewRequest(fasthttp.MethodGet, "/", test.in, nil)
		v := reflect.New(reflect.TypeOf(test.want)).Interface()
		if err := dec(ctx, v); err != nil {
			t.Error(err)
		} else {
			if diff := cmp.Diff(test.want, reflect.ValueOf(v).Elem().Interface()); diff != "" {
				t.Errorf("%T: %s", test.want, diff)
			}
		}
	}
}

func TestDecodeError(t *testing.T) {
	tests := []struct {
		in   string
		val  any
		want error
	}{
		{"test\n", &struct{}{}, don.ErrUnsupportedMediaType},
		{"test\n", ptr(0), strconv.ErrSyntax},
		{"test\n", &unmarshaler{Err: io.EOF}, io.EOF},
	}

	dec := encoding.GetDecoder([]byte("text/plain"))
	if dec == nil {
		t.Fatal("decoder not found")
	}

	for _, test := range tests {
		ctx := httptest.NewRequest(fasthttp.MethodGet, "/", test.in, nil)
		if err := dec(ctx, test.val); err == nil {
			t.Error("should return error")
		} else if !errors.Is(err, test.want) {
			t.Errorf("should return error %q, got %q", test.want, err)
		}
	}
}

func ptr[T any](v T) *T {
	return &v
}

type unmarshaler struct {
	S   string
	Err error
}

func (m *unmarshaler) UnmarshalText(text []byte) error {
	if m.Err != nil {
		return m.Err
	}
	m.S = string(text)
	return nil
}
