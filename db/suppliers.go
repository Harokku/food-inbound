package db

import (
	"database/sql"
	"errors"
)

type Supplier struct {
	Id             string
	Name           string
	Address        string
	ReferenceName  string
	ReferenceMail  string
	ReferencePhone string
}

// Query Suppliers table and return a single row matching ID, populate ONLY Id, Name, Address
func (s Service) GetSupplier(supplier *Supplier, id string) error {
	sqlStatement := `SELECT name, address FROM suppliers where id=$1`
	row := s.Db.QueryRow(sqlStatement, id)
	switch err := row.Scan(&supplier.Name, &supplier.Address); err {
	case sql.ErrNoRows:
		return errors.New("no rows were retrieved")
	case nil:
		supplier.Id = id
		return nil
	default:
		return errors.New("error retrieving requested row")
	}
}

// TODO: Add pagination option
// Query Suppliers table and return all records
func (s Service) GetSuppliers(suppliers *[]Supplier) error {
	sqlStatement := `SELECT id, name, address, reference_name, reference_mail, reference_phone FROM suppliers`
	rows, err := s.Db.Query(sqlStatement)
	if err != nil {
		return err
	}

	defer rows.Close()

	var supplier Supplier
	for rows.Next() {
		err = rows.Scan(&supplier.Id, &supplier.Name, &supplier.Address, &supplier.ReferenceName, &supplier.ReferenceMail, &supplier.ReferencePhone)
		// TODO: Implement better error handling
		if err != nil {
			panic(err)
		}
		*suppliers = append(*suppliers, supplier)
	}

	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
