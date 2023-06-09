package database

import (
	"database/sql"
	"fmt"
	"log"

	models "github.com/jain-chetan/companyservice/model"

	"github.com/google/uuid"
)

func CreateConnection() *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=testdb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func CreateCompanyQuery(company models.Company) uuid.UUID {

	// create the postgres db connection
	db := CreateConnection()

	// close the db connection
	defer db.Close()

	var id uuid.UUID
	createTable := `CREATE TABLE IF NOT EXISTS company( ID uuid DEFAULT uuid_generate_v1(), NAME TEXT NOT NULL, DESCRIPTION TEXT NOT NULL, EMPLOYEES INT, REGISTERED BOOL, TYPE TEXT)`
	_ = db.QueryRow(createTable)
	sqlStatement := `INSERT INTO company (name, description,employees,registered,type) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	// execute the sql statement
	err := db.QueryRow(sqlStatement, company.Name, company.Description, company.Employees, company.Registered, company.Type).Scan(&id)

	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted %v", id)

	// return the id
	return id
}

// get one company from the DB by its id
func GetCompanyQuery(id uuid.UUID) (models.Company, error) {
	// create the postgres db connection
	db := CreateConnection()

	// close the db connection
	defer db.Close()

	// create a company of models.company type
	var company models.Company

	// create the select sql query
	sqlStatement := `SELECT * FROM company WHERE id=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to company
	err := row.Scan(&company.ID, &company.Name, &company.Description, &company.Employees, &company.Registered, &company.Type)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return company, err
	case nil:
		return company, nil
	default:
		log.Printf("Unable to scan the row. %v", err)
	}

	// return empty company on error
	return company, err
}

// update company in the DB
func PatchCompanyQuery(id uuid.UUID, company models.Company) uuid.UUID {

	// create the postgres db connection
	db := CreateConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE company SET name=$2, description=$3, employees=$4, registered=$5, type=$6 WHERE id=$1`

	// execute the sql statement
	_, err := db.Exec(sqlStatement, id, company.Name, company.Description, company.Employees, company.Registered, company.Type)

	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
	}

	return id
}

// delete company in the DB
func DeleteCompanyQuery(id uuid.UUID) uuid.UUID {

	// create the postgres db connection
	db := CreateConnection()

	// close the db connection
	defer db.Close()

	sqlStatement := `DELETE FROM company WHERE id=$1`

	// execute the sql statement
	_, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
	}

	return id
}
