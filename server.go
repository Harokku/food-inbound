package main

import (
	"database/sql"
	"fmt"
	"food-inbound/db"
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
	dbConn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	checkErrorAndPanic(err)

	defer dbConn.Close()

	err = dbConn.Ping()
	checkErrorAndPanic(err)
	fmt.Println("Correctly pinged DB")

	// Echo server definition
	e := echo.New()

	// Static endpoint to serve API doc
	// TODO: Write API doc
	e.Static("/docs", "static/docs")

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Pong")
	})
	// GET a single supplier based on ID
	e.GET("/suppliers/:id", func(c echo.Context) error {
		id := c.Param("id")
		supplier := db.Supplier{}
		db := db.Service{Db: dbConn}
		err := db.GetSupplier(&supplier, id)
		checkErrorAndPanic(err)
		return c.JSON(http.StatusOK, supplier)
	})

	// GET all records of Suppliers table
	e.GET("/suppliers", func(c echo.Context) error {
		var suppliers []db.Supplier
		db := db.Service{Db: dbConn}
		err := db.GetSuppliers(&suppliers)
		checkErrorAndPanic(err)
		return c.JSON(http.StatusOK, suppliers)
	})

	e.Logger.Fatal(e.Start(":" + port))
}

// Default error check with fatal if err != nil
func checkErrorAndPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
