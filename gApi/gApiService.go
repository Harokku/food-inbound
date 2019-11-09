package gApi

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"
)

type GService struct {
	srv     *sheets.Service
	sheetId string
}

// Represent Google API service with auth and sheet ID read from env
// GOOGLE_API is the auth secret
// SHEETS_ID is the sheet to read from
func (s *GService) Service() error {
	secret := os.Getenv("GOOGLE_API")
	if secret == "" {
		panic("Can't read secret from env.")
	}
	conf, err := google.JWTConfigFromJSON([]byte(secret), sheets.SpreadsheetsScope)
	checkErrorAndPanic(err)

	srv, err := sheets.NewService(context.TODO(), option.WithHTTPClient(conf.Client(context.TODO())))
	checkErrorAndPanic(err)

	s.srv = srv
	s.sheetId = os.Getenv("SHEET_ID")
	return nil
}

// Read selected range and return the result
// r string: the range to retrieve in the !A1 format (Sheet!A1:B2)
// Return 2D array with results
func (s GService) ReadRange(r string) ([][]interface{}, error) {
	res, err := s.srv.Spreadsheets.Values.Get(s.sheetId, r).Do()
	if err != nil {
		return nil, err
	}

	return res.Values, nil
}

// Append data after selected range and return the result
// r string: the range after witch append data in the !A1 format (Sheet!A1:B2)
// data [][]interface{}: 2D array with data to append
// Return int with return code (HTTPStatusCode)
func (s GService) Append(r string, data [][]interface{}) (int, error) {
	var values = sheets.ValueRange{
		Values: data,
	}
	res, err := s.srv.Spreadsheets.Values.Append(s.sheetId, r, &values).ValueInputOption("RAW").Do()
	if err != nil {
		return 0, err
	}
	fmt.Printf("GApi responded with: %d\n", res.HTTPStatusCode)
	return res.HTTPStatusCode, nil
}

// Default error check with fatal if err != nil
func checkErrorAndPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
