//go:build !solution

package ledger

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type myLedger struct {
	dataBase *sql.DB
}

func New(ctx context.Context, dsn string) (Ledger, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	ledger := &myLedger{db}
	_, err = db.ExecContext(ctx, "CREATE TABLE accounts (id VARCHAR(100) PRIMARY KEY,balance BIGINT)")
	if err != nil {
		return nil, err
	}
	return ledger, nil
}

func (ledger *myLedger) CreateAccount(ctx context.Context, id ID) error {
	_, err := ledger.dataBase.ExecContext(ctx, "INSERT INTO accounts (id,balance) VALUES ($1,$2)", id, 0)
	if err != nil {
		return err
	}
	return nil
}
func (ledger *myLedger) GetBalance(ctx context.Context, id ID) (Money, error) {
	var money int64
	err := ledger.dataBase.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id=$1", id).Scan(&money)
	if err != nil {
		return -1, err
	}
	return Money(money), nil
}
func (ledger *myLedger) Deposit(ctx context.Context, id ID, amount Money) error {
	if amount < 0 {
		return ErrNegativeAmount
	}
	tx, err := ledger.dataBase.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		errRB := tx.Rollback()
		if errRB != nil && !errors.Is(sql.ErrTxDone, errRB) {
			return
		}
	}(tx)
	var res sql.Result
	res, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id=$2", amount, id)
	if err != nil {
		return err
	}
	var rowsAffected int64
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return tx.Commit()
}
func (ledger *myLedger) Withdraw(ctx context.Context, id ID, amount Money) error {
	if amount < 0 {
		return ErrNegativeAmount
	}

	tx, err := ledger.dataBase.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		errRB := tx.Rollback()
		if errRB != nil && !errors.Is(sql.ErrTxDone, errRB) {
			return
		}
	}(tx)
	var balance int64
	err = tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id=$1 FOR UPDATE", id).Scan(&balance)
	if err != nil {
		return err
	}
	if balance < int64(amount) {
		return ErrNoMoney
	}
	balance -= int64(amount)
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = $1 WHERE id = $2", balance, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}
func (ledger *myLedger) Transfer(ctx context.Context, from, to ID, amount Money) error {
	if amount < 0 {
		return ErrNegativeAmount
	}
	tx, err := ledger.dataBase.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		errRB := tx.Rollback()
		if errRB != nil && !errors.Is(sql.ErrTxDone, errRB) {
			return
		}
	}(tx)

	minSwapped := from
	maxSwapped := to
	if from > to {
		minSwapped = to
		maxSwapped = from
	}
	var balance int64
	err = tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id=$1 FOR UPDATE", minSwapped).Scan(&balance)
	if err != nil {
		return err
	}
	var balance2 int64
	err = tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id=$1 FOR UPDATE", maxSwapped).Scan(&balance2)
	if err != nil {
		return err
	}
	if minSwapped == from {
		if balance < int64(amount) {
			return ErrNoMoney
		}
		_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, from)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, to)
		if err != nil {
			return err
		}
	} else {
		if balance2 < int64(amount) {
			return ErrNoMoney
		}
		_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, to)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, from)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
func (ledger *myLedger) Close() error {
	return ledger.dataBase.Close()
}
