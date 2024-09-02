//go:build !solution

package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type MyDao struct {
	db *sql.DB
}

func CreateDao(ctx context.Context, dsn string) (Dao, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(ctx, `
CREATE TABLE dao (
    ID SERIAL PRIMARY KEY,
    Name text
)`)
	if err != nil {
		return &MyDao{db: db}, err
	}
	return &MyDao{db: db}, nil
}

func (d *MyDao) Create(ctx context.Context, u *User) (UserID, error) {
	var id int
	err := d.db.QueryRowContext(ctx, `INSERT INTO dao (Name) VALUES ($1) RETURNING ID`, u.Name).Scan(&id)
	if err != nil {
		return 0, err
	}
	fmt.Println(u.Name, id)
	return UserID(id), nil
}

func (d *MyDao) Update(ctx context.Context, u *User) error {
	res, err := d.db.ExecContext(ctx, `UPDATE dao SET Name = $1 WHERE ID = $2`, u.Name, u.ID)
	if err != nil {
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return errors.New("no found id")
	}
	return err
}

func (d *MyDao) Delete(ctx context.Context, id UserID) error {
	_, err := d.db.ExecContext(ctx, `DELETE FROM dao WHERE ID = $1`, id)
	return err
}

func (d *MyDao) Lookup(ctx context.Context, id UserID) (User, error) {
	user := User{ID: id}
	err := d.db.QueryRowContext(ctx, `SELECT Name FROM dao WHERE ID = $1`, id).Scan(&user.Name)
	return user, err
}

func (d *MyDao) List(ctx context.Context) ([]User, error) {
	rows, err := d.db.QueryContext(ctx, `SELECT Name, ID FROM dao`)
	if err != nil {
		return nil, err
	}
	var users []User
	for rows.Next() {
		users = append(users, User{})
		err := rows.Scan(&users[len(users)-1].Name, &users[len(users)-1].ID)
		if err != nil {
			return users, err
		}
	}
	if err = rows.Err(); err != nil {
		return users, nil
	}
	return users, nil
}

func (d *MyDao) Close() error {
	err := d.db.Close()
	return err
}
