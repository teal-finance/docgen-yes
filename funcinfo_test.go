package docgen

import (
	"reflect"
	"runtime"
	"testing"
)

func TestGetFuncInfo(t *testing.T) {
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

func Test_getCallerFrame(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want *runtime.Frame
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCallerFrame(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCallerFrame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPkgName(t *testing.T) {
	type args struct {
		file string
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
			if got := getPkgName(tt.args.file); got != tt.want {
				t.Errorf("getPkgName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFuncComment(t *testing.T) {
	type args struct {
		file string
		line int
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
			if got := getFuncComment(nil, tt.args.file, tt.args.line); got != tt.want {
				t.Errorf("getFuncComment() = %v, want %v", got, tt.want)
			}
		})
	}
}
