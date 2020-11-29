package dbm

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Instance struct {
	db 				*sqlx.Conn
	tx 				*sqlx.Tx
	Closed			bool
}

// Create a new connection to database and start a transaction
func Conn (
	ctx context.Context,
	dbName string,
	driver string,
	user string,
	password string,
	host string,
	port string,
	) (*Instance, error) {
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, password, host, port, dbName,
	)

	db, err := sqlx.ConnectContext(ctx, driver, dataSourceName)
	if  err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if  err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(0)

	conn, err := db.Connx(ctx)
	if err != nil {
		return nil, err
	}

	i :=  &Instance{
		db: conn,
	}

	err = i.Begin(ctx, &sql.TxOptions{ Isolation: sql.LevelDefault })
	if err != nil {
		return nil, err
	}

	return i, nil
}

// Begins a transaction
func (i *Instance) Begin (ctx context.Context, opts *sql.TxOptions) (err error) {
	i.tx, err = i.db.BeginTxx(ctx, opts)
	return
}

// Commit a transaction
func (i *Instance) Commit () error {
	err := i.tx.Commit()
	if err == nil {
		i.tx = nil
	}
	return err
}

// Rollback a transaction
func (i *Instance) Rollback () error {
	err := i.tx.Rollback()
	if err == nil {
		i.tx = nil
	}
	return err
}

// Close current instance
func (i *Instance) Close () error {
	if i.Closed {
		return nil
	}

	err := i.db.Close()
	if err == nil {
		i.Closed = true
	}

	return err
}

// Get executes a query that is expected to return at most one row in a structure.
func (i *Instance) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return i.tx.GetContext(ctx, dest, query, args...)
}

// GetAll executes a query that is expected to return one or many rows in a array of structure.
func (i *Instance) GetAll(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return i.tx.SelectContext(ctx, dest, query, args...)
}

/*
	Query executes a query that returns rows, typically a SELECT.
	The args are for any placeholder parameters in the query.
*/
func (i *Instance) Query(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return i.tx.QueryxContext(ctx, query, args...)
}

/*
	QueryRow executes a query that is expected to return at most one row.
	QueryRow always returns a non-nil value. Errors are deferred until
	Row's Scan method is called.
	If the query selects no rows, the *sql.Row's Scan will return sql.ErrNoRows.
	Otherwise, the *sql.Row's Scan scans the first selected row and discards
	the rest.
*/
func (i *Instance) QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return i.tx.QueryRowxContext(ctx, query, args...)
}

/*
	Exec executes a query that doesn't return rows.
	For example: an INSERT and UPDATE.
*/
func (i *Instance) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return i.tx.ExecContext(ctx, query, args...)
}
