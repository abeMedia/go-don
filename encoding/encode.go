package encoding

import (
	"bytes"

	"github.com/abemedia/go-don/internal/byteconv"
	"github.com/valyala/fasthttp"
)

type ResponseEncoder = func(ctx *fasthttp.RequestCtx, v any) error

// RegisterEncoder registers a response encoder on a given media type.
func RegisterEncoder(enc ResponseEncoder, mime string, aliases ...string) {
	encoders[mime] = enc
	for _, alias := range aliases {
		encoders[alias] = enc
	}
}

// GetEncoder returns the response encoder for a given media type.
func GetEncoder(mime []byte) ResponseEncoder {
	var contentType []byte
	for len(mime) > 0 {
		if i := bytes.IndexByte(mime, ','); i >= 0 {
			contentType = mime[:i]
			mime = mime[i+1:]
		} else {
			contentType = mime
			mime = nil
		}
		if i := bytes.IndexByte(contentType, ';'); i > 0 {
			contentType = contentType[:i]
		}
		if enc := encoders[byteconv.Btoa(trim(contentType))]; enc != nil {
			return enc
		}
	}
	return nil
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

func trim(b []byte) []byte {
	start := 0
	for ; start < len(b); start++ {
		if asciiSpace[b[start]] == 0 {
			break
		}
	}

	stop := len(b)
	for ; stop > start; stop-- {
		if asciiSpace[b[stop-1]] == 0 {
			break
		}
	}

	return b[start:stop]
}

var encoders = map[string]ResponseEncoder{}
