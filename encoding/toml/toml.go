// Package toml provides encoding and decoding of TOML data.
package toml

import (
	"github.com/abemedia/go-don/encoding"
	"github.com/pelletier/go-toml"
	"github.com/valyala/fasthttp"
)

func decode(ctx *fasthttp.RequestCtx, v any) error {
	return toml.NewDecoder(ctx.RequestBodyStream()).Decode(v)
}

func encode(ctx *fasthttp.RequestCtx, v any) error {
	ctx.SetContentType("application/toml; charset=utf-8")
	return toml.NewEncoder(ctx).Encode(v)
}

func init() {
	mediaType := "application/toml"

	encoding.RegisterDecoder(decode, mediaType)
	encoding.RegisterEncoder(encode, mediaType)
}
