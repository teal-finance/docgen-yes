package docgen

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"
)

// MarkupDoc describes a document to be generated.
type MarkupDoc struct {
	Opts          MarkupOpts
	Router        chi.Router
	Doc           Doc
	Routes        map[string]DocRouter // Pattern : DocRouter
	FormattedHTML string
	RouteHTML     string
}

// MarkupOpts describes the options for a MarkupDoc.
type MarkupOpts struct {
	// ProjectPath is the base Go import path of the project
	ProjectPath string

	// Intro text included at the top of the generated markdown file.
	Intro string

	// RouteText contains HTML generated from Route metadata
	RouteText string

	// ForceRelativeLinks to be relative even if they're not on github
	ForceRelativeLinks bool

	// URLMap allows specifying a map of package import paths to their link sources
	// Used for mapping vendored dependencies to their upstream sources
	// For example:
	// map[string]string{"github.com/my/package/vendor/go-chi/chi/": "https://github.com/go-chi/chi/blob/master/"}
	URLMap map[string]string
}

// MarkupRoutesDoc builds a document based on routes in a given router with given option set.
func MarkupRoutesDoc(r chi.Router, opts MarkupOpts) string {
	// Goal: rewrite this class to have a single exported function(?)
	//    that returns a formatted HTML document based on a MarkupOpts
	//     1) flatten router, build docs w/o recursion
	//     2) Alternatively, build the JSON or Markdown doc and convert to HTML

	mu := &MarkupDoc{
		Opts:   opts,
		Router: r,
		Doc: Doc{Router: DocRouter{
			Middlewares: []DocMiddleware{},
			Routes:      map[string]DocRoute{},
		}},
		Routes:        map[string]DocRouter{},
		FormattedHTML: "",
		RouteHTML:     "",
	}
	if err := mu.generate(); err != nil {
		return fmt.Sprintf("ERROR: %s\n", err.Error())
	}

	formattedHTML := mu.String()
	return formattedHTML
}

// String pretty prints the document.
func (mu *MarkupDoc) String() string {
	return mu.FormattedHTML
}

// Generate builds the document.
func (mu *MarkupDoc) generate() error {
	if mu.Router == nil {
		return errors.New("docgen: router is nil")
	}

	doc, err := BuildDoc(mu.Router)
	if err != nil {
		return err
	}

	mu.Doc = doc
	//  mu.Buf = &bytes.Buffer{}
	mu.Routes = make(map[string]DocRouter)

	mu.writeRoutes()

	r := strings.NewReplacer(
		"{title}", "go-chi Docgen",
		"{css}", BassCSS(),
		"{intro}", "Ding! Routes are Done!",
		"{routes}", mu.RouteHTML,
		"{favicon.ico}", FaviconIcoData(),
	)

	htmlString := r.Replace(BaseTemplate())
	mu.FormattedHTML = htmlString

	return nil
}

// writeRoutes generates the string for the Routes.
func (mu *MarkupDoc) writeRoutes() {
	routesHeader := Head(2, "Routes")
	mu.RouteHTML += routesHeader
	// mu.Buf.WriteString(routesHeader)

	// Build a route tree that consists of the full route pattern
	// and the part of the tree for just that specific route, stored
	// in routes map on the markdown struct. This is the structure we
	// are going to render to markdown.
	dr := mu.Doc.Router
	ar := DocRouter{
		Middlewares: []DocMiddleware{},
		Routes:      map[string]DocRoute{},
	}
	buildRoutesMap(mu, "", &ar, &ar, &dr)

	routePaths := []string{}
	for pat := range mu.Routes {
		routePaths = append(routePaths, pat)
	}
	sort.Strings(routePaths)

	for _, pat := range routePaths {
		dr := mu.Routes[pat]
		mu.RouteHTML += "<div>" + P(pat)

		printRouter(mu, 0, dr)

		mu.RouteHTML += "</div>"
	}
}

// Generate the markup to render the above structure.
func printRouter(mu *MarkupDoc, depth int, dr DocRouter) {
	// Middlewares
	middleWares := make([]string, len(dr.Middlewares))
	for j, mw := range dr.Middlewares {
		middleWares[j] = ListItem(fmt.Sprintf("[%s](%s)", mw.Func, mu.githubSourceURL(mw.File, mw.Line)))
	}
	middleWaresList := UnorderedList(strings.Join(middleWares, ""))
	mu.RouteHTML += Div(Head(3, "Middlewares") + middleWaresList)
	// Routes
	routeListItems := make([]string, len(dr.Routes))
	ri := -1
	for _, rt := range dr.Routes {
		ri++

		// RECURSION AAAAHHHHH NOOOOO
		if rt.Router != nil {
			printRouter(mu, depth+1, *rt.Router)
		} else {
			// Route Handler Methods
			methods := make([]string, len(rt.Handlers))
			mi := -1
			for meth, dh := range rt.Handlers {
				mi++

				innerMiddles := make([]string, len(dh.Middlewares))
				imi := -1
				// Handler middlewares
				for _, mw := range dh.Middlewares {
					imi++
					innerMiddles[imi] = ListItem(fmt.Sprintf("[%s](%s)", mw.Func, mu.githubSourceURL(mw.File, mw.Line)))
				}
				innerMiddlesList := UnorderedList(strings.Join(innerMiddles, ""))
				handlerEndpoint := fmt.Sprintf("[%s](%s)", dh.Func, mu.githubSourceURL(dh.File, dh.Line))
				methods[mi] = ListItem(meth + " " + handlerEndpoint + "<br />" + Div(innerMiddlesList))
			}
			methodList := UnorderedList(strings.Join(methods, ""))
			routeListItems[ri] = ListItem(rt.Pattern + "<br />" + methodList)
		}
	}
	routeList := UnorderedList(strings.Join(routeListItems, ""))
	mu.RouteHTML += Head(3, "Middlewares") + Div(middleWaresList) + Head(3, "Routes") + Div(routeList)
}

func buildRoutesMap(mu *MarkupDoc, parentPattern string, ar, nr, dr *DocRouter) {
	nr.Middlewares = append(nr.Middlewares, dr.Middlewares...)

	for pat, rt := range dr.Routes {
		pattern := parentPattern + pat

		nr.Routes = DocRoutes{}

		if rt.Router != nil {
			nnr := &DocRouter{}
			nr.Routes[pat] = DocRoute{
				Pattern:  pat,
				Handlers: rt.Handlers,
				Router:   nnr,
			}
			buildRoutesMap(mu, pattern, ar, nnr, rt.Router)
		} else if len(rt.Handlers) > 0 {
			nr.Routes[pat] = DocRoute{
				Pattern:  pat,
				Handlers: rt.Handlers,
				Router:   nil,
			}

			// Remove the trailing slash if the handler is a subroute for "/"
			routeKey := pattern
			if pat == "/" && len(routeKey) > 1 {
				routeKey = routeKey[:len(routeKey)-1]
			}
			mu.Routes[routeKey] = copyDocRouter(*ar)
		} else {
			panic("not possible")
		}
	}
}

func (mu *MarkupDoc) githubSourceURL(file string, line int) string {
	// Currently, we only automatically link to source for github projects
	if strings.Index(file, "github.com/") != 0 && !mu.Opts.ForceRelativeLinks {
		return ""
	}
	if mu.Opts.ProjectPath == "" {
		return ""
	}
	for pkg, url := range mu.Opts.URLMap {
		if idx := strings.Index(file, pkg); idx >= 0 {
			pos := idx + len(pkg)
			url = strings.TrimRight(url, "/")
			filepath := strings.TrimLeft(file[pos:], "/")
			return fmt.Sprintf("%s/%s#L%d", url, filepath, line)
		}
	}
	if idx := strings.Index(file, mu.Opts.ProjectPath); idx >= 0 {
		// relative
		pos := idx + len(mu.Opts.ProjectPath)
		return fmt.Sprintf("%s#L%d", file[pos:], line)
	}
	// absolute
	return fmt.Sprintf("https://%s#L%d", file, line)
}
