package api

import (
	"database/sql"
	dbRef "food-inbound/db"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func GetSupplier(dbConn *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		supplier := dbRef.Supplier{}
		db := dbRef.Service{Db: dbConn}
		err := db.GetSupplier(&supplier, id)
		checkErrorAndPanic(err)
		return c.JSON(http.StatusOK, supplier)
	}

}

func GetSuppliers(dbConn *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var suppliers []dbRef.Supplier
		db := dbRef.Service{Db: dbConn}
		err := db.GetSuppliers(&suppliers)
		checkErrorAndPanic(err)
		return c.JSON(http.StatusOK, suppliers)
	}
}

func PostSupplier(dbConn *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		s := new(dbRef.Supplier)
		if err := c.Bind(s); err != nil {
			return c.String(http.StatusBadRequest, "Error binding POST body")
		}

		db := dbRef.Service{Db: dbConn}
		var err error
		s.Id, err = db.PostSupplier(*s)
		checkErrorAndPanic(err)

		return c.JSON(http.StatusCreated, s)
	}
}

// Default error check with fatal if err != nil
func checkErrorAndPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
