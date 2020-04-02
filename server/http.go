package server

import (
	"net/http"
	"time"

	"github.com/jeffbean/drawfh/game"
)

func NewServer(addr string) *http.Server {
	g := game.NewGameServer()
	mux := http.NewServeMux()

	mux.HandleFunc("/", g.HomeHandler)
	mux.HandleFunc("/room/create", g.RoomCreate)

	return &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mux, // Pass our instance of gorilla/mux in.
	}
}
