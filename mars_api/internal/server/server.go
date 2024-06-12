package server

import (
	"github.com/google/mars_api/internal/api"
	"github.com/google/mars_api/internal/conf"
	"github.com/google/mars_api/internal/storage"
)

type Server struct {
	storage   *storage.Storage
	apiServer *api.ApiServer
	conf      conf.Config
}

func NewServer(storage *storage.Storage, conf conf.Config) *Server {
	return &Server{storage: storage, apiServer: api.NewApiServer(storage, conf), conf: conf}
}

func (s *Server) Run() error {
	return s.apiServer.Run()
}
