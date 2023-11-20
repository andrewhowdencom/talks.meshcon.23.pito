package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

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

type Option func(srv *http.Server)

func New(o ...Option) *http.Server {
	srv := &http.Server{
		Addr:    ":80",
		Handler: http.HandlerFunc(rabbit),
	}

	for _, opt := range o {
		opt(srv)
	}

	return srv
}

func WithListenAddr(addr string) Option {
	return func(srv *http.Server) {
		srv.Addr = addr
	}
}

func rabbit(w http.ResponseWriter, r *http.Request) {
	// Bail out of the request is from a more nefarious character
	usr := r.Header.Get("User-Agent")
	if strings.Contains(usr, "Elmar Fudd") {
		w.WriteHeader(http.StatusGone)
		w.Write([]byte("Try again next time, Elmar."))

		return
	}

	w.Write([]byte(fmt.Sprintf(`(\_/)
(o.O) < %s
(")(")
`, Quotes[rand.Intn(len(Quotes))])))
}
