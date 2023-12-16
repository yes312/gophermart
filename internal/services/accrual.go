package services

import (
	"context"
	"encoding/json"
	"fmt"
	db "gophermart/internal/database"
	"gophermart/models"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

const numberOfWorkers = 3

type accrual struct {
	accrualSysremAdress      string
	accrualRequestInterval   int
	accuralPuttingDBInterval int
	storage                  db.StoragerDB
	logger                   *zap.SugaredLogger
}

func NewAccrual(accrualSysremAdress string, accrualRequestInterval int, accuralPuttingDBInterval int, storage db.StoragerDB, logger *zap.SugaredLogger) *accrual {
	return &accrual{
		accuralPuttingDBInterval: accuralPuttingDBInterval,
		accrualSysremAdress:      accrualSysremAdress,
		accrualRequestInterval:   accrualRequestInterval,
		storage:                  storage,
		logger:                   logger,
	}

}
func (a *accrual) RunAccrualRequester(ctx context.Context, wg *sync.WaitGroup) {

	orders := make(chan string, 1000)
	ordersFromAccrual := make(chan models.OrderStatusNew, 1000)

	wg.Add(1)
	go a.collectOrders(ctx, orders, wg)

	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go a.worker(ctx, orders, ordersFromAccrual, wg)
	}

	wg.Add(1)
	go a.putOrdersInDB(ctx, ordersFromAccrual, wg)

}

func (a *accrual) collectOrders(ctx context.Context, orders chan<- string, wg *sync.WaitGroup) {

	defer wg.Done()
	defer close(orders)
	ticker := time.NewTicker(time.Duration(a.accrualRequestInterval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:

			result, err := a.storage.WithRetry(ctx, a.storage.GetNewProcessedOrders(ctx))
			if err != nil {
				a.logger.Errorf("ошибка при получении новых заказов %w", err)
			}

			if res, ok := result.([]string); ok {
				for _, v := range res {
					// log.Println("нашли ордера на отправку в accural: ", v)
					orders <- v
				}
			}

		}
	}
}

func (a *accrual) worker(ctx context.Context, in chan string, out chan models.OrderStatusNew, wg *sync.WaitGroup) {
	defer wg.Done()

	client := resty.New()
	// url := fmt.Sprint(a.accrualSysremAdress, "/api/orders/")
	for {
		select {
		case <-ctx.Done():
			return
		case orderNumber, ok := <-in:
			if !ok {
				return
			}
			url := fmt.Sprint(a.accrualSysremAdress, "/api/orders/", orderNumber)
			// log.Println("адрес запроса: ", url)
			resp, err := client.R().
				SetContext(ctx).
				// SetHeader("Content-Type", "application/json").
				Get(url)

			if err != nil {
				a.logger.Errorf("ошибка при выполнении response: %w", err)
			} else {
				var order models.OrderStatusNew

				if resp.StatusCode() != 200 {
					a.logger.Errorf("wrong status code: %d order: %s", resp.StatusCode(), orderNumber)
					continue
				}
				// log.Println("BODY: ", string(resp.Body()))
				if err := json.Unmarshal(resp.Body(), &order); err != nil {
					a.logger.Errorf("Ошибка при декодировании JSON: %w", err)
					continue
				} else {
					// fmt.Println("получено из accrual: ", order)
					out <- order
				}

			}

		}
	}
}

func (a *accrual) putOrdersInDB(ctx context.Context, ordersFromAccrual chan models.OrderStatusNew, wg *sync.WaitGroup) {

	defer close(ordersFromAccrual)
	defer wg.Done()
	var ordersList []models.OrderStatusNew
	go func() {
		for v := range ordersFromAccrual {
			ordersList = append(ordersList, v)
			// log.Println("добавили в ordersList: ", v, ordersList)
		}
	}()

	ticker := time.NewTicker(time.Duration(a.accuralPuttingDBInterval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

			if len(ordersList) != 0 {

				var mutex sync.Mutex
				mutex.Lock()
				ordersListCopy := make([]models.OrderStatusNew, len(ordersList))
				copy(ordersListCopy, ordersList)
				ordersList = nil
				// log.Println("ordersListCopy: ,будем его сохранять в базе", ordersListCopy)
				mutex.Unlock()

				if _, err := a.storage.WithRetry(ctx, a.storage.PutStatuses(ctx, &ordersListCopy)); err != nil {
					a.logger.Error("ошибка при сохранении статусов", err)
				} else {
					a.logger.Info("статусы сохранены")
				}
			}

		}
	}

}
