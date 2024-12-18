package test

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/encoding"
	_ "github.com/abemedia/go-don/encoding/text" // default encoding
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type EncodingOptions[T any] struct {
	Mime   string
	Parsed T
	Raw    string
}

func Encoding[T any](t *testing.T, opt EncodingOptions[T]) {
	t.Helper()
	t.Run("Decode", func(t *testing.T) {
		t.Helper()
		Decode(t, opt)
	})
	t.Run("Encode", func(t *testing.T) {
		t.Helper()
		Encode(t, opt)
	})
}

func Decode[T any](t *testing.T, opt EncodingOptions[T]) {
	t.Helper()

	var diff string

	api := don.New(nil)
	api.Post("/", don.H(func(ctx context.Context, req T) (any, error) {
		diff = cmp.Diff(opt.Parsed, req, ignoreUnexported[T]())
		return nil, nil
	}))

	ctx := httptest.NewRequest(http.MethodPost, "/", opt.Raw, map[string]string{"Content-Type": opt.Mime})
	api.RequestHandler()(ctx)

	if diff != "" {
		t.Error(diff)
	}

	if ctx.Response.StatusCode() != http.StatusNoContent {
		t.Errorf("expected success status: %v", &ctx.Response)
	}
}

func Encode[T any](t *testing.T, opt EncodingOptions[T]) {
	t.Helper()

	api := don.New(nil)
	api.Post("/", don.H(func(ctx context.Context, req any) (T, error) {
		return opt.Parsed, nil
	}))

	ctx := httptest.NewRequest(http.MethodPost, "/", "", map[string]string{"Accept": opt.Mime})
	api.RequestHandler()(ctx)

	if diff := cmp.Diff(opt.Raw, string(ctx.Response.Body()), ignoreUnexported[T]()); diff != "" {
		t.Error(diff)
	}

	if ctx.Response.StatusCode() != http.StatusOK {
		t.Errorf("expected success status: %v", &ctx.Response)
	}
}

func ignoreUnexported[T any]() cmp.Option {
	t := reflect.TypeOf(*new(T))
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	return cmpopts.IgnoreUnexported(reflect.New(t).Elem().Interface())
}

func BenchmarkEncoding[T any](b *testing.B, opt EncodingOptions[T]) {
	b.Run("Decode", func(b *testing.B) {
		b.Helper()
		BenchmarkDecode(b, opt)
	})
	b.Run("Encode", func(b *testing.B) {
		b.Helper()
		BenchmarkEncode(b, opt)
	})
}

func BenchmarkDecode[T any](b *testing.B, opt EncodingOptions[T]) {
	b.Helper()

	dec := encoding.GetDecoder([]byte(opt.Mime))
	if dec == nil {
		b.Fatal("decoder not found")
	}

	rd := strings.NewReader(opt.Raw)
	ctx := httptest.NewRequest("POST", "/", "", nil)
	ctx.Request.SetBodyStream(rd, len(opt.Raw))

	v := new(T)
	if val := reflect.ValueOf(v).Elem(); val.Kind() == reflect.Pointer {
		val.Set(reflect.New(val.Type().Elem()))
	}

	for i := 0; i < b.N; i++ {
		rd.Seek(0, io.SeekStart) //nolint:errcheck
		dec(ctx, v)              //nolint:errcheck
	}
}

func BenchmarkEncode[T any](b *testing.B, opt EncodingOptions[T]) {
	b.Helper()

	enc := encoding.GetEncoder([]byte(opt.Mime))
	if enc == nil {
		b.Fatal("encoder not found")
	}

	ctx := httptest.NewRequest("POST", "/", "", nil)

	for i := 0; i < b.N; i++ {
		ctx.Response.ResetBody()
		enc(ctx, opt.Parsed) //nolint:errcheck
	}
}
