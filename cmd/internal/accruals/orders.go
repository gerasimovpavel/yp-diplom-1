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

type accrual struct {
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
	for _, o := range orders {
		wg.Add(1)
		go func(ord *model.Order) {
			var a accrual

			defer wg.Done()
			httpc := resty.New().SetBaseURL(config.Options.AccrualSystemAddress)
			req := httpc.R().
				SetContext(context.Background()).SetResult(&a)

			resp, err := req.Get(fmt.Sprintf("/api/orders/%s", ord.Number))
			body := resp.Body()
			if err != nil {
				logger.Logger.Error(err.Error())
			}
			if resp.StatusCode() != 200 {
				logger.Logger.Warn(fmt.Sprintf("status: %d body: %s", resp.StatusCode(), body))
			}
			o1, _ := json.Marshal(a)
			o2, _ := json.Marshal(ord)
			logger.Logger.Debug("MAGIC???")
			logger.Logger.Debug(string(o1))
			logger.Logger.Debug(string(o2))

			if a.Order == ord.Number && a.Status != ord.Status {
				logger.Logger.Debug("MAGIC!!!")
				ord.Accrual = a.Accrual
				ord.Status = a.Status
				_, err = storage.Stor.SetOrder(ctx, ord)
				if err != nil {
					logger.Logger.Error(err.Error())
				}
				err = storage.Stor.UpdateBalance(ctx, userID)
				if err != nil {
					logger.Logger.Error(err.Error())
				}
			}
		}(o)

	}
	logger.Logger.Info("check order accruals end...")
	wg.Wait()
}
