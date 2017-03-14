package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lib/pq"
)

// User models database user
type User struct {
	ID       int
	Username string
	Age      int
}

const (
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbHost     = "localhost"
	dbPort     = "5432"
	dbName     = "test"
)

func main() {

	dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	fmt.Println("# Creating table")
	err = startDB(db)
	checkErr(err)

	fmt.Println("# Inserting values")
	err = loadData(db)
	checkErr(err)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM userinfo")
		checkErr(err)
		users := []*User{}
		for rows.Next() {
			user := &User{}
			err = rows.Scan(&user.ID, &user.Username, &user.Age)
			checkErr(err)
			users = append(users, user)
		}
		json.NewEncoder(w).Encode(users)
	})

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func startDB(db *sql.DB) error {
	query := `CREATE TEMP TABLE "userinfo" (
						  "uid" SERIAL PRIMARY KEY,
						  "username" varchar(30) UNIQUE NOT NULL,
							"age" integer
						)`
	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

func loadData(db *sql.DB) error {
	users := []User{
		User{
			Username: "demo1", Age: 19,
		},
		User{
			Username: "demo2", Age: 23,
		},
		User{
			Username: "demo3", Age: 22,
		},
		User{
			Username: "demo4", Age: 33,
		},
		User{
			Username: "demo5", Age: 43,
		},
		User{
			Username: "demo6", Age: 24,
		},
	}
	txn, err := db.Begin()
	if err != nil {
		checkErr(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn("userinfo", "username", "age"))
	if err != nil {
		checkErr(err)
	}

	for _, user := range users {
		_, err = stmt.Exec(user.Username, int64(user.Age))
		if err != nil {
			checkErr(err)
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		checkErr(err)
	}

	err = stmt.Close()
	if err != nil {
		checkErr(err)
	}

	err = txn.Commit()
	if err != nil {
		checkErr(err)
	}
	if err != nil {
		return err
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
