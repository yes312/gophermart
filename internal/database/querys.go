package db

import (
	"context"
	"database/sql"
)

func (storage *Storage) GetUser(ctx context.Context, login string) func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		getUserQuery := `SELECT user_id, hash from users WHERE login=$1;`
		var user User
		err := tx.QueryRowContext(ctx, getUserQuery, login).Scan(&user.Login, &user.Hash)

		return user, err
	}

}

func (storage *Storage) AddUser(ctx context.Context, login, hash string) func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		addUserQuery := `INSERT INTO users(user_id, hash) VALUES ($1, $2)`

		_, err := tx.ExecContext(ctx, addUserQuery, login, hash)
		return nil, err
	}
}

func (s *Storage) GetOrder(ctx context.Context, order string) func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		getOrderQuery := `
		SELECT orders.number, users.login, users.uid
		FROM orders
		LEFT JOIN users ON user_uid = uid
		WHERE orders.number = $1
	`

		var orderUserID OrderUserID
		err := tx.QueryRowContext(ctx, getOrderQuery, order).Scan(&orderUserID.OrderNumber, &orderUserID.UserID)
		if err != nil {
			return OrderUserID{}, err
		}

		return orderUserID, nil
	}
}

func (s *Storage) AddOrder(ctx context.Context, number string, uid string) func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		addOrderQuery := `INSERT INTO orders(number, user_uid) VALUES ($1, $2)`

		_, err := tx.ExecContext(ctx, addOrderQuery, number, uid)

		return nil, err
	}
}

func (s *Storage) GetOrders(ctx context.Context, userID string) func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		query := `SELECT orders.number, billing.status, billing.accrual, billing.uploaded_at
		FROM orders 
		JOIN billing ON orders.number = billing.number
		WHERE orders.user_id = $1
		AND billing.time = (
		SELECT MAX(time)
		FROM billing
		WHERE number = orders.number
	)`

		rows, err := s.DB.QueryContext(ctx, query, userID)
		defer rows.Close()
		if err != nil {
			return nil, err
		}

		orderStatusList := make([]OrderStatus, 0)
		for rows.Next() {
			var ordS OrderStatus
			err := rows.Scan(&ordS.Number, &ordS.Status, &ordS.Accrual, &ordS.UploadedAt)
			if err != nil {
				return nil, err
			}
			orderStatusList = append(orderStatusList, ordS)
		}

		return orderStatusList, nil
	}
}

func (s *Storage) GetBalance(ctx context.Context, user_id string) func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		getBalanceQuery := `
		SELECT
		SUM(CASE WHEN billing.status = 'PROCESSED' THEN billing.accrual ELSE 0 END) AS PROCESSED,
		SUM(CASE WHEN billing.status = 'WITHDRAWN' THEN billing.accrual ELSE 0 END) AS WITHDRAWN
		FROM orders 
		JOIN billing ON orders.number = billing.order_number
		WHERE orders.user_id =$1  AND billing.status IN ('PROCESSED', 'WITHDRAWN');`

		var balance Balance

		err := s.DB.QueryRowContext(ctx, getBalanceQuery, user_id).Scan(&balance.Current, &balance.Withdraw)

		if err != nil {
			return Balance{}, err
		}

		return balance, err
	}
}

// TODO retry сделать сделать через middleware!!
func (s *Storage) WithdrawBalance(ctx context.Context, userID string, orderSum OrderSum) func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		getBalanceQuery := `
		SELECT
		SUM(CASE WHEN billing.status = 'PROCESSED' THEN billing.accrual ELSE 0 END) AS PROCESSED,
		SUM(CASE WHEN billing.status = 'WITHDRAWN' THEN billing.accrual ELSE 0 END) AS WITHDRAWN
		FROM orders 
		JOIN billing ON orders.number = billing.order_number
		WHERE orders.user_id =$1  AND billing.status IN ('PROCESSED', 'WITHDRAWN');`

		var balance Balance

		err := s.DB.QueryRowContext(ctx, getBalanceQuery, userID).Scan(&balance.Current, &balance.Withdraw)

		if err != nil {
			return nil, err
		}
		if balance.Current-balance.Withdraw < orderSum.Sum {
			return nil, ErrNotEnoughFunds
		}

		addOrderQuery := `INSERT INTO billing (order_number, status, accrual, uploaded_at, time)
		VALUES ($1, 'WITHDRAWN', $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
	`
		_, err = s.DB.ExecContext(ctx, addOrderQuery, orderSum.OrderNumber, orderSum.Sum)

		return nil, err
	}
}

func (s *Storage) GetWithdrawals(ctx context.Context, userID string) func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		queryWithdrawals := `SELECT orders.number, billing.accrual AS sum, billing.uploaded_at AS processed_at
			FROM orders
			JOIN billing ON orders.number = billing.order_number
			WHERE orders.user_id = $1
			AND billing.status = 'WITHDRAWN'
			ORDER BY billing.uploaded_at ';`

		rows, err := s.DB.QueryContext(ctx, queryWithdrawals, userID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		withdrawalsList := make([]Withdrawal, 0)
		for rows.Next() {
			var w Withdrawal
			err := rows.Scan(&w.OrderNumber, &w.Sum, &w.ProcessedAt)
			if err != nil {
				return nil, err
			}
			withdrawalsList = append(withdrawalsList, w)
		}

		return withdrawalsList, nil
	}
}
