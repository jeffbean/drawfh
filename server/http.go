package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jeffbean/drawfh/game"
)

func NewServer(addr string) *http.Server {
	router := NewGameHandler()

	return &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
}

func NewGameHandler() http.Handler {
	router := mux.NewRouter()
	g := game.NewGameServer()

	// TODO: a page to create a room via a html form, that uses another route to post the form.
	router.HandleFunc("/room/create", g.RoomCreate).Methods("POST", "GET")

	router.HandleFunc("/room/{id}", g.RoomDetail).Methods("GET")

	// Routes Below expect json forms
	router.HandleFunc("/room/{id}/join", g.RoomJoin).Methods("POST")
	router.HandleFunc("/room/{id}/delete", g.RoomDelete).Methods("POST")
	return router
}
