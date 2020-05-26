package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func NewServer(addr string) *http.Server {
	router := mux.NewRouter()

	return &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
}
