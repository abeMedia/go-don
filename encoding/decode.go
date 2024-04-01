package encoding

import (
	"bytes"

	"github.com/abemedia/go-don/internal/byteconv"
	"github.com/valyala/fasthttp"
)

type RequestDecoder = func(ctx *fasthttp.RequestCtx, v any) error

// RegisterDecoder registers a request decoder for a given media type.
func RegisterDecoder(dec RequestDecoder, mime string, aliases ...string) {
	decoders[mime] = dec
	for _, alias := range aliases {
		decoders[alias] = dec
	}
}

// GetDecoder returns the request decoder for a given media type.
func GetDecoder(mime []byte) RequestDecoder {
	if i := bytes.IndexByte(mime, ';'); i > 0 {
		mime = mime[:i]
	}
	return decoders[byteconv.Btoa(bytes.TrimSpace(mime))]
}

var decoders = map[string]RequestDecoder{}
