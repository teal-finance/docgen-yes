package docgen

import (
	"reflect"
	"testing"
)

func Test_copyDocRouter(t *testing.T) {
	type args struct {
		dr DocRouter
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
			if got := copyDocRouter(tt.args.dr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("copyDocRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}
