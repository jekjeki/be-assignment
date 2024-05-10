package database

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func ConnectDatabase() {
	host := "localhost"
	port, _ := strconv.Atoi("5432")
	user := "postgres"
	db_name := "concrete"
	pass := "postgres123"

	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, db_name, pass)
	db, errSql := sql.Open("postgres", psqlSetup)
	if errSql != nil {
		fmt.Println("There is an error while connecting to the database ", errSql)
		fmt.Println(errSql.Error())
	} else {
		Db = db
		fmt.Println("Successfully connected to database!")
	}
}
