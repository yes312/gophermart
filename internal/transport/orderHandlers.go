package transport

import (
	"context"
	"database/sql"
	"encoding/json"
	db "gophermart/internal/database"
	jwtpackage "gophermart/pkg/jwt"
	"gophermart/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

const (
	ApplicationJSON = "application/json"
)

type handlersData struct {
	ctx      context.Context
	storage  db.StoragerDB
	logger   *zap.SugaredLogger
	TokenExp time.Duration
	// навверное AuthToken можно(нужно) сделать через интерфейс
	AuthToken jwtpackage.Token
}

type authData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func New(ctx context.Context, storage db.StoragerDB, logger *zap.SugaredLogger) *handlersData {
	return &handlersData{
		ctx:     ctx,
		storage: storage,
		logger:  logger,
	}
}

func (h *handlersData) UploadOrders(w http.ResponseWriter, r *http.Request) {

	ordersNumber := chi.URLParam(r, "ordersNumber")

	valid, err := utils.IsValidOrderNumber(ordersNumber)
	if err != nil {
		http.Error(w, "wrong order number", http.StatusBadRequest)
		return
	}

	if !valid {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	key := UserID("user")
	UserID := r.Context().Value(key).(string)

	orderUserIDInterface, err := h.storage.WithRetry(h.ctx, h.storage.GetOrder(h.ctx, ordersNumber))
	orderUserID, _ := orderUserIDInterface.(db.OrderUserID)

	switch {
	case err == sql.ErrNoRows:

		_, err = h.storage.WithRetry(h.ctx, h.storage.AddOrder(h.ctx, orderUserID.OrderNumber, orderUserID.UserID))
		if err != nil {
			h.logger.Errorf("ошибка при добавлении заказа %w", ordersNumber)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.logger.Infof("заказ %w успешно добавлен", ordersNumber)

	case err != nil:
		h.logger.Errorf("Ошибка при получении заказа %w", ordersNumber)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	default:
		if orderUserID.OrderNumber == ordersNumber {
			if orderUserID.UserID == UserID {

				h.logger.Infof("заказ %s уже был загружен этим пользователем %s", ordersNumber, UserID)
				setResponseHeaders(w, ApplicationJSON, http.StatusOK)

			} else {

				h.logger.Infof("заказ %s уже был загружен другим пользователем %s", ordersNumber, UserID)
				setResponseHeaders(w, ApplicationJSON, http.StatusConflict)

			}
		}
	}

	h.logger.Infof("заказ %s уже был загружен этим пользователем %s", ordersNumber, UserID)
	setResponseHeaders(w, ApplicationJSON, http.StatusOK)

}

func (h *handlersData) GetUploadedOrders(w http.ResponseWriter, r *http.Request) {

	key := UserID("user")
	userID := r.Context().Value(key).(string)

	ordersInterface, err := h.storage.WithRetry(h.ctx, h.storage.GetOrders(h.ctx, userID))

	orders, _ := ordersInterface.([]db.OrderStatus)

	switch {
	case err == sql.ErrNoRows:

		h.logger.Info("нет данных о заказах")
		http.Error(w, err.Error(), http.StatusNoContent)
		return

	case err != nil:
		h.logger.Errorf("Ошибка запроса к базе: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	default:

		encoder := json.NewEncoder(w)
		err := encoder.Encode(orders)
		if err != nil {
			h.logger.Errorf("Ошибка маршалинга: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		setResponseHeaders(w, ApplicationJSON, http.StatusOK)
	}

}
