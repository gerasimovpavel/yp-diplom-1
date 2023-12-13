package storage

import (
	"context"
)

func createTables(w *PgWorker) error {
	//users
	ctx := context.Background()
	_, err := w.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS public.users
		(
			user_id uuid NOT NULL,
			login character varying(128) COLLATE pg_catalog."default" NOT NULL,
			password character varying(128) COLLATE pg_catalog."default" NOT NULL,
			CONSTRAINT users_pkey PRIMARY KEY (user_id),
			CONSTRAINT login UNIQUE (login)
		)
	`)
	if err != nil {
		return err
	}
	//orders
	_, err = w.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS public.orders
			(
				"number" character varying(20) COLLATE pg_catalog."sdefault",
				status character varying(50) COLLATE pg_catalog."default",
				accrual double precision DEFAULT 0,
				uploaded_at timestamp with time zone,
				user_id uuid,
				CONSTRAINT number UNIQUE ("number")
			)
    `)

	if err != nil {
		return err
	}
	//accruals
	_, err = w.Exec(ctx, `


CREATE TABLE IF NOT EXISTS public.accruals
(
    "order" character varying(128) COLLATE pg_catalog."default" NOT NULL,
    summa double precision
)

TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS "order"
    ON public.accruals USING btree
    ("order" COLLATE pg_catalog."default" varchar_ops ASC NULLS LAST)
    WITH (deduplicate_items=True)
    TABLESPACE pg_default;
			`)
	if err != nil {
		return err
	}

	_, err = w.Exec(ctx, `


CREATE TABLE IF NOT EXISTS public.withdrawals
(
    "order" character varying(128) COLLATE pg_catalog."default" NOT NULL,
    summa double precision
)

TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS "order"
    ON public.withdrawals USING btree
    ("order" COLLATE pg_catalog."default" varchar_ops ASC NULLS LAST)
    WITH (deduplicate_items=True)
    TABLESPACE pg_default;
			`)
	if err != nil {
		return err
	}
	_, err = w.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS public.balance
		(
			user_id uuid,
			accrual double precision,
			withdraw double precision,
			CONSTRAINT user_id UNIQUE (user_id)
		)
		`)
	if err != nil {
		return err
	}
	return nil
}
