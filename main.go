package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func init() {
	os.Setenv("JWT_SECRET", "secretjwt")
}

func main() {
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("Seeding the database")

		seedAccounts(store)
	}
	// fmt.Printf("%+v\n",store)
	server := NewAPIServer(":8080", store)

	server.Run()
}

func seedAccount(store Storage, firstname, lastName, pw string) *Account {
	acc, err := NewSeedAccount(7625482,"John", "Johnson", "secret")
	if err != nil {
		log.Fatal(err)
	}
	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	
	return acc
}

func seedAccounts(s Storage) {
	acc := seedAccount(s, "test", "user", "secret")
	fmt.Println("seed acc:", acc)
}
