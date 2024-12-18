package encoding_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/abemedia/go-don/encoding"
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/valyala/fasthttp"
)

func TestRegisterDecoder(t *testing.T) {
	dec := func(ctx *fasthttp.RequestCtx, v any) error {
		b := ctx.Request.Body()
		if len(b) == 0 {
			return io.EOF
		}
		reflect.ValueOf(v).Elem().SetBytes(b)
		return nil
	}

	encoding.RegisterDecoder(dec, "mime", "alias")

	for _, v := range []string{"mime", "alias", "mime; charset=utf-8", "alias; charset=utf-8"} {
		decode := encoding.GetDecoder([]byte(v))
		if decode == nil {
			t.Error("decoder not found")
			continue
		}

		req := httptest.NewRequest("", "", v, nil)

		var b []byte
		if err := decode(req, &b); err != nil {
			t.Error(err)
		} else if string(b) != v {
			t.Error("should decode request")
		}

		req.Request.SetBody(nil)
		if err := decode(req, &b); err == nil {
			t.Error("should return error")
		}
	}
}
