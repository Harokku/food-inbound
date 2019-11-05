package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	// Heroku port from env variable
	port := os.Getenv("PORT")

	// Heroku Postgres connection and ping
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	checkErrorAndPanic(err)

	defer db.Close()

	err = db.Ping()
	checkErrorAndPanic(err)
	fmt.Println("Correctly pinged DB")

	// TODO: Cancel this aftre functionality check on Heroku
	// Test select to ckeck for DB connection
	sqlStatement := `SELECT id, name, address FROM suppliers`
	var id, name, address string
	row := db.QueryRow(sqlStatement)
	switch err := row.Scan(&id, &name, &address); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println(id, name, address)
	default:
		panic(err)
	}

	//TODO: Line to check if HEROKU db path exist, clean after check
	fmt.Printf("DBURL: %s", os.Getenv("DATABASE_URL"))
	e := echo.New()

	e.Static("/docs", "static/docs")

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Pong")
	})
	// TODO: Return sample DB row to check for Heroku postgres
	e.GET("/db", func(c echo.Context) error {
		// TODO: Cancel this aftre functionality check on Heroku
		// Test select to ckeck for DB connection
		sqlStatement := `SELECT id, name, address FROM suppliers`
		var id, name, address string
		row := db.QueryRow(sqlStatement)
		switch err := row.Scan(&id, &name, &address); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return c.String(http.StatusBadRequest, "No rows were found!")
		case nil:
			returnString := id + " - " + name + " - " + address
			return c.String(http.StatusOK, returnString)
		default:
			return c.String(http.StatusBadRequest, "No rows were found!")
		}

	})

	e.Logger.Fatal(e.Start(":" + port))
}

func checkErrorAndPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
