package accruals

import (
	"context"
	"fmt"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/config"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/logger"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/model"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/storage"
	"github.com/go-resty/resty/v2"
	"sync"
)

func CheckAccruals() {
	ctx := context.Background()
	orders, err := storage.Stor.ProcessingOrders(ctx)
	if err != nil {
		logger.Logger.Sugar().Error(err)
		return
	}
	wg := sync.WaitGroup{}
	for _, o := range orders {
		wg.Add(1)
		go func(ord *model.Order) {
			var order model.Order
			var body string

			defer wg.Done()
			httpc := resty.New()

			req := httpc.R().
				SetContext(context.Background()).
				SetResult(&order).SetBody(body)

			resp, err := req.Get(fmt.Sprintf("%s/api/orders/%s", config.Options.AccrualSystemAddress, ord.Number))

			if err != nil {
				logger.Logger.Sugar().Error(err)
				return
			}
			if resp.StatusCode() != 200 {
				logger.Logger.Warn(fmt.Sprintf("status: %d body: %s", resp.StatusCode(), body))
				return
			}
			if order.Number == ord.Number && order.Status != ord.Status {
				_, err = storage.Stor.SetOrder(ctx, &order)
				if err != nil {
					logger.Logger.Sugar().Error(err)
					return
				}
			}
		}(o)

	}
	wg.Wait()
}
