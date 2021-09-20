package docgen

import (
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestBuildDoc(t *testing.T) {
	type args struct {
		r chi.Routes
	}
	tests := []struct {
		name    string
		args    args
		want    Doc
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildDoc(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildDoc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildDoc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildDocRouter(t *testing.T) {
	type args struct {
		r chi.Routes
	}
	tests := []struct {
		name string
		args args
		want DocRouter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildDocRouter(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildDocRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetFuncInfo(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want FuncInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFuncInfo(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFuncInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
