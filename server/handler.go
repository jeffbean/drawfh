package server

import (
	"io"
	"net/http"
)

type drawServer struct {
}

// Handler is the main draw server handler
func Handler() http.Handler {
	return &drawServer{}
}

func (d *drawServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world v2!")
}
