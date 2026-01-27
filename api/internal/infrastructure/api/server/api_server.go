package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	http   *http.Server
}

func NewServer() *Server {
	engine := gin.New()

	engine.Use(gin.Recovery(), gin.Logger())

	return &Server{
		engine: engine,
	}
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}

func (s *Server) Start(addr string) error {
	s.http = &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.http == nil {
		return nil
	}
	return s.http.Shutdown(ctx)
}
