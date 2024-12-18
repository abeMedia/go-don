package text_test

import (
	"errors"
	"io"
	"testing"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/encoding"
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/google/go-cmp/cmp"
	"github.com/valyala/fasthttp"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		in   any
		want string
	}{
		{"test", "test"},
		{[]byte("test"), "test"},
		{int(5), "5"},
		{int8(5), "5"},
		{int16(5), "5"},
		{int32(5), "5"},
		{int64(5), "5"},
		{uint(5), "5"},
		{uint8(5), "5"},
		{uint16(5), "5"},
		{uint32(5), "5"},
		{uint64(5), "5"},
		{float32(5.1), "5.1"},
		{float64(5.1), "5.1"},
		{true, "true"},
		{errors.New("test"), "test"},
		{marshaler{s: "test"}, "test"},
		{&marshaler{s: "test"}, "test"},
		{stringer{}, "test"},
		{&stringer{}, "test"},
	}

	enc := encoding.GetEncoder([]byte("text/plain"))
	if enc == nil {
		t.Fatal("encoder not found")
	}

	for _, test := range tests {
		ctx := httptest.NewRequest(fasthttp.MethodGet, "/", "", nil)
		if err := enc(ctx, test.in); err != nil {
			t.Error(err)
		} else {
			if diff := cmp.Diff(test.want, string(ctx.Response.Body())); diff != "" {
				t.Errorf("%T: %s", test.in, diff)
			}
		}
	}
}

func TestEncodeError(t *testing.T) {
	tests := []struct {
		in   any
		want error
	}{
		{&struct{}{}, don.ErrNotAcceptable},
		{marshaler{err: io.EOF}, io.EOF},
	}

	enc := encoding.GetEncoder([]byte("text/plain"))
	if enc == nil {
		t.Fatal("encoder not found")
	}

	for _, test := range tests {
		ctx := httptest.NewRequest(fasthttp.MethodGet, "/", "", nil)
		if err := enc(ctx, test.in); err == nil {
			t.Error("should return error")
		} else if !errors.Is(err, test.want) {
			t.Errorf("should return error %q, got %q", test.want, err)
		}
	}
}

type marshaler struct {
	s   string
	err error
}

func (m marshaler) MarshalText() ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []byte(m.s), nil
}

type stringer struct{}

func (m stringer) String() string {
	return "test"
}
