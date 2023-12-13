package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

func (storage *Storage) GetUser(ctx context.Context, login string) dbOperation {

	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		getUserQuery := `SELECT user_id, hash from users WHERE user_id=$1;`
		var user User
		err := tx.QueryRowContext(ctx, getUserQuery, login).Scan(&user.Login, &user.Hash)

		return user, err
	}

}

func (storage *Storage) AddUser(ctx context.Context, UserID, hash string) dbOperation {
	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		addUserQuery := `INSERT INTO users(user_id, hash) VALUES ($1, $2)`

		_, err := tx.ExecContext(ctx, addUserQuery, UserID, hash)
		return nil, err
	}
}

// func (storage *Storage) GetOrder(ctx context.Context, order string) dbOperation {
// 	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
// 		getOrderQuery := `
// 		SELECT orders.number, users.user_id
// 		FROM orders
// 		LEFT JOIN users ON orders.user_id = users.user_id
// 		WHERE orders.number = $1`

// 		var orderUserID OrderUserID
// 		err := tx.QueryRowContext(ctx, getOrderQuery, order).Scan(&orderUserID.OrderNumber, &orderUserID.UserID)
// 		if err != nil {
// 			return OrderUserID{}, err
// 		}

// 		return orderUserID, nil
// 	}
// }

func (storage *Storage) AddOrder(ctx context.Context, orderNumber string, userID string) dbOperation {

	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		getOrderQuery := `SELECT number, user_id FROM orders
						  WHERE orders.number = $1`

		var orderUserID OrderUserID
		err := tx.QueryRowContext(ctx, getOrderQuery, orderNumber).Scan(&orderUserID.OrderNumber, &orderUserID.UserID)
		log.Println("ERROR!!!", err)
		switch {
		case err == sql.ErrNoRows:
			t := time.Now()
			addOrderQuery := `INSERT INTO orders(number, user_id, uploaded_at) VALUES ($1, $2, $3)`
			_, err = tx.ExecContext(ctx, addOrderQuery, orderNumber, userID, t)
			if err != nil {
				return OrderUserID{}, err
			}

			addBillingQuery := `INSERT INTO billing (order_number, status, accrual, uploaded_at, time)
								VALUES ($1, 'NEW', 0, $2, CURRENT_TIMESTAMP)`
			_, err = tx.ExecContext(ctx, addBillingQuery, orderNumber, t)
			log.Println("ERROR!!!222 ", err)
			if err != nil {
				return OrderUserID{}, err
			}

		case err != nil:
			return OrderUserID{}, err
		}

		return orderUserID, err
	}
}

func (storage *Storage) GetOrders(ctx context.Context, userID string) dbOperation {

	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		query := `SELECT orders.number, billing.status, billing.accrual as accrual, billing.uploaded_at
				 FROM orders 
				 JOIN billing ON orders.number = billing.order_number
				 WHERE orders.user_id = $1
				 AND billing.time = (
				 SELECT MAX(time)
				 FROM billing
				 WHERE billing.order_number = orders.number
				 AND billing.status != 'WITHDRAWN')
				 ORDER BY billing.time ASC,orders.uploaded_at ASC`

		rows, err := tx.QueryContext(ctx, query, userID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var orderStatusList []OrderStatusNew
		for rows.Next() {
			var ordS OrderStatusNew
			err := rows.Scan(&ordS.Number, &ordS.Status, &ordS.Accrual, &ordS.UploadedAt)
			if err != nil {
				return nil, err
			}
			ordS.Accrual = ordS.Accrual / 100
			orderStatusList = append(orderStatusList, ordS)
		}

		if err := rows.Err(); err != nil {
			return orderStatusList, err
		}

		return orderStatusList, nil
	}
}

func (storage *Storage) GetBalance(ctx context.Context, userID string) dbOperation {
	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		getBalanceQuery := `	
		SELECT
		COALESCE(SUM(CASE WHEN billing.status = 'PROCESSED' THEN billing.accrual ELSE 0 END),0) AS PROCESSED,
		COALESCE(SUM(CASE WHEN billing.status = 'WITHDRAWN' THEN billing.accrual ELSE 0 END),0) AS WITHDRAWN
		FROM orders 
		JOIN billing ON orders.number = billing.order_number
		WHERE orders.user_id =$1  AND billing.status IN ('PROCESSED', 'WITHDRAWN');`

		var balance Balance

		err := tx.QueryRowContext(ctx, getBalanceQuery, userID).Scan(&balance.Current, &balance.Withdraw)

		if err != nil {
			return Balance{}, err
		}

		balance.Current = balance.Current / 100
		balance.Withdraw = balance.Withdraw / 100
		return balance, err
	}
}

func (storage *Storage) WithdrawBalance(ctx context.Context, userID string, orderSum OrderSum) dbOperation {

	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		getBalanceQuery := `
		SELECT
		COALESCE(SUM(CASE WHEN billing.status = 'PROCESSED' THEN billing.accrual ELSE 0 END),0) AS PROCESSED,
		COALESCE(SUM(CASE WHEN billing.status = 'WITHDRAWN' THEN billing.accrual ELSE 0 END),0) AS WITHDRAWN
		FROM orders 
		JOIN billing ON orders.number = billing.order_number
		WHERE orders.user_id =$1  AND billing.status IN ('PROCESSED', 'WITHDRAWN');`

		var balance Balance

		err := tx.QueryRowContext(ctx, getBalanceQuery, userID).Scan(&balance.Current, &balance.Withdraw)

		if err != nil {
			return nil, err
		}
		balance.Current = balance.Current / 100
		balance.Withdraw = balance.Withdraw / 100

		if balance.Current-balance.Withdraw < orderSum.Sum {
			return nil, ErrNotEnoughFunds
		}

		// ======
		getOrderQuery := `SELECT number, user_id FROM orders
						  WHERE orders.number = $1`

		var orderUserID OrderUserID
		err = tx.QueryRowContext(ctx, getOrderQuery, orderSum.OrderNumber).Scan(&orderUserID.OrderNumber, &orderUserID.UserID)

		switch {
		case err == sql.ErrNoRows:
			t := time.Now()
			addOrderQuery := `INSERT INTO orders(number, user_id, uploaded_at) VALUES ($1, $2, $3)`
			_, err = tx.ExecContext(ctx, addOrderQuery, orderSum.OrderNumber, userID, t)
			if err != nil {
				return OrderUserID{}, err
			}
		case err != nil:
			return OrderUserID{}, fmt.Errorf("ошибка при получении ордера: %w", err)
		}

		if orderUserID.UserID != userID {
			return OrderUserID{}, fmt.Errorf("нельзя вывести деньги другому пользователю %s", orderUserID.UserID)
		}

		//======

		addOrderQuery := `INSERT INTO billing (order_number, status, accrual, uploaded_at, time)
		VALUES ($1, 'WITHDRAWN', $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
	`
		_, err = tx.ExecContext(ctx, addOrderQuery, orderSum.OrderNumber, orderSum.Sum*100)

		return nil, err
	}
}

func (storage *Storage) GetWithdrawals(ctx context.Context, userID string) dbOperation {
	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		queryWithdrawals := `SELECT orders.number, billing.accrual AS sum, billing.uploaded_at AS processed_at
			FROM orders
			JOIN billing ON orders.number = billing.order_number
			WHERE orders.user_id = $1
			AND billing.status = 'WITHDRAWN'
			ORDER BY billing.uploaded_at ';`

		rows, err := tx.QueryContext(ctx, queryWithdrawals, userID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var withdrawalsList []Withdrawal
		for rows.Next() {
			var w Withdrawal
			err := rows.Scan(&w.OrderNumber, &w.Sum, &w.ProcessedAt)
			if err != nil {
				return nil, err
			}
			w.Sum = w.Sum / 100
			withdrawalsList = append(withdrawalsList, w)
		}
		if err := rows.Err(); err != nil {
			return withdrawalsList, err
		}
		return withdrawalsList, nil
	}
}

func (storage *Storage) GetNewProcessedOrders(ctx context.Context) dbOperation {
	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		query := `
		SELECT b.order_number
		FROM billing b
		WHERE (b.order_number, b.time) IN (
			SELECT mt.order_number, MAX(mt.max_time) AS max_time
			FROM (
				SELECT order_number, MAX(time) AS max_time
				FROM billing
				GROUP BY order_number
			) AS mt
			WHERE b.order_number = mt.order_number
			GROUP BY mt.order_number
		) AND b.status IN ('PROCESSING', 'NEW') ;
		`
		var ordersList []string
		rows, err := tx.QueryContext(ctx, query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var orderNumber string

		for rows.Next() {

			err := rows.Scan(&orderNumber)
			if err != nil {
				return nil, err
			}
			ordersList = append(ordersList, orderNumber)
		}
		if err := rows.Err(); err != nil {
			return ordersList, err
		}
		return ordersList, nil
	}
}

func (storage *Storage) PutStatuses(ctx context.Context, orderStatus *[]OrderStatusNew) dbOperation {
	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		t := time.Now()
		builder := strings.Builder{}
		builder.WriteString("INSERT INTO billing (order_number, status, accrual, uploaded_at, time)\n")
		builder.WriteString("VALUES\n")
		for m, v := range *orderStatus {

			builder.WriteString(fmt.Sprintf("(%s,'%s',%v,%v,%s)", v.Number, v.Status, v.Accrual*100, "$1", "CURRENT_TIMESTAMP"))

			if m == len(*orderStatus)-1 {
				builder.WriteString("\n")
			} else {
				builder.WriteString(",\n")
			}

		}
		builder.WriteString("ON CONFLICT (order_number, status)\n")
		builder.WriteString("DO UPDATE SET order_number = EXCLUDED.order_number, status = EXCLUDED.status, accrual = EXCLUDED.accrual,  uploaded_at = EXCLUDED.uploaded_at;")
		query := builder.String()

		_, err := tx.ExecContext(ctx, query, t)

		// В целях отладки
		fmt.Println("В целях отладки", query)

		return OrderUserID{}, err
	}
}

// этот метод написан для тестирования
func (storage *Storage) GetBilling(ctx context.Context) dbOperation {
	return func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		query := `SELECT order_number, status, accrual AS accrual, uploaded_at, time FROM billing;`
		rows, err := tx.QueryContext(ctx, query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var billingList []Billing
		for rows.Next() {
			var b Billing
			err := rows.Scan(&b.OrderNumber, &b.Status, &b.Accrual, &b.UploadedAt, &b.Time)
			if err != nil {
				return nil, err
			}
			b.Accrual = b.Accrual / 100
			billingList = append(billingList, b)
		}
		if err := rows.Err(); err != nil {
			return billingList, err
		}
		return billingList, nil
	}
}
