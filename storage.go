package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
	Transfer(int, int, int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=go_bank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}
func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}
func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account(
		id serial primary key,
		first_name varchar(100),
		last_name varchar(100),
		number serial,
		encrypted_password varchar(100), 
		balance serial,
		created_at timestamp
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(account *Account) error {
	query := `insert into account 
	(first_name, last_name, number,encrypted_password, balance, created_at)
	values ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(
		query,
		account.FirstName,
		account.LastName,
		account.Number,
		account.EncryptedPassword,
		account.Balance,
		account.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("delete from account where id = $1", id)
	return err
}
func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}
func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from account")
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("select * from account where id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found ", id)
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query("select * from account where number = $1", number)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account with number %d not found ", number)
}

func (s *PostgresStore) Transfer(from, to, amount int) error {
	ctx := context.Background()
	currentBalance, err := s.GetBalanceByNumber(from)
	if err != nil {
		return err
	}
	if currentBalance < amount {
		return fmt.Errorf("Insufficient funds on the account")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("update account set balance = balance - $1 where number = $2", amount, from)
	if err != nil {
		tx.Rollback()
		log.Println("transaction rollback")
		return err
	}
	_, err = s.db.Exec("update account set balance = balance + $1 where number = $2", amount, to)
	if err != nil {
		tx.Rollback()
		log.Println("transaction rollback")
		return err

	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
func(s *PostgresStore) GetBalanceByNumber(number int) (int,error){
	var currentBalance int

	row,err := s.db.Query("select balance from account where number = $1",number)
	if err !=nil {
		return 0,err
	}
	for row.Next(){
		row.Scan(&currentBalance)
	}

	return currentBalance,nil
} 
func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt)
	return account, err
}
