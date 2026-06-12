package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"subscription-vault-service/internal/clients/postgres/sqlc"
	"time"

	uuid "github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

const (
	insertOneTime  = 1000
	requestTimeout = time.Second * 5
)

var defaultTxOpt = &sql.TxOptions{Isolation: sql.LevelRepeatableRead}

//go:generate mockgen -source=postgres.go -destination=postgres_mock.go -package=postgres IDB,IQuerier

type IQuerier interface {
	sqlc.Querier
}

type IDB interface {
	ExecTx(*sql.TxOptions, func(context.Context, IQuerier) error) error
	Querier() IQuerier
	CtxWithCancel() (context.Context, context.CancelFunc)
}

type DB struct {
	ctx  context.Context
	conn *sql.DB
	sqlc *sqlc.Queries
}

type Client struct {
	db IDB
}

func NewSQLConn(ctx context.Context, user, password, host, port, dbname string) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable opening db connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	go func() {
		<-ctx.Done()
		err = db.Close()
		if err != nil {
			panic(err)
		}
	}()

	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)

	return db, nil
}

func NewClient(ctx context.Context, conn *sql.DB) *Client {
	return buildClient(&DB{
		ctx:  ctx,
		sqlc: sqlc.New(conn),
		conn: conn,
	})
}

func buildClient(db IDB) *Client {
	return &Client{
		db: db,
	}
}

func (db *DB) CtxWithCancel() (context.Context, context.CancelFunc) {
	return context.WithTimeout(db.ctx, requestTimeout)
}

func (db *DB) ExecTx(txOpt *sql.TxOptions, withTx func(context.Context, IQuerier) error) (err error) {
	ctx, cancel := db.CtxWithCancel()
	defer cancel()

	var tx *sql.Tx
	tx, err = db.conn.BeginTx(ctx, txOpt)
	if err != nil {
		return
	}

	defer func() {
		errRB := tx.Rollback()
		if errRB != nil && !errors.Is(errRB, sql.ErrTxDone) {
			if err != nil {
				err = fmt.Errorf("ExecTx error: %w; Rollback error: %w", err, errRB)
			} else {
				err = fmt.Errorf("rollback error: %w", errRB)
			}
		}
	}()

	if err = withTx(ctx, db.sqlc.WithTx(tx)); err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}

func (db *DB) Querier() IQuerier {
	return db.sqlc
}

func isDuplicate(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr != nil && pqErr.Code == pqerror.UniqueViolation
}

func isNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func nullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{
			Valid: false,
		}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func fromNullString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullUUID(u *uuid.UUID) uuid.NullUUID {
	if u == nil {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{
		UUID:  *u,
		Valid: true,
	}
}

func fromNullUUID(nu uuid.NullUUID) *uuid.UUID {
	if nu.Valid {
		return &nu.UUID
	}
	return nil
}

func nullInt64(n *int64) sql.NullInt64 {
	if n == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Valid: true, Int64: *n}
}

func nullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func fromNullTime(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
