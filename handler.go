package don

import (
	"context"
	"net/http"

	"github.com/abemedia/go-don/encoding"
	"github.com/abemedia/httprouter"
	"github.com/valyala/fasthttp"
)

// StatusCoder allows you to customise the HTTP response code.
type StatusCoder interface {
	StatusCode() int
}

// Headerer allows you to customise the HTTP headers.
type Headerer interface {
	Header() http.Header
}

// Handle is the type for your handlers.
type Handle[T, O any] func(ctx context.Context, request T) (O, error)

// H wraps your handler function with the Go generics magic.
func H[T, O any](handle Handle[T, O]) httprouter.Handle {
	pool := newRequestPool(*new(T))
	decodeRequest := newRequestDecoder(*new(T))
	isNil := newNilCheck(*new(O))

	return func(ctx *fasthttp.RequestCtx, p httprouter.Params) {
		enc := encoding.GetEncoder(ctx.Request.Header.Peek(fasthttp.HeaderAccept))
		if enc == nil {
			handleError(ctx, ErrNotAcceptable)
			return
		}

		var res any

		req := pool.Get()
		err := decodeRequest(req, ctx, p)
		if err != nil {
			res = Error(err, getStatusCode(err, http.StatusBadRequest))
		} else {
			res, err = handle(ctx, *req)
			if err != nil {
				res = Error(err, 0)
			}
		}
		pool.Put(req)

		if h, ok := res.(Headerer); ok {
			for k, v := range h.Header() {
				ctx.Response.Header.Set(k, v[0])
			}
		}

		if sc, ok := res.(StatusCoder); ok {
			ctx.SetStatusCode(sc.StatusCode())
		}

		if err == nil && isNil(res) {
			res = nil
			ctx.Response.Header.SetContentLength(-3)
		}

		if err = enc(ctx, res); err != nil {
			handleError(ctx, err)
		}
	}
}

func handleError(ctx *fasthttp.RequestCtx, err error) {
	code := getStatusCode(err, http.StatusInternalServerError)
	if code < http.StatusInternalServerError {
		ctx.Error(err.Error(), code)
		return
	}
	ctx.Error(fasthttp.StatusMessage(code), code)
	ctx.Logger().Printf("%v", err)
}

func getStatusCode(i any, fallback int) int {
	if sc, ok := i.(StatusCoder); ok {
		return sc.StatusCode()
	}
	return fallback
}
