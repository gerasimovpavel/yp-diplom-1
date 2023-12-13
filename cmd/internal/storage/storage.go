package storage

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/model"
	"github.com/google/uuid"
	"time"
)

var Stor *PgStorage

type PgStorage struct {
	w *model.PgWorker
}

func NewPgStorage() (*PgStorage, error) {
	worker, err := model.NewPgWorker()
	if err != nil {
		return nil, err
	}

	err = createTables(worker)

	if err != nil {
		return nil, err
	}

	return &PgStorage{w: worker}, nil
}

func (pw *PgStorage) CreateUser(a *model.User) (*model.User, error) {

	sqlString := `INSERT INTO users (user_id, login, password) VALUES ($1,$2,$3)`

	userId := uuid.New()

	_, err := pw.w.Exec(context.Background(), sqlString, userId.String(), a.Login, a.PasswordHash())

	if err != nil {
		return a, err
	}
	a.UserID = userId.String()
	return a, nil
}

func (pw *PgStorage) GetUser(a *model.User) (*model.User, error) {
	users := []*model.User{}
	sqlString := `SELECT * FROM users WHERE login=$1 AND password=$2`

	err := pgxscan.Select(context.Background(), pw.w, &users, sqlString, a.Login, a.PasswordHash())

	if err != nil {
		return a, err
	}
	return users[0], nil
}

func (pw *PgStorage) GetOrder(number string) (*model.Order, error) {
	orders := []*model.Order{}
	sqlString := `SELECT * FROM orders WHERE number=$1`

	err := pgxscan.Select(context.Background(), pw.w, &orders, sqlString, number)

	if err != nil {
		return &model.Order{}, err
	}
	if len(orders) == 0 {
		return &model.Order{}, err
	}
	return orders[0], err
}

func (pw *PgStorage) GetOrderByUser(userId uuid.UUID) ([]*model.Order, error) {
	orders := []*model.Order{}
	sqlString := `SELECT * FROM orders WHERE user_id=$1`

	err := pw.w.Select(context.Background(), &orders, sqlString, userId)

	return orders, err
}

func (pw *PgStorage) SetOrders(o *model.Order) (*model.Order, error) {
	orders := []*model.Order{}
	sqlString := `INSERT INTO orders (number, user_id, status, uploaded_at) 
				  VALUES ($1,$2,$3,$4) 
				  ON CONFLICT (number) DO 
				    UPDATE SET status=excluded.status 
				    RETURNING number, user_id, status, uploaded_at`

	err := pw.w.Select(context.Background(), &orders, sqlString, o.Number, o.UserID, "NEW", o.UploadedAt)

	if err != nil {
		return o, err
	}

	if len(orders) == 0 {
		return o, nil
	}
	return orders[0], nil
}

func (pw *PgStorage) GetBalance(userId uuid.UUID) (*model.Balance, error) {
	balance := &model.Balance{}
	balances := []*model.Balance{}

	err := pw.w.Select(context.Background(), &balances,
		`
			SELECT accruals, withdraw FROM balance user_id=$1)
		`,
		userId)
	if err != nil {
		return balance, err
	}
	return balance, nil
}

func (pw *PgStorage) UpdateBalance(userId uuid.UUID) error {
	ctx := context.Background()
	err := pw.w.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = pw.w.Exec(ctx,
		`
			INSERT INTO balance (user_id, accrual)
				SELECT  o.user_id,
						SUM(a.summa) AS accrual
				FROM accruals a
				INNER JOIN orders O ON a.order = o.number
				WHERE o.user_id=$1
				GROUP BY o.user_id) AS ua
			ON CONFLICT (user_id) DO UPDATE SET accrual = excluded.accrual	
		`,
		userId)
	if err != nil {
		pw.w.Rollback(ctx)
		return err
	}

	_, err = pw.w.Exec(ctx,
		`
			INSERT INTO balance (user_id, accrual)
				SELECT  o.user_id,
						SUM(a.summa) AS withdraw
				FROM withdrawals a
				INNER JOIN orders O ON a.order = o.number
				WHERE o.user_id=$1
				GROUP BY o.user_id) AS ua
			ON CONFLICT (user_id) DO UPDATE SET accrual = excluded.accrual	
		`,
		userId)
	if err != nil {
		pw.w.Rollback(ctx)
		return err
	}
	err = pw.w.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (pw *PgStorage) SetWithdraw(w *model.Withdraw) (*model.Withdraw, error) {
	ctx := context.Background()

	pw.w.Begin(ctx)

	o, err := pw.GetOrder(w.Order)
	if err != nil {
		pw.w.Rollback(ctx)
		return w, err
	}
	t := time.Now()
	_, err = pw.w.Exec(ctx,
		`INSERT INTO withdrawals (order, summa, processed_at) VALUES ($1, $2, $3)`,
		w.Order, w.Sum, t,
	)
	if err != nil {
		pw.w.Rollback(ctx)
		return w, err
	}
	w.ProcessedAt = t

	err = pw.UpdateBalance(o.UserID)

	if err != nil {
		pw.w.Rollback(ctx)
		return w, err
	}

	pw.w.Commit(ctx)
	w.ProcessedAt = t
	return w, nil
}

func (pw *PgStorage) GetWithdrawals(userId uuid.UUID) ([]*model.Withdraw, error) {
	var w []*model.Withdraw
	err := pw.w.Select(context.Background(),
		&w,
		`SELECT * FROM withdrawals WHERE user_id=$1`,
		userId)
	if err != nil {
		return w, err
	}
	return w, nil
}
