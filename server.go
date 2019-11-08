package main

import (
	"database/sql"
	"fmt"
	"food-inbound/api"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	// Heroku port from env variable
	port := os.Getenv("PORT")

	// -----------------------
	// Database connection config
	// -----------------------

	// Heroku Postgres connection and ping
	dbConn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	checkErrorAndPanic(err)

	defer dbConn.Close()

	err = dbConn.Ping()
	checkErrorAndPanic(err)
	fmt.Println("Correctly pinged DB")

	// -----------------------
	// Echo server definition
	// -----------------------

	e := echo.New()

	// TODO: Implement logger with config

	// -----------------------
	// Routes definition
	// -----------------------

	// Static endpoint to serve API doc
	// TODO: Write API doc
	e.Static("/docs", "static/docs")

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Pong")
	})

	// -----------------------
	// /Suppliers
	// -----------------------

	// GET a single supplier based on ID
	e.GET("/suppliers/:id", api.GetSupplier(dbConn))

	// GET all records of Suppliers table
	e.GET("/suppliers", api.GetSuppliers(dbConn))

	// POST a new record to suppliers table
	e.POST("/suppliers", api.PostSupplier(dbConn))

	// PUT a new record (update based on Id)
	e.PUT("/suppliers/:id", api.PutSupplier(dbConn))

	// Delete record
	e.DELETE("/suppliers/:id", api.DeleteSuppliers(dbConn))

	// -----------------------
	// Server Start
	// -----------------------

	e.Logger.Fatal(e.Start(":" + port))
}

// Default error check with fatal if err != nil
func checkErrorAndPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
