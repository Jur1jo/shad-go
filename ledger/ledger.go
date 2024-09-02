//go:build !solution

package ledger

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type ledger struct {
	db *sql.DB
}

func (l *ledger) Close() error {
	return l.db.Close()
}

func New(ctx context.Context, dsn string) (Ledger, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.ExecContext(ctx, `
CREATE TABLE ledger (
    ID text UNIQUE,
    Money integer
)`)
	return &ledger{db: db}, err
}

func (l *ledger) CreateAccount(ctx context.Context, id ID) error {
	_, err := l.db.ExecContext(ctx, `INSERT INTO ledger (ID, Money) VALUES ($1, 0)`, id)
	return err
}

func (l *ledger) GetBalance(ctx context.Context, id ID) (Money, error) {
	var money Money
	err := l.db.QueryRowContext(ctx, `SELECT Money FROM ledger WHERE ID = $1`, id).Scan(&money)
	return money, err
}

func (l *ledger) Deposit(ctx context.Context, id ID, amount Money) error {
	if amount < 0 {
		return ErrNegativeAmount
	}
	res, err := l.db.ExecContext(ctx, `UPDATE ledger SET Money = Money + $1 WHERE ID = $2`, amount, id)
	if err != nil {
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return errors.New("no found id in database")
	}
	return nil
}

func (l *ledger) Withdraw(ctx context.Context, id ID, amount Money) error {
	if amount < 0 {
		return ErrNegativeAmount
	}
	var curMoney Money
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := tx.QueryRowContext(ctx, `SELECT money FROM ledger WHERE ID = $1 FOR UPDATE`, id).Scan(&curMoney); err != nil {
		return err
	}
	if curMoney-amount < 0 {
		return ErrNoMoney
	}
	_, err = tx.ExecContext(ctx, `UPDATE ledger SET Money = Money - $1 WHERE ID = $2`, amount, id)
	return tx.Commit()
}

func (l *ledger) Transfer(ctx context.Context, from, to ID, amount Money) error {
	if amount < 0 {
		return ErrNegativeAmount
	}
	var amountFrom Money
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if from < to {
		if err := tx.QueryRowContext(ctx, `SELECT money FROM ledger WHERE ID = $1 FOR UPDATE`, from).Scan(&amountFrom); err != nil {
			return err
		}
		var tmp Money
		if err := tx.QueryRowContext(ctx, `SELECT money FROM ledger WHERE ID = $1 FOR UPDATE `, to).Scan(&tmp); err != nil {
			return err
		}
	} else {
		if _, err := tx.ExecContext(ctx, `SELECT money FROM ledger WHERE ID = $1 FOR UPDATE `, to); err != nil {
			return err
		}
		if err := tx.QueryRowContext(ctx, `SELECT money FROM ledger WHERE ID = $1 FOR UPDATE`, from).Scan(&amountFrom); err != nil {
			return err
		}
	}
	if amountFrom-amount < 0 {
		return ErrNoMoney
	}
	if _, err := tx.ExecContext(ctx, `UPDATE ledger SET Money = Money - $1 WHERE ID = $2`, amount, from); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE ledger SET Money = Money + $1 WHERE ID = $2`, amount, to); err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}
