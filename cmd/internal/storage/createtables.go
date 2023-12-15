package storage

import (
	"context"
)

func createTables(w *PgWorker) error {
	//users
	ctx := context.Background()
	_, err := w.Exec(ctx, `

			CREATE TABLE public.accruals (
				"order" varchar(128) NOT NULL,
				summa float8 NULL
			);
			CREATE INDEX "order" ON public.accruals USING btree ("order" varchar_ops) WITH (deduplicate_items='true');
			
			CREATE TABLE public.balance (
				user_id uuid NULL,
				accrual numeric(19, 2) NULL DEFAULT 0,
				withdraw numeric(19, 2) NULL DEFAULT 0,
				"current" numeric(19, 2) NULL DEFAULT 0,
				CONSTRAINT user_id UNIQUE (user_id)
			);
			
			CREATE TABLE public.orders (
				"number" varchar(20) NULL,
				status varchar(50) NULL,
				accrual numeric(19, 2) NULL DEFAULT 0,
				uploaded_at timestamptz NULL,
				user_id uuid NULL,
				CONSTRAINT "number" UNIQUE (number)
			);
			
			CREATE TABLE public.users (
				user_id uuid NOT NULL,
				login varchar(128) NOT NULL,
				"password" varchar(128) NOT NULL,
				CONSTRAINT login UNIQUE (login),
				CONSTRAINT users_pkey PRIMARY KEY (user_id)
			);
			
			
			CREATE TABLE public.withdrawals (
				"order" varchar(128) NULL,
				summa numeric(19, 2) NULL,
				processed_at timestamptz NULL,
				user_id uuid NULL
			);
`)
	if err != nil {
		return err
	}
	return nil
}
