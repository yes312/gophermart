package app

import (
	"context"
	"gophermart/internal/config"
	db "gophermart/internal/database"
	transport "gophermart/internal/transport"
	"log"
	"net/http"

	"gophermart/pkg/logger"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// #ВопросМентору описание структуры сервера лучше тут или в пакете config?
type Server struct {
	server  *http.Server
	config  *config.Config
	mux     *chi.Mux
	storage Storager
	logger  *zap.SugaredLogger
}

var _ Storager = &db.Storage{}

type Storager interface {
	db.StoragerDB
}

func New(config *config.Config) *Server {

	return &Server{
		config: config,
		// mux:    chi.NewRouter(),
	}
}

func (s *Server) Start(ctx context.Context) error {

	log.Println("===Запуск сервера===")
	logger, err := logger.NewLogger(s.config.LoggerLevel)

	if err != nil {
		return err
	}
	s.logger = logger

	storage, err := db.New(ctx, s.config.DatabaseURI)
	if err != nil {
		return err
	}
	s.storage = storage

	s.ConfigureMux()

	s.server = &http.Server{
		Addr:    s.config.RunAdress,
		Handler: s.mux,
	}
	return s.server.ListenAndServe()
}

func (s *Server) Close() {

}

func (s *Server) ConfigureMux() *chi.Mux {
	r := chi.NewRouter()
	h := transport.New(s.storage)

	r.Route("/", func(r chi.Router) {

		r.Post("/api/user/register", http.HandlerFunc(h.Registration))

	})
	return r
}
