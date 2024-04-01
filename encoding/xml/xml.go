// Package xml provides encoding and decoding of XML data.
package xml

import (
	"encoding/xml"

	"github.com/abemedia/go-don/encoding"
	"github.com/valyala/fasthttp"
)

var Header = xml.Header

func decode(ctx *fasthttp.RequestCtx, v any) error {
	return xml.NewDecoder(ctx.RequestBodyStream()).Decode(v)
}

func encode(ctx *fasthttp.RequestCtx, v any) error {
	ctx.SetContentType("application/xml")
	_, _ = ctx.WriteString(Header)
	return xml.NewEncoder(ctx).Encode(v)
}

func init() {
	mediaType := "application/xml"
	alias := "text/xml"

	encoding.RegisterDecoder(decode, mediaType, alias)
	encoding.RegisterEncoder(encode, mediaType, alias)
}
