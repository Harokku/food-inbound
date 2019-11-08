package db

import (
	"database/sql"
	"errors"
)

type Supplier struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	ReferenceName  string `json:"reference_name"`
	ReferenceMail  string `json:"reference_mail"`
	ReferencePhone string `json:"reference_phone"`
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
			return err
		}
		*suppliers = append(*suppliers, supplier)
	}

	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

// Post a new supplier,
// Return a string indicating the new inserted UUID or an error
// TODO: Implement interface to make return database id type agnostic
func (s Service) PostSupplier(supplier Supplier) (string, error) {
	sqlStatement := `
		INSERT INTO suppliers (name, address, reference_name, reference_mail, reference_phone) 
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id
	`
	var res string
	err := s.Db.QueryRow(sqlStatement, supplier.Name, supplier.Address, supplier.ReferenceName, supplier.ReferenceMail, supplier.ReferencePhone).Scan(&res)
	if err != nil {
		return "", err
	}
	return res, nil
}

// Update selected record
// Pass the updated Supplier with Id set
// Beware to fill every field needed, omitted wil be nulled
func (s Service) PutSuppliers(supplier Supplier) error {
	sqlStatement := `
		UPDATE suppliers
		SET name=$2, address=$3, reference_name=$4, reference_mail=$5, reference_phone=$6
		WHERE id=$1
`
	_, err := s.Db.Exec(sqlStatement, supplier.Id, supplier.Name, supplier.Address, supplier.ReferenceName, supplier.ReferenceMail, supplier.ReferencePhone)
	if err != nil {
		return err
	}
	return nil
}

// Delete selected record
func (s Service) DeleteSupplier(id string) error {
	sqlStatement := `
		DELETE FROM suppliers
		WHERE id=$1
`
	_, err := s.Db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	return nil
}
