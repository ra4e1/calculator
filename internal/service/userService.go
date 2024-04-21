package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *DbService
}

func NewUserService(db *DbService) *UserService { //создание
	return &UserService{
		db: db,
	}
}

func generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}

func (u *UserService) Register(login string, password string) (int64, error) { // регистрация нового пользователя
	passwordHash, err := generate(password)
	if err != nil {
		return 0, err
	}

	user := &DbUser{Name: login, Password: passwordHash}
	id, err := u.db.InsertUser(user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *UserService) Login(login string, password string) (string, error) { // регистрация нового пользователя
	loginErr := fmt.Errorf("Login failed")
	dbuser, err := u.db.FindUserByName(login)
	if err != nil {
		return "", loginErr
	}
	if compare(dbuser.Password, password) != nil {
		return "", loginErr
	}

	auth := NewAuthService()
	token, err := auth.CreateToken(dbuser)
	if err != nil {
		return "", loginErr
	}
	return token, nil
}
