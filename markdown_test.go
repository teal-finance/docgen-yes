package docgen

import (
	"bytes"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestMarkdownRoutesDoc(t *testing.T) {
	type args struct {
		r    chi.Router
		opts MarkdownOpts
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
			if got := MarkdownRoutesDoc(tt.args.r, tt.args.opts); got != tt.want {
				t.Errorf("MarkdownRoutesDoc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkdownDoc_String(t *testing.T) {
	type fields struct {
		Opts   MarkdownOpts
		Router chi.Router
		Doc    Doc
		Routes map[string]DocRouter
		buf    *bytes.Buffer
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
			md := &MarkdownDoc{
				Opts:   tt.fields.Opts,
				Router: tt.fields.Router,
				Doc:    tt.fields.Doc,
				Routes: tt.fields.Routes,
				buf:    tt.fields.buf,
			}
			if got := md.String(); got != tt.want {
				t.Errorf("MarkdownDoc.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkdownDoc_Generate(t *testing.T) {
	type fields struct {
		Opts   MarkdownOpts
		Router chi.Router
		Doc    Doc
		Routes map[string]DocRouter
		buf    *bytes.Buffer
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
			md := &MarkdownDoc{
				Opts:   tt.fields.Opts,
				Router: tt.fields.Router,
				Doc:    tt.fields.Doc,
				Routes: tt.fields.Routes,
				buf:    tt.fields.buf,
			}
			if err := md.Generate(); (err != nil) != tt.wantErr {
				t.Errorf("MarkdownDoc.Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarkdownDoc_WriteIntro(t *testing.T) {
	type fields struct {
		Opts   MarkdownOpts
		Router chi.Router
		Doc    Doc
		Routes map[string]DocRouter
		buf    *bytes.Buffer
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := &MarkdownDoc{
				Opts:   tt.fields.Opts,
				Router: tt.fields.Router,
				Doc:    tt.fields.Doc,
				Routes: tt.fields.Routes,
				buf:    tt.fields.buf,
			}
			md.WriteIntro()
		})
	}
}

func TestMarkdownDoc_WriteRoutes(t *testing.T) {
	type fields struct {
		Opts   MarkdownOpts
		Router chi.Router
		Doc    Doc
		Routes map[string]DocRouter
		buf    *bytes.Buffer
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := &MarkdownDoc{
				Opts:   tt.fields.Opts,
				Router: tt.fields.Router,
				Doc:    tt.fields.Doc,
				Routes: tt.fields.Routes,
				buf:    tt.fields.buf,
			}
			md.WriteRoutes()
		})
	}
}

func TestMarkdownDoc_githubSourceURL(t *testing.T) {
	type fields struct {
		Opts   MarkdownOpts
		Router chi.Router
		Doc    Doc
		Routes map[string]DocRouter
		buf    *bytes.Buffer
	}
	type args struct {
		file string
		line int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := &MarkdownDoc{
				Opts:   tt.fields.Opts,
				Router: tt.fields.Router,
				Doc:    tt.fields.Doc,
				Routes: tt.fields.Routes,
				buf:    tt.fields.buf,
			}
			if got := md.githubSourceURL(tt.args.file, tt.args.line); got != tt.want {
				t.Errorf("MarkdownDoc.githubSourceURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
