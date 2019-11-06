package api

import (
	"database/sql"
	"food-inbound/db"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func GetSupplier(dbConn *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		supplier := db.Supplier{}
		db := db.Service{Db: dbConn}
		err := db.GetSupplier(&supplier, id)
		checkErrorAndPanic(err)
		return c.JSON(http.StatusOK, supplier)
	}

}

func GetSuppliers(dbConn *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var suppliers []db.Supplier
		db := db.Service{Db: dbConn}
		err := db.GetSuppliers(&suppliers)
		checkErrorAndPanic(err)
		return c.JSON(http.StatusOK, suppliers)
	}
}

// Default error check with fatal if err != nil
func checkErrorAndPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
