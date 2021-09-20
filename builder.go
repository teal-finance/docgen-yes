package docgen

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func BuildDoc(r chi.Routes) (Doc, error) {
	d := Doc{}

	// Walk and generate the router docs
	d.Router = buildDocRouter(r)

	return d, nil
}

func buildDocRouter(r chi.Routes) DocRouter {
	rts := r
	dr := DocRouter{
		Middlewares: []DocMiddleware{},
		Routes:      map[string]DocRoute{},
	}
	drts := DocRoutes{}
	dr.Routes = drts

	for _, mw := range rts.Middlewares() {
		dmw := DocMiddleware{
			FuncInfo: GetFuncInfo(mw),
		}
		dr.Middlewares = append(dr.Middlewares, dmw)
	}

	for _, rt := range rts.Routes() {
		drt := DocRoute{Pattern: rt.Pattern, Handlers: DocHandlers{}}

		if rt.SubRoutes != nil {
			subRoutes := rt.SubRoutes
			subDrts := buildDocRouter(subRoutes)
			drt.Router = &subDrts
		} else {
			hall := rt.Handlers["*"]
			for method, h := range rt.Handlers {
				if method != "*" && hall != nil && fmt.Sprintf("%v", hall) == fmt.Sprintf("%v", h) {
					continue
				}

				dh := DocHandler{
					Middlewares: []DocMiddleware{},
					Method:      method,
					FuncInfo:    FuncInfo{},
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

		drts[rt.Pattern] = drt
	}

	return dr
}
