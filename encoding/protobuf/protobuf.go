// Package protobuf provides encoding and decoding of Protocol Buffers data.
package protobuf

import (
	"reflect"
	"sync"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/encoding"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"
)

var (
	messageType = reflect.TypeOf((*proto.Message)(nil)).Elem()
	cache       sync.Map
)

func decode(ctx *fasthttp.RequestCtx, v any) error {
	typ := reflect.TypeOf(v)
	fn, ok := cache.Load(typ)
	if !ok {
		if typ.Elem().Implements(messageType) {
			fn = func(v any) any { return reflect.ValueOf(v).Elem().Interface() }
		} else {
			fn = func(v any) any { return v }
		}
		cache.Store(typ, fn)
	}
	m, ok := fn.(func(v any) any)(v).(proto.Message)
	if !ok {
		return don.ErrUnsupportedMediaType
	}
	return proto.Unmarshal(ctx.Request.Body(), m)
}

func encode(ctx *fasthttp.RequestCtx, v any) error {
	ctx.SetContentType("application/protobuf")
	m, ok := v.(proto.Message)
	if !ok {
		return don.ErrNotAcceptable
	}
	b, err := proto.Marshal(m)
	if err == nil {
		ctx.Response.SetBodyRaw(b)
	}
	return err
}

func init() {
	mediaType := "application/protobuf"
	alias := "application/x-protobuf"

	encoding.RegisterDecoder(decode, mediaType, alias)
	encoding.RegisterEncoder(encode, mediaType, alias)
}
