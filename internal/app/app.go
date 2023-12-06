package app

import (
	"context"
	"gophermart/internal/config"
	db "gophermart/internal/database"
	"gophermart/internal/services"
	transport "gophermart/internal/transport/handlers"
	"log"
	"net/http"
	"sync"

	jwtpackage "gophermart/pkg/jwt"
	"gophermart/pkg/logger"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

const dbName = "gm"

// #ВопросМентору описание структуры сервера лучше тут или в пакете config?
type Server struct {
	ctx     context.Context
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

func New(ctx context.Context, config *config.Config) *Server {

	return &Server{
		ctx:    ctx,
		config: config,
	}
}

func (s *Server) Start(ctx context.Context, wg *sync.WaitGroup) error {

	log.Println("===Запуск сервера===")
	logger, err := logger.NewLogger(s.config.LoggerLevel)

	if err != nil {
		return err
	}
	s.logger = logger

	storage, err := db.New(ctx, s.config.DatabaseURI, dbName)
	if err != nil {
		return err
	}
	s.storage = storage

	s.mux = s.ConfigureMux()

	s.server = &http.Server{
		Addr:    s.config.RunAdress,
		Handler: s.mux,
	}

	a := services.NewAccrual(s.config.AccrualSysremAdress, s.config.AccrualRequestInterval, s.config.AccuralPuttingDBInterval, s.storage, s.logger)
	go a.RunAccrualRequester(ctx, wg)

	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	s.logger.Info("===Завершение работы сервера===")
	err := s.server.Shutdown(s.ctx)
	if err != nil {
		return err
	}
	err = s.storage.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) ConfigureMux() *chi.Mux {

	router := chi.NewRouter()
	handler := transport.New(s.ctx, s.storage, s.logger)
	handler.AuthToken = *jwtpackage.NewToken(s.config.TokenExp, s.config.Key)

	router.Route("/", func(r chi.Router) {

		r.Post("/api/user/register", handler.Registration)
		r.Post("/api/user/login", handler.Login)

		r.Post("/api/user/orders/{ordersNumber}", handler.AuthMiddleware(handler.UploadOrders))
		r.Get("/api/user/orders", handler.AuthMiddleware(handler.GetUploadedOrders))

		r.Get("/api/user/balance", handler.AuthMiddleware(handler.GetBalance))
		r.Post("/api/user/balance/withdraw", handler.AuthMiddleware(handler.WithdrawBalance))

		r.Get("/api/user/withdrawals", handler.GetWithdrawals)

	})

	return router
}
