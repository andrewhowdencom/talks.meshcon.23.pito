package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

const pkgName = "rabbit.local/server"

var Quotes = []string{
	"Eh, What's up, Doc?",
	"Carrots are devine… You get a dozen for a dime, It’s maaaa-gic!",
	"For shame, doc. Hunting rabbits with an elephant gun. Why don’t you shoot yourself an elephant?",
	"Gee, ain’t I a stinker?",
	"I bet you say that to all the wabbits",
	"If it’s the Captain’s Mess, let him clean it up.",
	"Well, it’s 5 o’clock somewhere.",
	"Oh, well, we almost had a romantic ending!",
}

type Srv struct {
	t trace.Tracer

	*http.Server
}

type Option func(srv *Srv)

func New(o ...Option) *http.Server {

	srv := &Srv{
		t: noop.NewTracerProvider().Tracer(pkgName),
	}

	srv.Server = &http.Server{
		Addr:    ":80",
		Handler: http.HandlerFunc(srv.rabbit),
	}

	for _, opt := range o {
		opt(srv)
	}

	return srv.Server
}

func WithListenAddr(addr string) Option {
	return func(srv *Srv) {
		srv.Addr = addr
	}
}

func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(srv *Srv) {
		srv.t = tp.Tracer(pkgName)
	}
}

func (s *Srv) rabbit(w http.ResponseWriter, r *http.Request) {
	_, span := s.t.Start(r.Context(), "rabbit")
	defer span.End()

	// Set default status
	span.SetAttributes(semconv.HTTPStatusCode(http.StatusOK))

	// Track some critical user information. ⚠️ In Europe, this is personal data and needs to be treated
	// accordingly.
	span.SetAttributes(semconv.UserAgentOriginal(r.Header.Get("User-Agent")))

	// Bail out of the request is from a more nefarious character
	usr := r.Header.Get("User-Agent")
	if strings.Contains(usr, "Elmar Fudd") {
		// Override the status
		span.SetAttributes(semconv.HTTPStatusCode(http.StatusGone))
		w.WriteHeader(http.StatusGone)
		w.Write([]byte("Try again next time, Elmar."))

		return
	}

	w.Write([]byte(fmt.Sprintf(`(\_/)
(o.O) < %s
(")(")
`, Quotes[rand.Intn(len(Quotes))])))
}
