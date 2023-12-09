package transport

import (
	"database/sql"
	"encoding/json"
	"errors"
	db "gophermart/internal/database"
	"gophermart/utils"
	"net/http"
)

func (h *handlersData) GetBalance(w http.ResponseWriter, r *http.Request) {

	key := UserID("user")
	userID := r.Context().Value(key).(string)

	balanceInterface, err := h.storage.WithRetry(h.ctx, h.storage.GetBalance(h.ctx, userID))
	setContentType(w, ApplicationJSON)
	if err != nil {
		h.logger.Errorf("ошибка при получении баланса: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balance, _ := balanceInterface.(db.Balance)

	//  #ВОПРОСМЕНТОРУ может выделить маршалинг и отправку в JSON в отдельную функцию?
	encoder := json.NewEncoder(w)
	err = encoder.Encode(balance)
	if err != nil {
		h.logger.Errorf("Ошибка маршалинга: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	setResponseHeaders(w, ApplicationJSON, http.StatusOK)

}

func (h *handlersData) WithdrawBalance(w http.ResponseWriter, r *http.Request) {

	// 	Возможные коды ответа:
	// 200 — успешная обработка запроса;
	// 401 — пользователь не авторизован;
	// 402 — на счету недостаточно средств;
	// 422 — неверный номер заказа;
	// 500 — внутренняя ошибка сервера.

	key := UserID("user")
	userID := r.Context().Value(key).(string)

	var data db.OrderSum
	setContentType(w, ApplicationJSON)
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// #ВОПРОСМЕНТОРУ Этот кусочек ниже повторяется в других хендлерах. стоит выделить его в отдельную функцию?
	// может в middleware?
	valid, err := utils.IsValidOrderNumber(data.OrderNumber)
	if err != nil {
		http.Error(w, "wrong order number", http.StatusUnprocessableEntity)
		return
	}
	if !valid {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	_, err = h.storage.WithRetry(h.ctx, h.storage.WithdrawBalance(h.ctx, userID, data))
	switch {
	case errors.Is(err, db.ErrNotEnoughFunds):

		h.logger.Errorf("На пользователя %s счете недостаточно баллов", userID)
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return

	case err != nil:

		h.logger.Errorf("ошибка получения данных из БД %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	default:
		h.logger.Infof("баллы по заказу %s списаны с баланса пользователя %s", data.OrderNumber, userID)
		setResponseHeaders(w, ApplicationJSON, http.StatusOK)
	}
}

func (h *handlersData) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	// 204 — нет ни одного списания.
	// 401 — пользователь не авторизован.
	// 500 — внутренняя ошибка сервера.

	key := UserID("user")
	userID := r.Context().Value(key).(string)

	withdrawalsInterface, err := h.storage.WithRetry(h.ctx, h.storage.GetWithdrawals(h.ctx, userID))
	setContentType(w, ApplicationJSON)
	switch {
	case errors.Is(err, sql.ErrNoRows):

		h.logger.Info("нет данных о выводе средств")
		http.Error(w, err.Error(), http.StatusNoContent)
		return

	case err != nil:
		h.logger.Errorf("Ошибка запроса к базе: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	default:
		withdrawals, _ := withdrawalsInterface.([]db.Withdrawal)
		encoder := json.NewEncoder(w)
		err := encoder.Encode(withdrawals)
		if err != nil {
			h.logger.Errorf("Ошибка маршалинга: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		setResponseHeaders(w, ApplicationJSON, http.StatusOK)
	}

}
