package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	envHTTPAddr = "HTTP_ADDR"
)

type Server struct {
	srv *http.Server
}

func NewServer(mux *http.ServeMux) *Server {
	srv := http.Server{
		Addr:         os.Getenv(envHTTPAddr),
		IdleTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
	}

	if srv.Addr == "" {
		log.Fatal().
			Err(fmt.Errorf("environment variable %q is not set", envHTTPAddr)).
			Send()
		return nil
	}

	return &Server{
		srv: &srv,
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.srv.Shutdown(ctx)
}
