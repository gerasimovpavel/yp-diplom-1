package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/model"
	"github.com/google/uuid"
	"time"
)

var Stor *PgStorage

type PgStorage struct {
	w *PgWorker
}

func NewPgStorage() (*PgStorage, error) {
	worker, err := NewPgWorker()
	if err != nil {
		return nil, err
	}

	err = createTables(worker)

	if err != nil {
		return nil, err
	}

	return &PgStorage{w: worker}, nil
}

func (pw *PgStorage) CreateUser(ctx context.Context, a *model.User) (*model.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	sqlString := `INSERT INTO users (user_id, login, password) VALUES ($1,$2,$3)`

	userID := uuid.New()

	_, err := pw.w.Exec(ctx, sqlString, userID.String(), a.Login, a.PasswordHash())

	if err != nil {
		return a, err
	}
	a.UserID = userID.String()
	return a, nil
}

func (pw *PgStorage) GetUser(ctx context.Context, a *model.User) (*model.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	users := []*model.User{}
	sqlString := `SELECT * FROM users WHERE login=$1 AND password=$2`

	err := pgxscan.Select(ctx, pw.w, &users, sqlString, a.Login, a.PasswordHash())

	if err != nil {
		return a, err
	}
	return users[0], nil
}

func (pw *PgStorage) GetOrder(ctx context.Context, number string) (*model.Order, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	orders := []*model.Order{}
	sqlString := `SELECT * FROM orders WHERE number=$1`

	err := pw.w.Select(ctx, &orders, sqlString, number)

	if err != nil {
		return &model.Order{}, err
	}
	if len(orders) == 0 {
		return &model.Order{}, err
	}
	return orders[0], err
}

func (pw *PgStorage) GetOrderByUser(ctx context.Context, userID uuid.UUID) ([]*model.Order, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	orders := []*model.Order{}
	sqlString := `SELECT * FROM orders WHERE user_id=$1`

	err := pw.w.Select(ctx, &orders, sqlString, userID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return orders, err
	}
	return orders, nil
}

func (pw *PgStorage) SetOrder(ctx context.Context, o *model.Order) (*model.Order, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	orders := []*model.Order{}
	sqlString := `INSERT INTO orders (number, user_id, status, uploaded_at) 
				  VALUES ($1,$2,$3,$4) 
				  ON CONFLICT (number) DO 
				    UPDATE SET status=excluded.status 
				    RETURNING number, user_id, status, uploaded_at`

	err := pw.w.Select(ctx, &orders, sqlString, o.Number, o.UserID, "NEW", o.UploadedAt)

	if err != nil {
		return o, err
	}

	if len(orders) == 0 {
		return o, nil
	}
	return orders[0], nil
}

func (pw *PgStorage) GetBalance(ctx context.Context, userID uuid.UUID) (*model.Balance, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	balance := &model.Balance{}
	balances := []*model.Balance{}

	err := pw.w.Select(ctx, &balances,
		`
			SELECT accruals, withdraw FROM balance user_id=$1)
		`,
		userID)
	if err != nil {
		return balance, err
	}
	return balance, nil
}

func (pw *PgStorage) UpdateBalance(ctx context.Context, userID uuid.UUID) error {
	var err error

	if ctx == nil {
		ctx = context.Background()
	}

	allowCommit := true
	if ctx.Value("tx") == nil {
		allowCommit = false
		ctx, err = pw.w.Begin(ctx)
		if err != nil {
			return err
		}
	}

	_, err = pw.w.Exec(ctx,
		`
			INSERT INTO balance (user_id, accrual)
				SELECT  o.user_id,
						SUM(a.summa) AS accrual
				FROM accruals a
				INNER JOIN orders O ON a.order = o.number
				WHERE o.user_id=$1
				GROUP BY o.user_id
			ON CONFLICT (user_id) DO UPDATE SET accrual = excluded.accrual	
		`,
		userID)
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
				GROUP BY o.user_id
			ON CONFLICT (user_id) DO UPDATE SET accrual = excluded.accrual	
		`,
		userID)

	if err != nil {
		if allowCommit {
			pw.w.Rollback(ctx)
		}
		return err
	}
	if allowCommit {
		pw.w.Commit(ctx)
	}
	return nil
}

func (pw *PgStorage) SetWithdraw(ctx context.Context, w *model.Withdraw) (*model.Withdraw, error) {
	var err error
	if ctx == nil {
		ctx = context.Background()
	}

	allowCommit := true
	if ctx.Value("tx") == nil {
		allowCommit = false
		ctx, err = pw.w.Begin(ctx)
		if err != nil {
			return w, err
		}
	}

	o, err := pw.GetOrder(ctx, w.Order)
	if err != nil {
		if allowCommit {
			pw.w.Rollback(ctx)
		}
		return w, err
	}
	t := time.Now()
	_, err = pw.w.Exec(ctx,
		`INSERT INTO withdrawals ("order", summa, processed_at) VALUES ($1, $2, $3)`,
		w.Order, w.Sum, t,
	)
	if err != nil {
		pw.w.Rollback(ctx)
		return w, err
	}
	w.ProcessedAt = t

	err = pw.UpdateBalance(ctx, o.UserID)

	if err != nil {
		if allowCommit {
			pw.w.Rollback(ctx)
		}
		return w, err
	}
	if allowCommit {
		pw.w.Commit(ctx)
	}
	w.ProcessedAt = t
	return w, nil
}

func (pw *PgStorage) GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]*model.Withdraw, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var w []*model.Withdraw
	err := pw.w.Select(ctx,
		&w,
		`SELECT * FROM withdrawals WHERE user_id=$1`,
		userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return w, err
	}
	return w, nil
}

func (pw *PgStorage) ProcessingOrders(ctx context.Context) ([]*model.Order, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var o []*model.Order
	err := pw.w.Select(ctx,
		&o,
		`SELECT * FROM orders WHERE status IN ('NEW', 'PROCESSING', 'REGISTERED')`,
	)
	if err != nil {
		return o, err
	}
	return o, nil
}
