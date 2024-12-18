package recaptcha

import (
	"errors"
	"net/http"
	"time"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/internal/byteconv"
	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

const endpoint = "https://www.google.com/recaptcha/api/siteverify"

var (
	ErrNoToken  = errors.New("missing token")
	ErrLowScore = errors.New("score too low")
)

//nolint:tagliatelle
type Response struct {
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	Action      string    `json:"action"`
	ErrorCodes  []string  `json:"error-codes"`
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
}

func Verify(secret, response, remoteip string) error {
	if response == "" {
		return ErrNoToken
	}

	form := fasthttp.AcquireArgs()
	form.Add("secret", secret)
	form.Add("response", response)
	form.Add("remoteip", remoteip)
	defer fasthttp.ReleaseArgs(form)

	_, b, err := fasthttp.Post(nil, endpoint, form)
	if err != nil {
		return err
	}

	var res Response
	if err = json.Unmarshal(b, &res); err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.ErrorCodes[0]) //nolint:goerr113
	}

	if res.Score < 0.5 {
		return ErrLowScore
	}

	return nil
}

func Middleware(secret string) don.Middleware {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			remoteip := ctx.RemoteIP().String()
			token := byteconv.Btoa(ctx.Request.Header.Peek("Recaptcha-Token"))

			if err := Verify(secret, token, remoteip); err != nil {
				ctx.Error(err.Error(), http.StatusForbidden)
				return
			}

			next(ctx)
		}
	}
}
