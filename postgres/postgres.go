package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"pg-upsert-returning/domain"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type pgCustomerRepo struct {
	dbHandler *sql.DB
}

func NewPgCustomerRepo(connString string) (*pgCustomerRepo, error) {
	db, err := sql.Open("postgres", connString)
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

	row := repo.dbHandler.QueryRowContext(ctx, query, id)
	err = row.Scan(&res.CTID, &res.XMAX, &res.ID)
	return
}

func (repo *pgCustomerRepo) UpsertCustomerDoNothing(ctx context.Context, id uuid.UUID) (res domain.UpsertedRow, err error) {
	query :=
		"INSERT INTO customers(customer_id)	VALUES($1) ON CONFLICT DO NOTHING;" +
			"SELECT ctid, xmax, id FROM customers WHERE customer_id = $1"

	query = strings.ReplaceAll(query, "$1", fmt.Sprintf("'%s'", id.String()))
	row := repo.dbHandler.QueryRowContext(ctx, query)
	err = row.Scan(&res.CTID, &res.XMAX, &res.ID)
	return
}

func hash(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64())
}

func (repo *pgCustomerRepo) UpsertCustomerLock(ctx context.Context, id uuid.UUID) (res domain.UpsertedRow, err error) {
	tx, err := repo.dbHandler.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	_, err = tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock($1)", hash(id.String()))
	if err != nil {
		return
	}

	r := tx.QueryRowContext(ctx, "SELECT ctid, xmax, id FROM customers WHERE customer_id = $1", id)
	err = r.Scan(&res.CTID, &res.XMAX, &res.ID)

	if err != nil && err == sql.ErrNoRows {
		q := "INSERT INTO customers(customer_id) VALUES($1) RETURNING ctid, xmax, id"
		row := tx.QueryRowContext(ctx, q, id)
		err = row.Scan(&res.CTID, &res.XMAX, &res.ID)
	}

	if err == nil {
		err = tx.Commit()
	} else {
		err = tx.Rollback()
	}

	return
}

func (repo *pgCustomerRepo) UpsertCustomerConflict(ctx context.Context, id uuid.UUID) (res domain.UpsertedRow, err error) {
	query :=
		"INSERT INTO customers(customer_id) VALUES($1) " +
			"ON CONFLICT (customer_id) DO UPDATE SET customer_id = excluded.customer_id RETURNING ctid, xmax, id"

	row := repo.dbHandler.QueryRowContext(ctx, query, id)
	err = row.Scan(&res.CTID, &res.XMAX, &res.ID)
	return
}
