// Package msgpack provides encoding and decoding of MessagePack data.
package msgpack

import (
	"github.com/abemedia/go-don/encoding"
	"github.com/valyala/fasthttp"
	"github.com/vmihailenco/msgpack/v5"
)

func decode(ctx *fasthttp.RequestCtx, v any) error {
	dec := msgpack.GetDecoder()
	dec.UsePreallocateValues(true)
	dec.Reset(ctx.RequestBodyStream())
	err := dec.Decode(v)
	msgpack.PutDecoder(dec)
	return err
}

func encode(ctx *fasthttp.RequestCtx, v any) error {
	ctx.SetContentType("application/msgpack")
	b, err := msgpack.Marshal(v)
	if err == nil {
		ctx.Response.SetBodyRaw(b)
	}
	return err
}

func init() {
	mediaType := "application/msgpack"
	aliases := []string{"application/x-msgpack", "application/vnd.msgpack"}

	encoding.RegisterDecoder(decode, mediaType, aliases...)
	encoding.RegisterEncoder(encode, mediaType, aliases...)
}
