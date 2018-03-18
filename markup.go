package docgen

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/go-chi/chi"
)

// MarkupDoc describes a document to be generated
type MarkupDoc struct {
	Opts          MarkupOpts
	Router        chi.Router
	Doc           Doc
	Routes        map[string]DocRouter // Pattern : DocRouter
	FormattedHTML string
	RouteHTML     string

	Buf *bytes.Buffer
}

// MarkupOpts describes the options for a MarkupDoc
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

	// RunHTTPServer determines if the generated HTML will be hosted on a server //TODO: web server import (quickweb)
	RunHTTPServer bool

	// HTTPServerPort determines the numerical port for the web server
	HTTPServerPort int
}

// MarkupRoutesDoc builds a document based on routes in a given router with given option set
func MarkupRoutesDoc(r chi.Router, opts MarkupOpts) string {
	mu := &MarkupDoc{Router: r, Opts: opts}
	if err := mu.Generate(); err != nil {
		return fmt.Sprintf("ERROR: %s\n", err.Error())
	}
	return mu.String()
}

// String pretty prints the document
func (mu *MarkupDoc) String() string {
	return mu.FormattedHTML
}

// Generate builds the document
func (mu *MarkupDoc) Generate() error {
	if mu.Router == nil {
		return errors.New("docgen: router is nil")
	}

	doc, err := BuildDoc(mu.Router)
	if err != nil {
		return err
	}

	mu.Doc = doc
	mu.Buf = &bytes.Buffer{}
	mu.Routes = make(map[string]DocRouter)

	mu.writeRoutes()

	r := strings.NewReplacer(
		"{title}", mu.Opts.ProjectPath,
		"{css}", MilligramMinCSS(),
		"{intro}", mu.Opts.Intro,
		"{routes}", mu.Opts.RouteText,
		"{favicon.ico}", FaviconIcoData(),
	)

	htmlString := r.Replace(BaseTemplate())
	mu.FormattedHTML = htmlString

	return nil
}

// writeRoutes generates the string for the Routes
func (mu *MarkupDoc) writeRoutes() {
	routesHeader := Head(2, "Routes")
	mu.Buf.WriteString(routesHeader)

	// Build a route tree that consists of the full route pattern
	// and the part of the tree for just that specific route, stored
	// in routes map on the markdown struct. This is the structure we
	// are going to render to markdown.
	dr := mu.Doc.Router
	ar := DocRouter{}
	buildRoutesMap(mu, "", &ar, &ar, &dr)

	routePaths := []string{}
	for pat := range mu.Routes {
		routePaths = append(routePaths, pat)
	}
	sort.Strings(routePaths)

	for _, pat := range routePaths {
		dr := mu.Routes[pat]
		mu.Buf.WriteString(fmt.Sprintf("<details>\n"))
		mu.Buf.WriteString(fmt.Sprintf("<summary>`%s`</summary>\n", pat))

		printRouter(mu, 0, dr)

		mu.Buf.WriteString(fmt.Sprintf("</details>\n"))
	}

	mu.Buf.WriteString(fmt.Sprintf("\n"))
	mu.Buf.WriteString(fmt.Sprintf("Total # of routes: %d\n", len(mu.Routes)))

	// TODO: total number of handlers..
	//return "oops"
}

// Generate the markdown to render the above structure
func printRouter(mu *MarkupDoc, depth int, dr DocRouter) {

	// Middlewares
	middleWares := make([]string, len(dr.Middlewares))
	for j, mw := range dr.Middlewares {
		middleWares[j] = ListItem(fmt.Sprintf("[%s](%s)", mw.Func, mu.githubSourceURL(mw.File, mw.Line)))
	}
	middleWaresList := UnorderedList(strings.Join(middleWares, ""))

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

				// Handler endpoint
				handlerEndpoint := fmt.Sprintf("[%s](%s)", dh.Func, mu.githubSourceURL(dh.File, dh.Line))
				mu.Buf.WriteString(handlerEndpoint)

				methods[mi] = ListItem(meth + " " + handlerEndpoint + "<br />" + Div(innerMiddlesList))
			}
			methodList := UnorderedList(strings.Join(methods, ""))
			routeListItems[ri] = ListItem(rt.Pattern + "<br />" + methodList)
		}
	}
	routeList := UnorderedList(strings.Join(routeListItems, ""))
	mu.RouteHTML = Head(3, "Middlewares") + Div(middleWaresList) + Head(3, "Routes") + Div(routeList)
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
