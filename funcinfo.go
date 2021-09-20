package docgen

import (
	"go/parser"
	"go/token"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// FuncInfo describes a function's metadata.
type FuncInfo struct {
	Pkg          string `json:"pkg"`
	Func         string `json:"func"`
	Comment      string `json:"comment"`
	File         string `json:"file,omitempty"`
	Line         int    `json:"line,omitempty"`
	Anonymous    bool   `json:"anonymous,omitempty"`
	Unresolvable bool   `json:"unresolvable,omitempty"`
}

// GetFuncInfo returns a FuncInfo object for a given interface.
func GetFuncInfo(i interface{}) FuncInfo {
	fi := FuncInfo{
		Pkg:          "",
		Func:         "",
		Comment:      "",
		File:         "",
		Line:         0,
		Anonymous:    false,
		Unresolvable: false,
	}
	frame := getCallerFrame(i)
	goPathSrc := filepath.Join(os.Getenv("GOPATH"), "src")

	if frame == nil {
		fi.Unresolvable = true

		return fi
	}

	pkgName := getPkgName(frame.File)
	if pkgName == "chi" {
		fi.Unresolvable = true
	}
	funcPath := frame.Func.Name()

	idx := strings.Index(funcPath, "/"+pkgName)
	if idx > 0 {
		fi.Pkg = funcPath[:idx+1+len(pkgName)]
		fi.Func = funcPath[idx+2+len(pkgName):]
	} else {
		fi.Func = funcPath
	}

	if strings.Index(fi.Func, ".func") > 0 {
		fi.Anonymous = true
	}

	fi.File = frame.File
	fi.Line = frame.Line
	if filepath.HasPrefix(fi.File, goPathSrc) {
		fi.File = fi.File[len(goPathSrc)+1:]
	}

	// Check if file info is unresolvable
	if !strings.Contains(funcPath, pkgName) {
		fi.Unresolvable = true
	}

	if !fi.Unresolvable {
		fi.Comment = getFuncComment(i, frame.File, frame.Line)
	}

	return fi
}

func getCallerFrame(i interface{}) *runtime.Frame {
	val := reflect.ValueOf(i)
	var pc uintptr

	switch val.Kind() {
	case reflect.Func:
		pc = val.Pointer()
	case reflect.Ptr:
		typ := reflect.TypeOf(i)
		handlerType := reflect.TypeOf(new(http.Handler)).Elem()

		if typ.Implements(handlerType) {
			if method, ok := typ.Elem().MethodByName("ServeHTTP"); ok {
				pc = method.Func.Pointer()
			}
		}
	}

	frames := runtime.CallersFrames([]uintptr{pc})
	if frames == nil {
		return nil
	}
	frame, _ := frames.Next()
	if frame.Entry == 0 {
		return nil
	}
	return &frame
}

func getPkgName(file string) string {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, file, nil, parser.PackageClauseOnly)
	if err != nil {
		return ""
	}
	if astFile.Name == nil {
		return ""
	}
	return astFile.Name.Name
}

func getFuncComment(i interface{}, file string, line int) string {
	fset := token.NewFileSet()

	astFile, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return ""
	}

	if len(astFile.Comments) == 0 {
		return ""
	}
	typNames := strings.Split(reflect.TypeOf(i).String(), ".")
	typName := typNames[len(typNames)-1]

	for _, cmt := range astFile.Comments {
		if strings.HasPrefix(cmt.Text(), typName+" ") {
			return cmt.Text()
		}
	}

	for _, cmt := range astFile.Comments {
		if fset.Position(cmt.End()).Line+1 == line {
			return cmt.Text()
		}
	}

	return ""
}
