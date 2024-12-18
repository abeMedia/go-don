package encoding_test

import (
	"io"
	"testing"

	"github.com/abemedia/go-don/encoding"
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/valyala/fasthttp"
)

func TestRegisterEncoder(t *testing.T) {
	enc := func(ctx *fasthttp.RequestCtx, v any) error {
		b := v.([]byte)
		if len(b) == 0 {
			return io.EOF
		}
		ctx.Response.SetBodyRaw(b)
		return nil
	}

	encoding.RegisterEncoder(enc, "response-encoder", "response-encoder-alias")

	for _, v := range []string{"response-encoder", "response-encoder-alias"} {
		encode := encoding.GetEncoder([]byte(v))
		if encode == nil {
			t.Error("encoder not found")
			continue
		}

		req := httptest.NewRequest("", "", v, nil)

		if err := encode(req, []byte(v)); err != nil {
			t.Error(err)
		} else if string(req.Response.Body()) != v {
			t.Error("should encode response")
		}

		if err := encode(req, []byte{}); err == nil {
			t.Error("should return error")
		}
	}
}

func TestGetEncoder_MultipleContentTypes(t *testing.T) {
	encFn := func(ctx *fasthttp.RequestCtx, v any) error {
		return nil
	}

	encoding.RegisterEncoder(encFn, "application/xml")

	enc := encoding.GetEncoder([]byte("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"))
	if enc == nil {
		t.Fatal("encoder not found")
	}

	enc = encoding.GetEncoder([]byte("application/xhtml+xml"))
	if enc != nil {
		t.Fatal("encoder should not be found")
	}
}
