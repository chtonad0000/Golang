//go:build !solution

package dao

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type MyDao struct {
	DataBase *sql.DB
}

func CreateDao(ctx context.Context, dsn string) (Dao, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	myDao := &MyDao{
		DataBase: db,
	}
	_, err = db.ExecContext(ctx, "CREATE TABLE users (id SERIAL PRIMARY KEY,name VARCHAR(255))")
	if err != nil {
		return nil, err
	}
	return myDao, nil
}

func (myDao *MyDao) Create(ctx context.Context, u *User) (UserID, error) {
	var id int
	err := myDao.DataBase.QueryRowContext(ctx, "INSERT INTO users (name) VALUES ($1) RETURNING id", u.Name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return UserID(id), nil
}
func (myDao *MyDao) Update(ctx context.Context, u *User) error {
	res, err := myDao.DataBase.ExecContext(ctx, "UPDATE users SET name = $1 WHERE id = $2", u.Name, u.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (myDao *MyDao) Delete(ctx context.Context, id UserID) error {
	_, err := myDao.DataBase.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
func (myDao *MyDao) Lookup(ctx context.Context, id UserID) (User, error) {
	var userID int
	var name string
	err := myDao.DataBase.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id = $1", id).Scan(&userID, &name)
	if err != nil {
		return User{}, err
	}
	return User{ID: UserID(userID), Name: name}, nil
}
func (myDao *MyDao) List(ctx context.Context) ([]User, error) {
	rows, err := myDao.DataBase.QueryContext(ctx, "SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		errScan := rows.Scan(&user.ID, &user.Name)
		if errScan != nil {
			return nil, errScan
		}
		users = append(users, user)
	}
	return users, nil
}
func (myDao *MyDao) Close() error {
	return myDao.DataBase.Close()
}
