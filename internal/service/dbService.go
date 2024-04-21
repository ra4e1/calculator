package service

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DbUser struct {
	ID       int64
	Name     string
	Password string
}

type DbExpression struct {
	ID        int64
	UserId    int64
	Answer    float64
	Ready     int
	Err       string
	Expresion string
}

type DbService struct {
	dbname  string
	db      *sql.DB
	context context.Context
}

func NewDbService(dbname string) *DbService { //создание
	return &DbService{
		dbname: dbname,
	}
}

func (srv *DbService) createTables() error {
	const usersTable = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		name TEXT UNIQUE,
		password TEXT
	);`

	const stateTable = `
	CREATE TABLE IF NOT EXISTS expressions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		ready   INTEGER NOT NULL, 
		answer Float,
		err TEXT,
		expression TEXT
	);`

	if _, err := srv.db.ExecContext(srv.context, usersTable); err != nil {
		return err
	}
	if _, err := srv.db.ExecContext(srv.context, stateTable); err != nil {
		return err
	}

	return nil
}

func (srv *DbService) InsertUser(user *DbUser) (int64, error) {
	var q = `
	INSERT INTO users (name, password) values ($1, $2)
	`
	result, err := srv.db.ExecContext(srv.context, q, user.Name, user.Password)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (srv *DbService) FindUserByName(name string) (*DbUser, error) {
	user := &DbUser{}
	var q = "SELECT id, name, password FROM users WHERE name=$1"
	err := srv.db.QueryRowContext(srv.context, q, name).Scan(&user.ID, &user.Name, &user.Password)
	return user, err
}

func (srv *DbService) InsertExpression(expression *DbExpression) (int64, error) {
	var q = `
	INSERT INTO expressions (user_id, ready, answer, err, expression) values ($1, $2, $3, $4, $5)
	`
	result, err := srv.db.ExecContext(
		srv.context,
		q,
		expression.UserId,
		expression.Ready,
		expression.Answer,
		expression.Err,
		expression.Expresion,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (srv *DbService) UpdateExpression(id int64, expression *DbExpression) error {
	var q = `
	UPDATE expressions SET 
		user_id = $1, 
		answer = $2, 
		ready = $3, 
		err = $4, 
		expression = $5 
	WHERE
	    id = $6
	`
	_, err := srv.db.ExecContext(srv.context, q,
		expression.UserId,
		expression.Answer,
		expression.Ready,
		expression.Err,
		expression.Expresion,
		id,
	)
	return err
}

func (srv *DbService) FindExpressionById(userId, id int64) (*DbExpression, error) {
	var q = `
	SELECT id, user_id, answer, ready, err, expression FROM expressions WHERE user_id = $1 AND id = $2
	`
	expression := &DbExpression{}
	err := srv.db.QueryRowContext(srv.context, q, userId, id).Scan(&expression.ID, &expression.UserId, &expression.Answer, &expression.Ready, &expression.Err, &expression.Expresion)
	return expression, err
}

func (srv *DbService) FindExpressions(userId int64) ([]*DbExpression, error) {
	var q = `
	SELECT id, user_id, answer, ready, err, expression FROM expressions WHERE user_id = $1
	`
	rows, err := srv.db.QueryContext(srv.context, q, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expressions []*DbExpression
	for rows.Next() {
		exp := &DbExpression{}
		if err := rows.Scan(&exp.ID, &exp.UserId, &exp.Answer, &exp.Ready, &exp.Err, &exp.Expresion); err != nil {
			return expressions, err
		}
		expressions = append(expressions, exp)
	}

	if err = rows.Err(); err != nil {
		return expressions, err
	}

	return expressions, nil
}

func (srv *DbService) Open() (err error) {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", srv.dbname)
	if err != nil {
		return err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}
	srv.context = ctx
	srv.db = db

	if err = srv.createTables(); err != nil {
		return err
	}

	return nil
}

func (srv *DbService) Close() {
	if srv.db != nil {
		srv.db.Close()
	}
}
