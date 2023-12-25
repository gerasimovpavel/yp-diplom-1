package accruals

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/config"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/logger"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/model"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/storage"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"sync"
)

type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

func CheckAccruals(userID uuid.UUID) {
	ctx := context.Background()
	logger.Logger.Info("processing orders start...")
	orders, err := storage.Stor.ProcessingOrders(ctx, userID)
	if err != nil {
		logger.Logger.Sugar().Error(err)
		return
	}
	logger.Logger.Info("processing orders end...")
	wg := sync.WaitGroup{}
	logger.Logger.Info(fmt.Sprintf("orders count:%d", len(orders)))
	logger.Logger.Info("check order accruals start...")
	for _, ord := range orders {
		wg.Add(1)
		go func(order *model.Order) {
			accrual := Accrual{}

			defer wg.Done()
			httpc := resty.New().SetBaseURL(config.Options.AccrualSystemAddress)
			req := httpc.R().
				SetContext(context.Background()).SetResult(&accrual)

			resp, err := req.Get(fmt.Sprintf("/api/orders/%s", order.Number))
			body := resp.Body()
			if err != nil {
				logger.Logger.Error(err.Error())
			}
			if resp.StatusCode() != 200 {
				logger.Logger.Warn(fmt.Sprintf("status: %d body: %s", resp.StatusCode(), body))
			}
			accrualJSON, _ := json.Marshal(accrual)
			orderJSON, _ := json.Marshal(order)
			logger.Logger.Debug("MAGIC???")
			logger.Logger.Debug(string(accrualJSON))
			logger.Logger.Debug(string(orderJSON))

			if accrual.Order == order.Number && accrual.Status != order.Status {
				logger.Logger.Debug("MAGIC!!!")
				order.Accrual = accrual.Accrual
				order.Status = accrual.Status
				_, err = storage.Stor.SetOrder(ctx, order)
				if err != nil {
					logger.Logger.Error(err.Error())
				}
				err = storage.Stor.UpdateBalance(ctx, userID)
				if err != nil {
					logger.Logger.Error(err.Error())
				}
			}
		}(ord)

	}
	logger.Logger.Info("check order accruals end...")
	wg.Wait()
}
