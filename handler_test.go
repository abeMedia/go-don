package don_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/json"
	_ "github.com/abemedia/go-don/encoding/text"
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/abemedia/httprouter"
	"github.com/google/go-cmp/cmp"
	"github.com/valyala/fasthttp"
)

func TestHandlerRequest(t *testing.T) {
	type request struct {
		Path   string `path:"path"`
		Header string `header:"Header"`
		Query  string `query:"query"`
		JSON   string `json:"json"`
	}

	var got request

	want := request{
		Path:   "path",
		Header: "header",
		Query:  "query",
		JSON:   "json",
	}

	api := don.New(&don.Config{DefaultEncoding: "application/json"})

	api.Post("/:path", don.H(func(ctx context.Context, req request) (any, error) {
		got = req
		return nil, nil
	}))

	api.RequestHandler()(httptest.NewRequest(
		fasthttp.MethodPost,
		"/path?query=query",
		`{"json":"json"}`,
		map[string]string{"header": "header"},
	))

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}

func TestHandlerResponse(t *testing.T) {
	type request struct {
		url    string
		body   string
		header map[string]string
	}

	type response struct {
		Code   int
		Body   string
		Header map[string]string
	}

	tests := []struct {
		message string
		want    response
		config  *don.Config
		route   string
		handler httprouter.Handle
		request request
	}{
		{
			message: "should return no content",
			want: response{
				Code: fasthttp.StatusNoContent,
				Body: "",
				Header: map[string]string{
					"Content-Type": "text/plain; charset=utf-8",
				},
			},
			handler: don.H(func(ctx context.Context, req any) (any, error) {
				return nil, nil
			}),
		},
		{
			message: "should return null",
			want: response{
				Code: fasthttp.StatusOK,
				Body: "null\n",
				Header: map[string]string{
					"Content-Type": "application/json",
				},
			},
			config: &don.Config{DefaultEncoding: "application/json", DisableNoContent: true},
			handler: don.H(func(ctx context.Context, req any) (any, error) {
				return nil, nil
			}),
		},
		{
			message: "should set response headers",
			want: response{
				Code: fasthttp.StatusOK,
				Header: map[string]string{
					"Content-Type": "text/plain; charset=utf-8",
					"Foo":          "bar",
				},
			},
			handler: don.H(func(ctx context.Context, req any) (any, error) {
				return &headerer{}, nil
			}),
		},
		{
			message: "should return error on unprocessable request",
			want: response{
				Code:   fasthttp.StatusUnsupportedMediaType,
				Body:   "Unsupported Media Type",
				Header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			},
			handler: don.H(func(ctx context.Context, req struct{ Hello string }) (any, error) {
				return nil, nil
			}),
			request: request{body: `{"foo":"bar"}`},
		},
		{
			message: "should return error on unacceptable",
			want: response{
				Code:   fasthttp.StatusNotAcceptable,
				Body:   "Not Acceptable",
				Header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			},
			handler: don.H(func(ctx context.Context, req any) ([]string, error) {
				return []string{"foo", "bar"}, nil
			}),
			request: request{header: map[string]string{"Accept": "text/plain; charset=utf-8"}},
		},
		{
			message: "should return error on unsupported accept",
			want: response{
				Code:   fasthttp.StatusNotAcceptable,
				Body:   "Not Acceptable",
				Header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			},
			handler: don.H(func(ctx context.Context, req any) (any, error) {
				return nil, nil
			}),
			request: request{header: map[string]string{"Accept": "application/msword"}},
		},
		{
			message: "should return error on unsupported content type",
			want: response{
				Code:   fasthttp.StatusUnsupportedMediaType,
				Body:   "Unsupported Media Type",
				Header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			},
			handler: don.H(func(ctx context.Context, req any) (any, error) {
				return nil, nil
			}),
			request: request{
				body:   `foo`,
				header: map[string]string{"Content-Type": "application/msword"},
			},
		},
		{
			message: "should return error on invalid query",
			want: response{
				Code:   fasthttp.StatusBadRequest,
				Body:   "strconv.ParseInt: parsing \"foo\": invalid syntax",
				Header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			},
			handler: don.H(func(ctx context.Context, req struct {
				Test int `query:"test"`
			},
			) (any, error) {
				return req, nil
			}),
			request: request{url: "/?test=foo"},
		},
		{
			message: "should return error on invalid header",
			want: response{
				Code:   fasthttp.StatusBadRequest,
				Body:   "strconv.ParseInt: parsing \"foo\": invalid syntax",
				Header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			},
			handler: don.H(func(ctx context.Context, req struct {
				Test int `header:"Test"`
			},
			) (any, error) {
				return req, nil
			}),
			request: request{header: map[string]string{"Test": "foo"}},
		},
		{
			message: "should return error on invalid path element",
			want: response{
				Code:   fasthttp.StatusNotFound,
				Body:   "Not Found",
				Header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			},
			handler: don.H(func(ctx context.Context, req struct {
				Test int `path:"test"`
			},
			) (any, error) {
				return req, nil
			}),
			route:   "/:test",
			request: request{url: "/foo"},
		},
		{
			message: "should return error on invalid body",
			want: response{
				Code:   fasthttp.StatusBadRequest,
				Body:   "strconv.Atoi: parsing \"foo\": invalid syntax",
				Header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			},
			handler: don.H(func(ctx context.Context, req int) (any, error) {
				return req, nil
			}),
			request: request{body: "foo"},
		},
		{
			message: "should return internal server error",
			want: response{
				Code:   fasthttp.StatusInternalServerError,
				Body:   "test",
				Header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			},
			handler: don.H(func(ctx context.Context, req any) (any, error) {
				return nil, errors.New("test")
			}),
		},
	}

	for _, test := range tests {
		ctx := httptest.NewRequest(fasthttp.MethodPost, test.request.url, test.request.body, test.request.header)

		api := don.New(test.config)
		api.Post("/"+strings.TrimPrefix(test.route, "/"), test.handler)
		api.RequestHandler()(ctx)

		res := response{ctx.Response.StatusCode(), string(ctx.Response.Body()), map[string]string{}}
		ctx.Response.Header.VisitAll(func(key, value []byte) { res.Header[string(key)] = string(value) })

		if diff := cmp.Diff(test.want, res); diff != "" {
			t.Errorf("%s:\n%s", test.message, diff)
		}
	}
}

type headerer struct{}

func (h *headerer) String() string {
	return ""
}

func (h *headerer) Header() http.Header {
	return http.Header{"Foo": []string{"bar"}}
}

func BenchmarkHandler(b *testing.B) {
	type request struct {
		Path   string `path:"path"`
		Header string `header:"Header"`
		Query  string `query:"query"`
	}

	header := map[string]string{"Header": "header", "Accept": "text/plain"}
	ctx := httptest.NewRequest("POST", "/path?query=query", "", header)
	p := httprouter.Params{{Key: "path", Value: "path"}}

	b.Run("Request", func(b *testing.B) {
		h := don.H(func(ctx context.Context, req request) (any, error) { return nil, nil })
		for i := 0; i < b.N; i++ {
			h(ctx, p)
		}
	})

	b.Run("RequestPointer", func(b *testing.B) {
		h := don.H(func(ctx context.Context, req *request) (any, error) { return nil, nil })
		for i := 0; i < b.N; i++ {
			h(ctx, p)
		}
	})
}
