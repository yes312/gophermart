package transport

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	db "gophermart/internal/database"
	jwtpackage "gophermart/pkg/jwt"
	"gophermart/utils"
	"io"
	"log"
	"net/http"
	"time"

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

	// ordersNumber := chi.URLParam(r, "ordersNumber")

	body, err := (io.ReadAll(r.Body))
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	ordersNumber := string(body)
	log.Println("будет добавлен заказ", ordersNumber)

	valid, err := utils.IsValidOrderNumber(ordersNumber)
	if err != nil {
		http.Error(w, "wrong order number", http.StatusBadRequest)
		return
	}

	if !valid {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		h.logger.Errorf("путой юзер детектед")
		http.Error(w, "wrong user id", http.StatusUnauthorized)
		return
	}

	orderUserIDInterface, err := h.storage.WithRetry(h.ctx, h.storage.AddOrder(h.ctx, ordersNumber, userID))
	orderUserID, _ := orderUserIDInterface.(db.OrderUserID)
	log.Println("orderUserID", orderUserID)
	if err != nil {

		h.logger.Errorf("Ошибка при получении заказа %w", ordersNumber)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	if orderUserID.OrderNumber == ordersNumber {
		if orderUserID.UserID == userID {

			h.logger.Infof("заказ %s уже был загружен этим пользователем %s", ordersNumber, userID)
			setResponseHeaders(w, ApplicationJSON, http.StatusOK)
			return

		} else {

			h.logger.Infof("заказ %s уже был загружен другим пользователем %s", ordersNumber, userID)
			setResponseHeaders(w, ApplicationJSON, http.StatusConflict)
			return
		}
	}
	// +++++++++++++++++++++++++++
	ordersInterface, err := h.storage.WithRetry(h.ctx, h.storage.GetOrders(h.ctx, userID))

	orders, ok := ordersInterface.([]db.OrderStatus)
	h.logger.Info("===============================", orders, ok, err)
	// +++++++++++++++++++++++++++
	h.logger.Infof("заказ %s загружен пользователем %s", ordersNumber, userID)
	setResponseHeaders(w, ApplicationJSON, http.StatusAccepted)

}

func (h *handlersData) GetUploadedOrders(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		h.logger.Errorf("путой юзер детектед")
		http.Error(w, "wrong user id", http.StatusUnauthorized)
		return
	}

	ordersInterface, err := h.storage.WithRetry(h.ctx, h.storage.GetOrders(h.ctx, userID))

	orders, ok := ordersInterface.([]db.OrderStatus)
	h.logger.Info("===============================", orders, ok)
	switch {
	case errors.Is(err, sql.ErrNoRows) || !ok:
		// если данных нет, то эта лшибка не выпадает. т.к. err=nil
		h.logger.Info("нет данных о заказах")
		http.Error(w, err.Error(), http.StatusNoContent)
		return

	case err != nil:
		h.logger.Errorf("Ошибка запроса к базе: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	default:

		if len(orders) == 0 {
			h.logger.Info("нет данных о заказах")
			setResponseHeaders(w, ApplicationJSON, http.StatusNoContent)
			return
		}

		encoder := json.NewEncoder(w)
		err := encoder.Encode(orders)

		if err != nil {
			h.logger.Errorf("Ошибка маршалинга: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
