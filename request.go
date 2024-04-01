package don

import (
	"reflect"

	"github.com/abemedia/go-don/decoder"
	"github.com/abemedia/go-don/encoding"
	"github.com/abemedia/httprouter"
	"github.com/valyala/fasthttp"
)

type requestDecoder[V any] func(v *V, ctx *fasthttp.RequestCtx, p httprouter.Params) error

func newRequestDecoder[V any](v V) requestDecoder[V] {
	path, _ := decoder.NewCached(v, "path")
	query, _ := decoder.NewCached(v, "query")
	header, _ := decoder.NewCached(v, "header")

	if path == nil && query == nil && header == nil {
		return decodeBody[V]()
	}

	return decodeRequest(path, query, header)
}

func decodeRequest[V any](path, query, header *decoder.CachedDecoder[V]) requestDecoder[V] {
	body := decodeBody[V]()
	return func(v *V, ctx *fasthttp.RequestCtx, p httprouter.Params) error {
		if err := body(v, ctx, nil); err != nil {
			return err
		}

		val := reflect.ValueOf(v).Elem()

		if path != nil && len(p) > 0 {
			if err := path.DecodeValue((decoder.Params)(p), val); err != nil {
				return ErrNotFound
			}
		}

		if query != nil {
			if q := ctx.Request.URI().QueryArgs(); q.Len() > 0 {
				if err := query.DecodeValue((*decoder.Args)(q), val); err != nil {
					return err
				}
			}
		}

		if header != nil {
			if err := header.DecodeValue((*decoder.Header)(&ctx.Request.Header), val); err != nil {
				return err
			}
		}

		return nil
	}
}

func decodeBody[V any]() requestDecoder[V] {
	return func(v *V, ctx *fasthttp.RequestCtx, _ httprouter.Params) error {
		if ctx.Request.Header.ContentLength() == 0 || ctx.IsGet() || ctx.IsHead() {
			return nil
		}

		dec := encoding.GetDecoder(ctx.Request.Header.ContentType())
		if dec == nil {
			return ErrUnsupportedMediaType
		}

		return dec(ctx, v)
	}
}
