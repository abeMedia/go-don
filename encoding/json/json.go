// Package json provides encoding and decoding of JSON data.
package json

import (
	"github.com/abemedia/go-don/encoding"
	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

func decode(ctx *fasthttp.RequestCtx, v any) error {
	return json.Unmarshal(ctx.Request.Body(), v)
}

func encode(ctx *fasthttp.RequestCtx, v any) error {
	ctx.SetContentType("application/json")
	return json.NewEncoder(ctx).Encode(v)
}

func init() {
	mediaType := "application/json"

	encoding.RegisterDecoder(decode, mediaType)
	encoding.RegisterEncoder(encode, mediaType)
}
