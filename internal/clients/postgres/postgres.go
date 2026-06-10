package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"subscription-vault-service/internal/clients/postgres/sqlc"
	sp "subscription-vault-service/internal/secret_provider"
	"subscription-vault-service/internal/supports"
	"time"
)

const (
	insertOneTime  = 1000
	requestTimeout = time.Second * 5

	db_host_secret_path     = "db_host"
	db_port_secret_path     = "db_port"
	db_password_secret_path = "db_password"
	db_user_secret_path     = "db_user"
	db_name_secret_path     = "db_name"
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

func NewSQLConn(ctx context.Context, sp *sp.SecretProvider) (*sql.DB, error) {
	host, err := sp.ReadSecret(db_host_secret_path)
	if err != nil {
		return nil, err
	}
	port, err := sp.ReadSecret(db_port_secret_path)
	if err != nil {
		return nil, err
	}
	user, err := sp.ReadSecret(db_user_secret_path)
	if err != nil {
		return nil, err
	}
	password, err := sp.ReadSecret(db_password_secret_path)
	if err != nil {
		return nil, err
	}
	dbname, err := sp.ReadSecret(db_name_secret_path)
	if err != nil {
		return nil, err
	}

	if !supports.IsInContainer() {
		host = "localhost"
	}

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
