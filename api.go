// Package don provides a fast and efficient API framework.
package don

import (
	"bytes"
	"net"
	"net/http"

	"github.com/abemedia/httprouter"
	"github.com/valyala/fasthttp"
)

// DefaultEncoding contains the media type of the default encoding to fall back
// on if no `Accept` or `Content-Type` header was provided.
var DefaultEncoding = "text/plain"

type Middleware func(fasthttp.RequestHandler) fasthttp.RequestHandler

type Router interface {
	Get(path string, handle httprouter.Handle)
	Post(path string, handle httprouter.Handle)
	Put(path string, handle httprouter.Handle)
	Patch(path string, handle httprouter.Handle)
	Delete(path string, handle httprouter.Handle)
	Handle(method, path string, handle httprouter.Handle)
	Handler(method, path string, handle http.Handler)
	HandleFunc(method, path string, handle http.HandlerFunc)
	Group(path string) Router
	Use(mw ...Middleware)
}

type API struct {
	NotFound         fasthttp.RequestHandler
	MethodNotAllowed fasthttp.RequestHandler
	PanicHandler     func(*fasthttp.RequestCtx, any)

	router *httprouter.Router
	config *Config
	mw     []Middleware
}

type Config struct {
	// DefaultEncoding contains the mime type of the default encoding to fall
	// back on if no `Accept` or `Content-Type` header was provided.
	DefaultEncoding string

	// DisableNoContent controls whether a nil or zero value response should
	// automatically return 204 No Content with an empty body.
	DisableNoContent bool

	ForceDefaultEncoding bool
}

// New creates a new API instance.
func New(c *Config) *API {
	if c == nil {
		c = &Config{}
	}

	if c.DefaultEncoding == "" {
		c.DefaultEncoding = DefaultEncoding
	}

	return &API{
		router:           httprouter.New(),
		config:           c,
		NotFound:         E(ErrNotFound),
		MethodNotAllowed: E(ErrMethodNotAllowed),
	}
}

// Get is a shortcut for router.Handle(http.MethodGet, path, handle).
func (r *API) Get(path string, handle httprouter.Handle) {
	r.Handle(http.MethodGet, path, handle)
}

// Post is a shortcut for router.Handle(http.MethodPost, path, handle).
func (r *API) Post(path string, handle httprouter.Handle) {
	r.Handle(http.MethodPost, path, handle)
}

// Put is a shortcut for router.Handle(http.MethodPut, path, handle).
func (r *API) Put(path string, handle httprouter.Handle) {
	r.Handle(http.MethodPut, path, handle)
}

// Patch is a shortcut for router.Handle(http.MethodPatch, path, handle).
func (r *API) Patch(path string, handle httprouter.Handle) {
	r.Handle(http.MethodPatch, path, handle)
}

// Delete is a shortcut for router.Handle(http.MethodDelete, path, handle).
func (r *API) Delete(path string, handle httprouter.Handle) {
	r.Handle(http.MethodDelete, path, handle)
}

// Handle registers a new request handle with the given path and method.
func (r *API) Handle(method, path string, handle httprouter.Handle) {
	r.router.Handle(method, path, handle)
}

// Handler is an adapter which allows the usage of an http.Handler as a request handle.
func (r *API) Handler(method, path string, handle http.Handler) {
	r.router.Handler(method, path, handle)
}

// HandleFunc is an adapter which allows the usage of an http.HandlerFunc as a request handle.
func (r *API) HandleFunc(method, path string, handle http.HandlerFunc) {
	r.router.HandlerFunc(method, path, handle)
}

// Group creates a new sub-router with a common prefix.
func (r *API) Group(path string) Router {
	return &group{prefix: path, r: r}
}

// Use registers a middleware.
func (r *API) Use(mw ...Middleware) {
	r.mw = append(r.mw, mw...)
}

// RequestHandler creates a fasthttp.RequestHandler for the API.
func (r *API) RequestHandler() fasthttp.RequestHandler {
	r.router.NotFound = r.NotFound
	r.router.MethodNotAllowed = r.MethodNotAllowed
	r.router.PanicHandler = r.PanicHandler

	h := r.router.HandleFastHTTP
	for _, mw := range r.mw {
		h = mw(h)
	}

	return func(ctx *fasthttp.RequestCtx) {
		if r.config.ForceDefaultEncoding {
			ctx.Request.Header.SetContentType(r.config.DefaultEncoding)
			ctx.Request.Header.Set(fasthttp.HeaderAccept, r.config.DefaultEncoding)
		} else {
			contentType := ctx.Request.Header.ContentType()
			if len(contentType) == 0 || bytes.HasPrefix(contentType, anyEncoding) {
				ctx.Request.Header.SetContentType(r.config.DefaultEncoding)
			}
			accept := ctx.Request.Header.Peek(fasthttp.HeaderAccept)
			if len(accept) == 0 || bytes.HasPrefix(accept, anyEncoding) {
				ctx.Request.Header.Set(fasthttp.HeaderAccept, r.config.DefaultEncoding)
			}
		}

		h(ctx)

		// Content-Length of -3 means handler returned nil.
		if ctx.Response.Header.ContentLength() == -3 {
			ctx.Response.Header.Del(fasthttp.HeaderTransferEncoding)

			if !r.config.DisableNoContent {
				ctx.Response.SetBody(nil)

				if ctx.Response.StatusCode() == fasthttp.StatusOK {
					ctx.Response.SetStatusCode(fasthttp.StatusNoContent)
				}
			}
		}
	}
}

// ListenAndServe serves HTTP requests from the given TCP4 addr.
func (r *API) ListenAndServe(addr string) error {
	return newServer(r).ListenAndServe(addr)
}

// Serve serves incoming connections from the given listener.
func (r *API) Serve(ln net.Listener) error {
	return newServer(r).Serve(ln)
}

func newServer(r *API) *fasthttp.Server {
	return &fasthttp.Server{
		Handler:               r.RequestHandler(),
		StreamRequestBody:     true,
		NoDefaultContentType:  true,
		NoDefaultServerHeader: true,
	}
}

var anyEncoding = []byte("*/*")
