package docgen_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/forrest321/docgen"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestBaseTemplate(t *testing.T) {
	fmt.Println("Testing Base Template")

	//check Base Template for expected template fields: {title}, {css}, {favicon.ico}, {intro}, {routes}
	expectedFields := []string{"{title}", "{css}", "{favicon.ico}", "{intro}", "{routes}"}
	bt := docgen.BaseTemplate()
	for _, field := range expectedFields {
		assert.True(t, strings.Contains(bt, field))
	}

	fmt.Println("Base Template Testing Complete")
}

func TestListItem(t *testing.T) {
	fmt.Println("Testing List Item")

	itemText := "test"
	li := docgen.ListItem(itemText)
	assert.True(t, strings.Contains(li, "<li>"))
	assert.True(t, strings.Contains(li, "</li>"))
	assert.True(t, strings.Contains(li, itemText))

	fmt.Println("List Item Testing Complete")
}

func TestOrderedList(t *testing.T) {
	fmt.Println("Testing Ordered List")

	li := docgen.ListItem("test")
	ol := docgen.OrderedList(li)

	assert.True(t, strings.Contains(ol, "<ol>"))
	assert.True(t, strings.Contains(ol, "</ol>"))
	assert.True(t, strings.Contains(ol, li))

	fmt.Println("Ordered List Testing Complete")
}

func TestUnorderedList(t *testing.T) {
	fmt.Println("Testing Unordered List")

	li := docgen.ListItem("test")
	ul := docgen.UnorderedList(li)

	assert.True(t, strings.Contains(ul, "<ul>"))
	assert.True(t, strings.Contains(ul, "</ul>"))
	assert.True(t, strings.Contains(ul, li))

	fmt.Println("Ordered List Testing Complete")
}

func TestDiv(t *testing.T) {
	fmt.Println("Testing Div")
	testString := "test"
	div := docgen.Div(testString)
	assert.True(t, strings.Contains(div, "<div>"))
	assert.True(t, strings.Contains(div, "</div>"))
	assert.True(t, strings.Contains(div, testString))
	fmt.Println("Div Testing Complete")
}

func TestP(t *testing.T) {
	fmt.Println("Testing P")
	testString := "test"
	p := docgen.P(testString)
	assert.True(t, strings.Contains(p, "<p>"))
	assert.True(t, strings.Contains(p, "</p>"))
	assert.True(t, strings.Contains(p, testString))
	fmt.Println("P Testing Complete")
}

func TestHeading(t *testing.T) {
	fmt.Println("Testing Headings")
	testString := "test"

	//h1 - h6 are valid
	for index := 1; index <= 6; index++ {
		h := docgen.Head(index, testString)

		headerOpenTag := fmt.Sprintf("<h%v>", index)
		headerCloseTag := fmt.Sprintf("</h%v>", index)

		assert.True(t, strings.Contains(h, headerOpenTag))
		assert.True(t, strings.Contains(h, headerCloseTag))
		assert.True(t, strings.Contains(h, testString))
	}

	//anything below 1 should return an h1 tag
	hzero := docgen.Head(0, "zero | should return h1")
	assert.True(t, strings.Contains(hzero, "<h1>"))
	//anything above 6 should return an h6 tag
	hseven := docgen.Head(7, "seven | should return h6")
	assert.True(t, strings.Contains(hseven, "<h6>"))

	fmt.Println("Heading Testing Complete")
}

func TestCSS(t *testing.T) {
	//MilligramMinCSS
	css := docgen.MilligramMinCSS()
	assert.NotNil(t, css)
}

func setupRouter() chi.Router {

	var r, sr1, sr2, sr3, sr4, sr5, sr6 *chi.Mux
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
			sr1 = r.(*chi.Mux)
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
			})
			r.Route("/{hubID}", func(r chi.Router) {
				sr2 = r.(*chi.Mux)
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
					sr4 = r.(*chi.Mux)
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
					sr5 = r.(*chi.Mux)
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
			sr6 = r.(*chi.Mux)
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

	return r

	//fmt.Println(docgen.JSONRoutesDoc(r))

	//Markupdoc
	// docgen.PrintRoutes(r)

}

func TestMarkupRoutesDoc(t *testing.T) {
	type args struct {
		r    chi.Router
		opts docgen.MarkupOpts
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := docgen.MarkupRoutesDoc(tt.args.r, tt.args.opts); got != tt.want {
				t.Errorf("MarkupRoutesDoc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkupDoc_String(t *testing.T) {
	type fields struct {
		Opts          docgen.MarkupOpts
		Router        chi.Router
		Doc           docgen.Doc
		Routes        map[string]docgen.DocRouter
		FormattedHTML string
		RouteHTML     string
		buf           *bytes.Buffer
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := &docgen.MarkupDoc{
				Opts:          tt.fields.Opts,
				Router:        tt.fields.Router,
				Doc:           tt.fields.Doc,
				Routes:        tt.fields.Routes,
				FormattedHTML: tt.fields.FormattedHTML,
				RouteHTML:     tt.fields.RouteHTML,
				Buf:           tt.fields.buf,
			}
			if got := mu.String(); got != tt.want {
				t.Errorf("MarkupDoc.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkupDoc_generate(t *testing.T) {
	type fields struct {
		Opts          docgen.MarkupOpts
		Router        chi.Router
		Doc           docgen.Doc
		Routes        map[string]docgen.DocRouter
		FormattedHTML string
		RouteHTML     string
		buf           *bytes.Buffer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := &docgen.MarkupDoc{
				Opts:          tt.fields.Opts,
				Router:        tt.fields.Router,
				Doc:           tt.fields.Doc,
				Routes:        tt.fields.Routes,
				FormattedHTML: tt.fields.FormattedHTML,
				RouteHTML:     tt.fields.RouteHTML,
				Buf:           tt.fields.buf,
			}
			if err := mu.Generate(); (err != nil) != tt.wantErr {
				t.Errorf("MarkupDoc.generate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarkupDoc_writeRoutes(t *testing.T) {
	type fields struct {
		Opts          docgen.MarkupOpts
		Router        chi.Router
		Doc           docgen.Doc
		Routes        map[string]docgen.DocRouter
		FormattedHTML string
		RouteHTML     string
		buf           *bytes.Buffer
	}
	tests := []struct {
		name   string
		fields fields
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := &docgen.MarkupDoc{
				Opts:          tt.fields.Opts,
				Router:        tt.fields.Router,
				Doc:           tt.fields.Doc,
				Routes:        tt.fields.Routes,
				FormattedHTML: tt.fields.FormattedHTML,
				RouteHTML:     tt.fields.RouteHTML,
				Buf:           tt.fields.buf,
			}

			assert.NotNil(t, mu)
			//mu.writeRoutes()
		})
	}
}

func Test_printRouter(t *testing.T) {
	type args struct {
		mu    *docgen.MarkupDoc
		depth int
		dr    docgen.DocRouter
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//printRouter(tt.args.mu, tt.args.depth, tt.args.dr)
		})
	}
}

func Test_buildRoutesMap(t *testing.T) {
	type args struct {
		mu            *docgen.MarkupDoc
		parentPattern string
		ar            *docgen.DocRouter
		nr            *docgen.DocRouter
		dr            *docgen.DocRouter
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//buildRoutesMap(tt.args.mu, tt.args.parentPattern, tt.args.ar, tt.args.nr, tt.args.dr)
		})
	}
}
