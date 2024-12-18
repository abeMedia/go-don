// Package yaml provides encoding and decoding of YAML data.
package yaml

import (
	"github.com/abemedia/go-don/encoding"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v3"
)

func decode(ctx *fasthttp.RequestCtx, v any) error {
	return yaml.NewDecoder(ctx.RequestBodyStream()).Decode(v)
}

func encode(ctx *fasthttp.RequestCtx, v any) error {
	ctx.SetContentType("application/yaml")
	return yaml.NewEncoder(ctx).Encode(v)
}

func init() {
	mediaType := "application/yaml"
	aliases := []string{"text/yaml", "application/x-yaml", "text/x-yaml", "text/vnd.yaml"}

	encoding.RegisterDecoder(decode, mediaType, aliases...)
	encoding.RegisterEncoder(encode, mediaType, aliases...)
}
