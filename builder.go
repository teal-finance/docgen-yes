package docgen

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func BuildDoc(r chi.Routes) (Doc, error) {
	d := Doc{}

	if getGoPath() == "" {
		return d, errors.New("docgen: unable to determine your $GOPATH")
	}

	// Walk and generate the router docs
	d.Router = BuildDocRouter(r)

	return d, nil
}

func BuildDocRouter(r chi.Routes) DocRouter {
	if r == nil {
		return DocRouter{}
	}

	rts := r
	dr := DocRouter{
		Middlewares: []DocMiddleware{},
		Routes:      map[string]DocRoute{},
	}
	dr.Routes = DocRoutes{}

	for _, mw := range rts.Middlewares() {
		dmw := DocMiddleware{
			FuncInfo: GetFuncInfo(mw),
		}
		dr.Middlewares = append(dr.Middlewares, dmw)
	}

	for _, rt := range rts.Routes() {
		drt := DocRoute{
			Pattern:  rt.Pattern,
			Handlers: DocHandlers{},
			Router: &DocRouter{
				Middlewares: []DocMiddleware{},
				Routes:      map[string]DocRoute{},
			},
		}

		if rt.SubRoutes != nil {
			subRoutes := rt.SubRoutes
			subDrts := BuildDocRouter(subRoutes)
			drt.Router = &subDrts
		} else {
			hall := rt.Handlers["*"]
			for method, h := range rt.Handlers {
				if method != "*" && hall != nil && fmt.Sprint(hall) == fmt.Sprint(h) {
					continue
				}

				dh := DocHandler{
					Middlewares: []DocMiddleware{},
					Method:      method,
					FuncInfo: FuncInfo{
						Pkg:          "",
						Func:         "",
						Comment:      "",
						File:         "",
						Line:         0,
						Anonymous:    false,
						Unresolvable: false,
					},
				}

				var endpoint http.Handler
				chain, _ := h.(*chi.ChainHandler)

				if chain != nil {
					for _, mw := range chain.Middlewares {
						dh.Middlewares = append(dh.Middlewares, DocMiddleware{
							FuncInfo: GetFuncInfo(mw),
						})
					}
					endpoint = chain.Endpoint
				} else {
					endpoint = h
				}

				dh.FuncInfo = GetFuncInfo(endpoint)

				drt.Handlers[method] = dh
			}
		}

		dr.Routes[rt.Pattern] = drt
	}

	return dr
}
