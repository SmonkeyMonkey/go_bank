package main

import (
	"math/rand"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}
type LoginRequest struct {
	Nubmer   int64  `json:"number"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Number int64  `json:"number"`
	Token  string `json:"token"`
}
type TrasnferRequest struct {
	FromAccount int `json:"fromAccount"`
	ToAccount   int `json:"toAccount"`
	Amount      int `json:"amount"`
}
type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Number            int64     `json:"number"`
	EncryptedPassword string    `json:"-"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encrPswrd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:                rand.Intn(10000),
		FirstName:         firstName,
		LastName:          lastName,
		Number:            int64(rand.Intn(10000000)),
		EncryptedPassword: string(encrPswrd),
		CreatedAt:         time.Now().UTC(),
	}, nil
}

func NewSeedAccount(number int,firstName, lastName, password string) (*Account, error) {
	encrPswrd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:                255,
		FirstName:         firstName,
		LastName:          lastName,
		Number:            int64(number),
		EncryptedPassword: string(encrPswrd),
		Balance:           5000,
		CreatedAt:         time.Now().UTC(),
	}, nil
}
func (a *Account) ValidatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
}

type JWTClaims struct {
	UserID int `json:"userid"`
	jwt.RegisteredClaims
}
