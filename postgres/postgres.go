package postgres

import (
	"context"
	"hash/fnv"
	"pg-upsert-returning/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type pgCustomerRepo struct {
	dbHandler *pgxpool.Pool
}

func NewPgCustomerRepo(connString string) (*pgCustomerRepo, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &MyTid{&pgtype.TID{}},
			Name:  "tid",
			OID:   pgtype.TIDOID,
		})
		return nil
	}

	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &pgCustomerRepo{
		dbHandler: db,
	}, nil
}

func (repo *pgCustomerRepo) UpsertCustomerCte(ctx context.Context, id uuid.UUID) (res domain.UpsertedRow, err error) {
	query :=
		"WITH " +
			"search AS (SELECT ctid, xmax, id FROM customers WHERE customer_id = $1 LIMIT 1)," +
			"add AS (INSERT INTO customers (customer_id) SELECT $1 WHERE NOT EXISTS(SELECT id from search) RETURNING ctid, xmax, id)" +
			"SELECT ctid, xmax, id from add	UNION ALL SELECT ctid, xmax, id from search"

	err = repo.dbHandler.QueryRow(ctx, query, id).Scan(&res.CTID, &res.XMAX, &res.ID)
	return
}

func (repo *pgCustomerRepo) UpsertCustomerDoNothing(ctx context.Context, id uuid.UUID) (res domain.UpsertedRow, err error) {

	batch := &pgx.Batch{}
	batch.Queue("INSERT INTO customers(customer_id) VALUES($1) ON CONFLICT DO NOTHING", id)
	batch.Queue("SELECT ctid, xmax, id FROM customers WHERE customer_id = $1", id)
	results := repo.dbHandler.SendBatch(ctx, batch)

	_, err = results.Exec()
	if err != nil {
		return
	}

	err = results.QueryRow().Scan(&res.CTID, &res.XMAX, &res.ID)
	return
}

func hash(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64())
}

func (repo *pgCustomerRepo) UpsertCustomerLock(ctx context.Context, id uuid.UUID) (res domain.UpsertedRow, err error) {
	tx, err := repo.dbHandler.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", hash(id.String()))
	if err != nil {
		return
	}

	err = tx.QueryRow(ctx, "SELECT ctid, xmax, id FROM customers WHERE customer_id = $1", id).Scan(&res.CTID, &res.XMAX, &res.ID)

	if err != nil && err == pgx.ErrNoRows {
		q := "INSERT INTO customers(customer_id) VALUES($1) RETURNING ctid, xmax, id"
		err = tx.QueryRow(ctx, q, id).Scan(&res.CTID, &res.XMAX, &res.ID)
	}

	if err == nil {
		err = tx.Commit(ctx)
	} else {
		err = tx.Rollback(ctx)
	}

	return
}

func (repo *pgCustomerRepo) UpsertCustomerConflict(ctx context.Context, id uuid.UUID) (res domain.UpsertedRow, err error) {
	query :=
		"INSERT INTO customers(customer_id) VALUES($1) " +
			"ON CONFLICT (customer_id) DO UPDATE SET customer_id = excluded.customer_id RETURNING ctid, xmax, id"

	err = repo.dbHandler.QueryRow(ctx, query, id).Scan(&res.CTID, &res.XMAX, &res.ID)
	return
}
