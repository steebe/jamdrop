package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) Routes() http.Handler {
	root := mux.NewRouter()
	root.NotFoundHandler = http.RedirectHandler("/app", http.StatusPermanentRedirect)

	root.
		Methods(http.MethodGet).
		Path("/ping").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "pong") })

	root.
		Methods(http.MethodGet).
		Path("/app").
		Handler(s.RequireLoggedIn(s.Index()))

	root.
		Methods(http.MethodGet).
		PathPrefix("/public").
		Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	users := root.PathPrefix("/users").Subrouter()
	users.Use(s.RequireLoggedIn)
	users.
		Methods(http.MethodGet).
		Path("/follow").
		Handler(s.Share())

	users.
		Methods(http.MethodPost).
		Path("/queue").
		Handler(s.QueueSong())

	spotify := root.PathPrefix("/spotify").Subrouter()
	spotify.
		Methods(http.MethodGet).
		Path("/authorize").
		Handler(s.SpotifyAuthorize())

	spotify.
		Methods(http.MethodGet).
		Path("/authorize/callback").
		Handler(s.SpotifyAuthorizeCallback())

	if s.App.Config.Internal {
		jobs := root.PathPrefix("/jobs").Subrouter()
		jobs.
			Methods(http.MethodGet).
			Path("/eject_session_tokens").
			HandlerFunc(s.EjectSessionTokens)
	}

	printRoutes(root)
	return root
}

func printRoutes(router *mux.Router) {
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		methods, _ := route.GetMethods()
		if len(methods) == 0 {
			return nil
		}

		path, _ := route.GetPathTemplate()
		fmt.Printf("%v %s\n", methods, path)
		return nil
	})
}
