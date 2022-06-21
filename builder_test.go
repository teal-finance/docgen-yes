package docgen_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/teal-finance/docgen-yes"
)

func TestBuildDoc(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		routes  chi.Routes
		want    docgen.Doc
		wantErr bool
	}{{
		name:    "empty",
		routes:  nil,
		want:    docgen.Doc{},
		wantErr: false,
	}}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got, err := docgen.BuildDoc(c.routes)
			if (err != nil) != c.wantErr {
				t.Errorf("BuildDoc() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("BuildDoc() = %v, want %v", got, c.want)
			}
		})
	}
}

func Test_buildDocRouter(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		routes chi.Routes
		want   docgen.DocRouter
	}{{
		name:   "empty",
		routes: nil, want: docgen.DocRouter{},
	}}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			if got := docgen.BuildDocRouter(c.routes); !reflect.DeepEqual(got, c.want) {
				t.Errorf("buildDocRouter() = %v, want %v", got, c.want)
			}
		})
	}
}

type next struct{ called bool }

func (n *next) ServeHTTP(http.ResponseWriter, *http.Request) { n.called = true }

func Test_GetFuncInfo(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		i    interface{}
		want docgen.FuncInfo
	}{{
		name: "Unresolvable",
		i:    next{},
		want: docgen.FuncInfo{Unresolvable: true},
	}}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			if got := docgen.GetFuncInfo(c.i); !reflect.DeepEqual(got, c.want) {
				t.Errorf("GetFuncInfo() = %v, want %v", got, c.want)
			}
		})
	}
}
