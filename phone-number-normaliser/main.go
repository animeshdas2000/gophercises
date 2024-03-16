package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	dbname = "pnn"
)

func main() {
	// db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to DB: %v\n", err)
	// 	os.Exit(1)
	// }
	// err = resetDB(db, dbname)
	// if err != nil {
	// 	panic(err)
	// }
	// db.Close()
	// phoneNumbers := []string{
	// 	"1234567890",
	// 	"123 456 7891",
	// 	"(123) 456 789",
	// 	"(123) 456-7893",
	// 	"123-456-7894",
	// 	"123-456-7890",
	// 	"(123)456-7892",
	// }
	psqlInfo := fmt.Sprintf("%s%s", os.Getenv("DATABASE_URL"), dbname)
	db, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	must(createPhoneNumberTable(db))
	// for _, phoneNumber := range phoneNumbers {
	// 	_, err = insertPhone(db, phoneNumber)
	// 	must(err)
	// }
	id, err := insertPhone(db, "0000023223")
	must(err)
	ph, err := getPhone(db, id)
	must(err)
	fmt.Println(ph)
	phoneNumbers, err := allPhones(db)
	must(err)
	for _, ph := range phoneNumbers {
		fmt.Printf("%+v\n", ph)
	}
}

func getPhone(db *sql.DB, id int) (string, error) {
	var phoneNumber string
	statement := `SELECT * FROM phone_numbers where id=$1`
	err := db.QueryRow(statement, id).Scan(&id, &phoneNumber)
	if err != nil {
		return "", err
	}
	return phoneNumber, nil
}

type phone struct {
	id     int
	number string
}

func allPhones(db *sql.DB) ([]phone, error) {
	rows, err := db.Query("SELECT id,value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []phone
	for rows.Next() {
		var p phone
		if err := rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	return ret, nil

}
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func createPhoneNumberTable(db *sql.DB) error {
	statement := `
	CREATE TABLE IF NOT EXISTS phone_numbers (
		id SERIAL,
		value VARCHAR(255)
	)`
	_, err := db.Exec(statement)
	// fmt.Printf("%d", id)
	if err != nil {
		return err
	}
	return nil
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `
	INSERT INTO phone_numbers(value) VALUES($1) RETURNING id
	`
	var id int
	// args :=
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, err
}

// func resetDB(db *sql.DB, name string) error {

// 	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
// 	if err != nil {
// 		return err
// 	}
// 	return createDB(db, name)
// }

// func createDB(db *sql.DB, name string) error {
// 	_, err := db.Exec("CREATE DATABASE " + name)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func normalise(ph string) string {
	// Build a new string on the go
	var buf bytes.Buffer
	//loop over the string value
	for _, ch := range ph {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

// func normalise(ph string) string {
// 	re := regexp.MustCompile("[^0-9]")
// 	//re := regexp.MustCompile("\\D")
// 	return re.ReplaceAllLiteralString(ph, "")

// }
