package docgen_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/teal-finance/docgen-yes"
)

// RequestID comment goes here.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "requestID", "1")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func hubIndexHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s := fmt.Sprintf("/hubs/%s reqid:%s session:%s",
		chi.URLParam(r, "hubID"), ctx.Value("requestID"), ctx.Value("session.user"))
	w.Write([]byte(s))
}

// Generate docs for the MuxBig from chi/mux_test.go.
func TestMuxBig(t *testing.T) {
	// var sr1, sr2, sr3, sr4, sr5, sr6 *chi.Mux
	var r, sr3 *chi.Mux
	r = chi.NewRouter()
	r.Use(RequestID)

	// Some inline middleware, 1
	// We just love Go's ast tools
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), "session.user", "anonymous")
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})
		r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("fav"))
		})
		r.Get("/hubs/{hubID}/view", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			s := fmt.Sprintf("/hubs/%s/view reqid:%s session:%s", chi.URLParam(r, "hubID"),
				ctx.Value("requestID"), ctx.Value("session.user"))
			w.Write([]byte(s))
		})
		r.Get("/hubs/{hubID}/view/*", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			s := fmt.Sprintf("/hubs/%s/view/%s reqid:%s session:%s", chi.URLParamFromCtx(ctx, "hubID"),
				chi.URLParam(r, "*"), ctx.Value("requestID"), ctx.Value("session.user"))
			w.Write([]byte(s))
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), "session.user", "elvis")
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			s := fmt.Sprintf("/ reqid:%s session:%s", ctx.Value("requestID"), ctx.Value("session.user"))
			w.Write([]byte(s))
		})
		r.Get("/suggestions", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			s := fmt.Sprintf("/suggestions reqid:%s session:%s", ctx.Value("requestID"), ctx.Value("session.user"))
			w.Write([]byte(s))
		})

		r.Get("/woot/{wootID}/*", func(w http.ResponseWriter, r *http.Request) {
			s := fmt.Sprintf("/woot/%s/%s", chi.URLParam(r, "wootID"), chi.URLParam(r, "*"))
			w.Write([]byte(s))
		})

		r.Route("/hubs", func(r chi.Router) {
			_ = r.(*chi.Mux) // sr1
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
			})
			r.Route("/{hubID}", func(r chi.Router) {
				_ = r.(*chi.Mux) // sr2
				r.Get("/", hubIndexHandler)
				r.Get("/touch", func(w http.ResponseWriter, r *http.Request) {
					ctx := r.Context()
					s := fmt.Sprintf("/hubs/%s/touch reqid:%s session:%s", chi.URLParam(r, "hubID"),
						ctx.Value("requestID"), ctx.Value("session.user"))
					w.Write([]byte(s))
				})

				sr3 = chi.NewRouter()
				sr3.Get("/", func(w http.ResponseWriter, r *http.Request) {
					ctx := r.Context()
					s := fmt.Sprintf("/hubs/%s/webhooks reqid:%s session:%s", chi.URLParam(r, "hubID"),
						ctx.Value("requestID"), ctx.Value("session.user"))
					w.Write([]byte(s))
				})
				sr3.Route("/{webhookID}", func(r chi.Router) {
					_ = r.(*chi.Mux) // sr4
					r.Get("/", func(w http.ResponseWriter, r *http.Request) {
						ctx := r.Context()
						s := fmt.Sprintf("/hubs/%s/webhooks/%s reqid:%s session:%s", chi.URLParam(r, "hubID"),
							chi.URLParam(r, "webhookID"), ctx.Value("requestID"), ctx.Value("session.user"))
						w.Write([]byte(s))
					})
				})

				// TODO: /webooks is not coming up as a subrouter here...
				// we kind of want to wrap a Router... ?
				// perhaps add .Router() to the middleware inline thing..
				// and use that always.. or, can detect in that method..
				r.Mount("/webhooks", chi.Chain(func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "hook", true)))
					})
				}).Handler(sr3))

				// HMMMM.. only let Mount() for just a Router..?
				// r.Mount("/webhooks", Use(...).Router(sr3))
				// ... could this work even....?

				// HMMMMMMMMMMMMMMMMMMMMMMMM...
				// even if Mount() were to record all subhandlers mounted, we still couldn't get at the
				// routes

				r.Route("/posts", func(r chi.Router) {
					_ = r.(*chi.Mux) // sr5
					r.Get("/", func(w http.ResponseWriter, r *http.Request) {
						ctx := r.Context()
						s := fmt.Sprintf("/hubs/%s/posts reqid:%s session:%s", chi.URLParam(r, "hubID"),
							ctx.Value("requestID"), ctx.Value("session.user"))
						w.Write([]byte(s))
					})
				})
			})
		})

		r.Route("/folders/", func(r chi.Router) {
			_ = r.(*chi.Mux) // sr6
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				s := fmt.Sprintf("/folders/ reqid:%s session:%s",
					ctx.Value("requestID"), ctx.Value("session.user"))
				w.Write([]byte(s))
			})
			r.Get("/public", func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				s := fmt.Sprintf("/folders/public reqid:%s session:%s",
					ctx.Value("requestID"), ctx.Value("session.user"))
				w.Write([]byte(s))
			})
			r.Get("/in", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}).ServeHTTP)

			r.With(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "search", true)))
				})
			}).Get("/search", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("searching.."))
			})
		})
	})

	t.Logf("Done setting up chi.Router")
}

func TestPrintRoutes(t *testing.T) {
	type args struct {
		r chi.Routes
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docgen.PrintRoutes(tt.args.r)
		})
	}
}

func TestJSONRoutesDoc(t *testing.T) {
	cases := []struct {
		name string
		r    chi.Routes
		want string
	}{
		// TODO: Add test cases.
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := docgen.JSONRoutesDoc(c.r); got != c.want {
				t.Errorf("JSONRoutesDoc() = %v, want %v", got, c.want)
			}
		})
	}
}
