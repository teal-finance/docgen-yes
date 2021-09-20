package raml_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/micheartin/docgen-yes"
	"github.com/micheartin/docgen-yes/raml"
	yaml "gopkg.in/yaml.v2"
)

func TestWalkerRAML(t *testing.T) {
	r := Router()

	ramlDocs := &raml.RAML{
		Title:         "Big Mux",
		BaseURI:       "https://bigmux.example.com",
		Protocols:     []string{},
		MediaType:     "application/json",
		Version:       "v1.0",
		Documentation: []raml.Documentation{},
		Resources:     map[string]*raml.Resource{},
	}

	if err := chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		handlerInfo := docgen.GetFuncInfo(handler)
		resource := &raml.Resource{
			DisplayName:     "",
			Description:     handlerInfo.Comment,
			Responses:       map[int]raml.Response{},
			Body:            map[string]raml.Example{},
			Is:              []string{},
			Example:         "",
			SecuredBy:       []string{},
			URIParameters:   map[string]raml.Example{},
			QueryParameters: map[string]raml.Example{},
			Resources:       map[string]*raml.Resource{},
		}

		return ramlDocs.Add(method, route, resource)
	}); err != nil {
		t.Error(err)
	}

	if _, err := yaml.Marshal(ramlDocs); err != nil {
		t.Error(err)
	}
}

// Copy-pasted from _examples/raml. We can't simply import it, since it's main pkg.
func Router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	// RESTy routes for "articles" resource
	r.Route("/articles", func(r chi.Router) {
		r.With(paginate).Get("/", ListArticles)
		r.Post("/", CreateArticle)       // POST /articles
		r.Get("/search", SearchArticles) // GET /articles/search

		r.Route("/:articleID", func(r chi.Router) {
			r.Use(ArticleCtx)            // Load the *Article on the request context
			r.Get("/", GetArticle)       // GET /articles/123
			r.Put("/", UpdateArticle)    // PUT /articles/123
			r.Delete("/", DeleteArticle) // DELETE /articles/123
		})
	})

	// Mount the admin sub-router, the same as a call to
	// Route("/admin", func(r chi.Router) { with routes here })
	r.Mount("/admin", adminRouter())

	return r
}

type Article struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Article fixture data.
var articles = []*Article{
	{ID: "1", Title: "Hi"},
	{ID: "2", Title: "sup"},
}

// ArticleCtx middleware is used to load an Article object from
// the URL parameters passed through as the request. In case
// the Article could not be found, we stop here and return a 404.
func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID := chi.URLParam(r, "articleID")
		article, err := dbGetArticle(articleID)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, http.StatusText(http.StatusNotFound))
			return
		}
		ctx := context.WithValue(r.Context(), "article", article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Search Articles.
// Searches the Articles data for a matching article.
// It's just a stub, but you get the idea.
func SearchArticles(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, articles)
}

// List Articles.
// Returns an array of Articles.
func ListArticles(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, articles)
}

// Create new Article.
// Ppersists the posted Article and returns it
// back to the client as an acknowledgement.
func CreateArticle(w http.ResponseWriter, r *http.Request) {
	article := &Article{
		ID:    "",
		Title: "",
	}

	render.JSON(w, r, article)
}

// Get a specific Article.
func GetArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*Article)

	render.JSON(w, r, article)
}

// Update a specific Article.
// Updates an existing Article in our persistent store.
func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*Article)

	render.JSON(w, r, article)
}

// Delete a specific Article.
// Removes an existing Article from our persistent store.
func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*Article)

	render.JSON(w, r, article)
}

// A completely separate router for administrator routes.
func adminRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(AdminOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("admin: index"))
	})
	r.Get("/accounts", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("admin: list accounts.."))
	})
	r.Get("/users/:userId", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf("admin: view user id %v", chi.URLParam(r, "userId"))))
	})
	return r
}

// AdminOnly middleware restricts access to just administrators.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value("acl.admin").(bool)
		if !ok || !isAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

//--

// dbGetArticle is a helper function that mock some kind of storage.
func dbGetArticle(id string) (*Article, error) {
	for _, a := range articles {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, errors.New("article not found")
}
