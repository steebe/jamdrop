package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/broothie/queuecumber/app"
	"github.com/broothie/queuecumber/db"
	"github.com/broothie/queuecumber/spotify"
	"github.com/gorilla/sessions"
)

type Server struct {
	App      *app.App
	Logger   *log.Logger
	Spotify  *spotify.Spotify
	DB       *db.DB
	Sessions *sessions.CookieStore
}

func New(app *app.App) *Server {
	return &Server{
		App:      app,
		Logger:   app.Logger,
		Spotify:  app.Spotify,
		DB:       app.DB,
		Sessions: sessions.NewCookieStore([]byte(app.Config.SecretKey)),
	}
}

func Run(app *app.App) {
	New(app).Run()
}

func (s *Server) Run() {
	s.Logger.Printf("serving @ %s", s.App.Config.BaseURL())
	s.Logger.Panic(http.ListenAndServe(fmt.Sprintf(":%d", s.App.Config.Port), s.Handler()))
}

func (s *Server) Handler() http.Handler {
	return ApplyLoggerMiddleware(s.Routes(), s.Logger)
}

func (s *Server) AppRedirect(w http.ResponseWriter, r *http.Request) {
	s.Logger.Println("server.AppRedirect")
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func (s *Server) Error(w http.ResponseWriter, error string, code int) {
	s.Logger.Println(error)
	http.Error(w, error, code)
}

func (s *Server) JSON(w http.ResponseWriter, data interface{}) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
