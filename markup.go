package docgen

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/go-chi/chi"
)

// MarkupDoc describes an HTML document
type MarkupDoc struct {
	Opts   MarkupOpts
	Router chi.Router
	Doc    Doc
	Routes map[string]DocRouter // Pattern : DocRouter

	buf *bytes.Buffer
}

// MarkupOpts describes the options for a MarkupDoc
type MarkupOpts struct {
	// ProjectPath is the base Go import path of the project
	ProjectPath string

	// Intro text included at the top of the generated markdown file.
	Intro string

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

// MarkupRoutesDoc builds a document based on routes in a given router
func MarkupRoutesDoc(r chi.Router, opts MarkupOpts) string {
	mu := &MarkupDoc{Router: r, Opts: opts}
	if err := mu.Generate(); err != nil {
		return fmt.Sprintf("ERROR: %s\n", err.Error())
	}
	return mu.String()
}

// String pretty prints the document
func (mu *MarkupDoc) String() string {
	return mu.buf.String()
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
	mu.buf = &bytes.Buffer{}
	mu.Routes = make(map[string]DocRouter)

	//TODO: get template, build strings, replace holders in template
	//TODO: if option is set, build and run web server
	mu.WriteIntro()
	mu.WriteRoutes()

	return nil
}

// WriteIntro writes the package name and provided intro string
func (mu *MarkupDoc) WriteIntro() {
	pkgName := mu.Opts.ProjectPath
	mu.buf.WriteString(fmt.Sprintf("# %s\n\n", pkgName))

	intro := mu.Opts.Intro
	mu.buf.WriteString(fmt.Sprintf("%s\n\n", intro))
}

// WriteRoutes generates the string for the Routes
func (mu *MarkupDoc) WriteRoutes() {
	mu.buf.WriteString(fmt.Sprintf("## Routes\n\n"))

	var buildRoutesMap func(parentPattern string, ar, nr, dr *DocRouter)
	buildRoutesMap = func(parentPattern string, ar, nr, dr *DocRouter) {

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
				buildRoutesMap(pattern, ar, nnr, rt.Router)

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

	// Build a route tree that consists of the full route pattern
	// and the part of the tree for just that specific route, stored
	// in routes map on the markdown struct. This is the structure we
	// are going to render to markdown.
	dr := mu.Doc.Router
	ar := DocRouter{}
	buildRoutesMap("", &ar, &ar, &dr)

	// Generate the markdown to render the above structure
	var printRouter func(depth int, dr DocRouter)
	printRouter = func(depth int, dr DocRouter) {

		tabs := ""
		for i := 0; i < depth; i++ {
			tabs += "\t"
		}

		// Middlewares
		for _, mw := range dr.Middlewares {
			mu.buf.WriteString(fmt.Sprintf("%s- [%s](%s)\n", tabs, mw.Func, mu.githubSourceURL(mw.File, mw.Line)))
		}

		// Routes
		for _, rt := range dr.Routes {
			mu.buf.WriteString(fmt.Sprintf("%s- **%s**\n", tabs, rt.Pattern))

			if rt.Router != nil {
				printRouter(depth+1, *rt.Router)
			} else {
				for meth, dh := range rt.Handlers {
					mu.buf.WriteString(fmt.Sprintf("%s\t- _%s_\n", tabs, meth))

					// Handler middlewares
					for _, mw := range dh.Middlewares {
						mu.buf.WriteString(fmt.Sprintf("%s\t\t- [%s](%s)\n", tabs, mw.Func, mu.githubSourceURL(mw.File, mw.Line)))
					}

					// Handler endpoint
					mu.buf.WriteString(fmt.Sprintf("%s\t\t- [%s](%s)\n", tabs, dh.Func, mu.githubSourceURL(dh.File, dh.Line)))
				}
			}
		}
	}

	routePaths := []string{}
	for pat := range mu.Routes {
		routePaths = append(routePaths, pat)
	}
	sort.Strings(routePaths)

	for _, pat := range routePaths {
		dr := mu.Routes[pat]
		mu.buf.WriteString(fmt.Sprintf("<details>\n"))
		mu.buf.WriteString(fmt.Sprintf("<summary>`%s`</summary>\n", pat))
		mu.buf.WriteString(fmt.Sprintf("\n"))
		printRouter(0, dr)
		mu.buf.WriteString(fmt.Sprintf("\n"))
		mu.buf.WriteString(fmt.Sprintf("</details>\n"))
	}

	mu.buf.WriteString(fmt.Sprintf("\n"))
	mu.buf.WriteString(fmt.Sprintf("Total # of routes: %d\n", len(mu.Routes)))

	// TODO: total number of handlers..
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

func htmlTemplate() string {
	return `
		<html>
			<head>
				{title}
				{css}
			</head>
			<body>
				<h1>{title}</h1>
				<p class='intro'>{intro}</p>
				<p class='routes'>{routes}</p>
			</body>
		</html>
	`
}

func htmlUnorderedList() string {
	return `
		<ul>
			{listItems}
		</ul>
	`
}

func htmlOrderedList() string {
	return `
		<ol>
			{listItems}
		</ol>
	`
}

func listItem(text string) string {
	return "<li>" + text + "</li>"
}

func milligramMinCSS() string {
	return `	
		/*!
		* Milligram v1.3.0
		* https://milligram.github.io
		*
		* Copyright (c) 2017 CJ Patoilo
		* Licensed under the MIT license
		*/
	
	*,*:after,*:before{box-sizing:inherit}html{box-sizing:border-box;font-size:62.5%}body{color:#606c76;font-family:'Roboto', 'Helvetica Neue', 'Helvetica', 'Arial', sans-serif;font-size:1.6em;font-weight:300;letter-spacing:.01em;line-height:1.6}blockquote{border-left:0.3rem solid #d1d1d1;margin-left:0;margin-right:0;padding:1rem 1.5rem}blockquote *:last-child{margin-bottom:0}.button,button,input[type='button'],input[type='reset'],input[type='submit']{background-color:#9b4dca;border:0.1rem solid #9b4dca;border-radius:.4rem;color:#fff;cursor:pointer;display:inline-block;font-size:1.1rem;font-weight:700;height:3.8rem;letter-spacing:.1rem;line-height:3.8rem;padding:0 3.0rem;text-align:center;text-decoration:none;text-transform:uppercase;white-space:nowrap}.button:focus,.button:hover,button:focus,button:hover,input[type='button']:focus,input[type='button']:hover,input[type='reset']:focus,input[type='reset']:hover,input[type='submit']:focus,input[type='submit']:hover{background-color:#606c76;border-color:#606c76;color:#fff;outline:0}.button[disabled],button[disabled],input[type='button'][disabled],input[type='reset'][disabled],input[type='submit'][disabled]{cursor:default;opacity:.5}.button[disabled]:focus,.button[disabled]:hover,button[disabled]:focus,button[disabled]:hover,input[type='button'][disabled]:focus,input[type='button'][disabled]:hover,input[type='reset'][disabled]:focus,input[type='reset'][disabled]:hover,input[type='submit'][disabled]:focus,input[type='submit'][disabled]:hover{background-color:#9b4dca;border-color:#9b4dca}.button.button-outline,button.button-outline,input[type='button'].button-outline,input[type='reset'].button-outline,input[type='submit'].button-outline{background-color:transparent;color:#9b4dca}.button.button-outline:focus,.button.button-outline:hover,button.button-outline:focus,button.button-outline:hover,input[type='button'].button-outline:focus,input[type='button'].button-outline:hover,input[type='reset'].button-outline:focus,input[type='reset'].button-outline:hover,input[type='submit'].button-outline:focus,input[type='submit'].button-outline:hover{background-color:transparent;border-color:#606c76;color:#606c76}.button.button-outline[disabled]:focus,.button.button-outline[disabled]:hover,button.button-outline[disabled]:focus,button.button-outline[disabled]:hover,input[type='button'].button-outline[disabled]:focus,input[type='button'].button-outline[disabled]:hover,input[type='reset'].button-outline[disabled]:focus,input[type='reset'].button-outline[disabled]:hover,input[type='submit'].button-outline[disabled]:focus,input[type='submit'].button-outline[disabled]:hover{border-color:inherit;color:#9b4dca}.button.button-clear,button.button-clear,input[type='button'].button-clear,input[type='reset'].button-clear,input[type='submit'].button-clear{background-color:transparent;border-color:transparent;color:#9b4dca}.button.button-clear:focus,.button.button-clear:hover,button.button-clear:focus,button.button-clear:hover,input[type='button'].button-clear:focus,input[type='button'].button-clear:hover,input[type='reset'].button-clear:focus,input[type='reset'].button-clear:hover,input[type='submit'].button-clear:focus,input[type='submit'].button-clear:hover{background-color:transparent;border-color:transparent;color:#606c76}.button.button-clear[disabled]:focus,.button.button-clear[disabled]:hover,button.button-clear[disabled]:focus,button.button-clear[disabled]:hover,input[type='button'].button-clear[disabled]:focus,input[type='button'].button-clear[disabled]:hover,input[type='reset'].button-clear[disabled]:focus,input[type='reset'].button-clear[disabled]:hover,input[type='submit'].button-clear[disabled]:focus,input[type='submit'].button-clear[disabled]:hover{color:#9b4dca}code{background:#f4f5f6;border-radius:.4rem;font-size:86%;margin:0 .2rem;padding:.2rem .5rem;white-space:nowrap}pre{background:#f4f5f6;border-left:0.3rem solid #9b4dca;overflow-y:hidden}pre>code{border-radius:0;display:block;padding:1rem 1.5rem;white-space:pre}hr{border:0;border-top:0.1rem solid #f4f5f6;margin:3.0rem 0}input[type='email'],input[type='number'],input[type='password'],input[type='search'],input[type='tel'],input[type='text'],input[type='url'],textarea,select{-webkit-appearance:none;-moz-appearance:none;appearance:none;background-color:transparent;border:0.1rem solid #d1d1d1;border-radius:.4rem;box-shadow:none;box-sizing:inherit;height:3.8rem;padding:.6rem 1.0rem;width:100%}input[type='email']:focus,input[type='number']:focus,input[type='password']:focus,input[type='search']:focus,input[type='tel']:focus,input[type='text']:focus,input[type='url']:focus,textarea:focus,select:focus{border-color:#9b4dca;outline:0}select{background:url('data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' height='14' viewBox='0 0 29 14' width='29'><path fill='#d1d1d1' d='M9.37727 3.625l5.08154 6.93523L19.54036 3.625'/></svg>') center right no-repeat;padding-right:3.0rem}select:focus{background-image:url('data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' height='14' viewBox='0 0 29 14' width='29'><path fill='#9b4dca' d='M9.37727 3.625l5.08154 6.93523L19.54036 3.625'/></svg>')}textarea{min-height:6.5rem}label,legend{display:block;font-size:1.6rem;font-weight:700;margin-bottom:.5rem}fieldset{border-width:0;padding:0}input[type='checkbox'],input[type='radio']{display:inline}.label-inline{display:inline-block;font-weight:normal;margin-left:.5rem}.container{margin:0 auto;max-width:112.0rem;padding:0 2.0rem;position:relative;width:100%}.row{display:flex;flex-direction:column;padding:0;width:100%}.row.row-no-padding{padding:0}.row.row-no-padding>.column{padding:0}.row.row-wrap{flex-wrap:wrap}.row.row-top{align-items:flex-start}.row.row-bottom{align-items:flex-end}.row.row-center{align-items:center}.row.row-stretch{align-items:stretch}.row.row-baseline{align-items:baseline}.row .column{display:block;flex:1 1 auto;margin-left:0;max-width:100%;width:100%}.row .column.column-offset-10{margin-left:10%}.row .column.column-offset-20{margin-left:20%}.row .column.column-offset-25{margin-left:25%}.row .column.column-offset-33,.row .column.column-offset-34{margin-left:33.3333%}.row .column.column-offset-50{margin-left:50%}.row .column.column-offset-66,.row .column.column-offset-67{margin-left:66.6666%}.row .column.column-offset-75{margin-left:75%}.row .column.column-offset-80{margin-left:80%}.row .column.column-offset-90{margin-left:90%}.row .column.column-10{flex:0 0 10%;max-width:10%}.row .column.column-20{flex:0 0 20%;max-width:20%}.row .column.column-25{flex:0 0 25%;max-width:25%}.row .column.column-33,.row .column.column-34{flex:0 0 33.3333%;max-width:33.3333%}.row .column.column-40{flex:0 0 40%;max-width:40%}.row .column.column-50{flex:0 0 50%;max-width:50%}.row .column.column-60{flex:0 0 60%;max-width:60%}.row .column.column-66,.row .column.column-67{flex:0 0 66.6666%;max-width:66.6666%}.row .column.column-75{flex:0 0 75%;max-width:75%}.row .column.column-80{flex:0 0 80%;max-width:80%}.row .column.column-90{flex:0 0 90%;max-width:90%}.row .column .column-top{align-self:flex-start}.row .column .column-bottom{align-self:flex-end}.row .column .column-center{-ms-grid-row-align:center;align-self:center}@media (min-width: 40rem){.row{flex-direction:row;margin-left:-1.0rem;width:calc(100% + 2.0rem)}.row .column{margin-bottom:inherit;padding:0 1.0rem}}a{color:#9b4dca;text-decoration:none}a:focus,a:hover{color:#606c76}dl,ol,ul{list-style:none;margin-top:0;padding-left:0}dl dl,dl ol,dl ul,ol dl,ol ol,ol ul,ul dl,ul ol,ul ul{font-size:90%;margin:1.5rem 0 1.5rem 3.0rem}ol{list-style:decimal inside}ul{list-style:circle inside}.button,button,dd,dt,li{margin-bottom:1.0rem}fieldset,input,select,textarea{margin-bottom:1.5rem}blockquote,dl,figure,form,ol,p,pre,table,ul{margin-bottom:2.5rem}table{border-spacing:0;width:100%}td,th{border-bottom:0.1rem solid #e1e1e1;padding:1.2rem 1.5rem;text-align:left}td:first-child,th:first-child{padding-left:0}td:last-child,th:last-child{padding-right:0}b,strong{font-weight:bold}p{margin-top:0}h1,h2,h3,h4,h5,h6{font-weight:300;letter-spacing:-.1rem;margin-bottom:2.0rem;margin-top:0}h1{font-size:4.6rem;line-height:1.2}h2{font-size:3.6rem;line-height:1.25}h3{font-size:2.8rem;line-height:1.3}h4{font-size:2.2rem;letter-spacing:-.08rem;line-height:1.35}h5{font-size:1.8rem;letter-spacing:-.05rem;line-height:1.5}h6{font-size:1.6rem;letter-spacing:0;line-height:1.4}img{max-width:100%}.clearfix:after{clear:both;content:' ';display:table/*# sourceMappingURL=milligram.min.css.map */}.float-left{float:left}.float-right{float:right}
	`
}
