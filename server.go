package main

import (
	"database/sql"
	"fmt"
	"food-inbound/api"
	"food-inbound/gApi"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	// TODO: Implement real service
	// GApi service test
	srv := gApi.GService{}
	err := srv.Service()
	checkErrorAndPanic(err)
	resp, err := srv.ReadRange("Fornitori!A2:B3")
	checkErrorAndPanic(err)
	fmt.Printf("Returned data is: %v\n", resp)

	type result struct {
		nome      string
		indirizzo string
	}
	var fornitori []result
	for _, row := range resp {
		fmt.Printf("Row: %v\n", row)
		f := result{
			nome:      fmt.Sprintf("%v", row[0]),
			indirizzo: fmt.Sprintf("%v", row[1]),
		}
		fornitori = append(fornitori, f)
	}
	fmt.Printf("Fornitori: %v\n", fornitori)
	for _, row := range fornitori {
		fmt.Printf("Struct row: %s - %s\n", row.nome, row.indirizzo)
	}
	// Append Test
	res, err := srv.Append("Fornitori!A2:B4", [][]interface{}{{"Lindt", "Siur Sprungli"}})
	fmt.Printf("Append response: %v\n", res)

	// FIXME: Return used for dev purpose, remove when done
	return

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
