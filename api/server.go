package api

import (
	"log"
	"net/http"

	"github.com/suvrick/go-kiss-core/api/controllers/frame"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() {

	mux := http.NewServeMux()

	frame.NewFrameController(mux)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panic(err)
	}
}
