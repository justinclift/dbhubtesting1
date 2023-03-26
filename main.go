package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/sqlitebrowser/go-dbhub"
)

func main() {
	// Read in our configuration
	err := ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new DBHub.io API object
	db, err := dbhub.New(Conf.Api.APIKey) // Use the API key stored in our local config file
	if err != nil {
		log.Fatal(err)
	}
	if Conf.Api.Server != "" {
		db.ChangeServer(Conf.Api.Server) // If a server was given in our local config, use that instead of the default
	}
	db.ChangeVerifyServerCert(Conf.Api.VerifyCert)

	// Retrieve the list of standard databases in your account
	databases, err := db.Databases()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Standard databases:")
	for _, j := range databases {
		fmt.Printf("  * %s\n", j)
	}
	fmt.Println()

	// Retrieve the list of LIVE databases in your account
	liveDatabases, err := db.DatabasesLive()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Live databases:")
	for _, j := range liveDatabases {
		fmt.Printf("  * %s\n", j)
	}
	fmt.Println()

	// Read a test database into memory, ready for uploading as we need
	z, err := os.Open(filepath.Join("test_data", "Join Testing with index.sqlite"))
	if err != nil {
		log.Fatal(err)
	}
	defer z.Close()
	testDB, err := io.ReadAll(z)
	if err != nil {
		log.Fatal(err)
	}

	// Upload a standard database for doing stuff with
	dbName := "Upload STANDARD test.sqlite"
	for _, j := range databases {
		if j == dbName {
			// If the database already exists in the account, then remove it first
			err = db.Delete(dbName)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	err = db.Upload(dbName, dbhub.UploadInformation{}, &testDB)
	if err != nil {
		log.Fatal(err)
	}

	// Upload the test database as a LIVE database
	liveName := "Upload LIVE test.sqlite"
	for _, j := range liveDatabases {
		if j == liveName {
			// If the database already exists in the account, then remove it first
			err = db.Delete(liveName)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	err = db.UploadLive(liveName, &testDB)
	if err != nil {
		log.Fatal(err)
	}

	// Query the remote server
	// FIXME: Add a (reasonable) test for blob values
	r, err := db.Query("default", dbName, dbhub.Identifier{}, false,
		`SELECT table1.Name, table2.value
			FROM table1 JOIN table2
			USING (id)
			ORDER BY table1.id`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Query results: %v\n", r) // TODO: Format this a bit better
	fmt.Println()

	// Retrieve the list of tables in a remote standard database
	tables, err := db.Tables("default", dbName, dbhub.Identifier{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tables:")
	for _, j := range tables {
		fmt.Printf("  * %s\n", j)
	}
	fmt.Println()

	// Retrieve the list of views in the remote standard database
	views, err := db.Views("default", dbName, dbhub.Identifier{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Views:")
	for _, j := range views {
		fmt.Printf("  * %s\n", j)
	}
	fmt.Println()

	// Retrieve the list of indexes in the remote standard database
	indexes, err := db.Indexes("default", dbName, dbhub.Identifier{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Indexes:")
	for _, j := range indexes {
		fmt.Printf("  * '%s' on table '%s'\n", j.Name, j.Table)
	}
	fmt.Println()

	// Retrieve the column info for a table or view in the remote standard database
	table := "table1"
	columns, err := db.Columns("default", dbName, dbhub.Identifier{}, table)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Columns on table or view '%s':\n", table)
	for _, j := range columns {
		fmt.Printf("  * '%v':\n", j.Name)
		fmt.Printf("      Cid: %v\n", j.Cid)
		fmt.Printf("      Data Type: %v\n", j.DataType)
		fmt.Printf("      Default Value: %v\n", j.DfltValue)
		fmt.Printf("      Not Null: %v\n", j.NotNull)
		fmt.Printf("      Primary Key: %v\n", j.Pk)
	}
	fmt.Println()

	// Retrieve the remote standard database file
	dbStream, err := db.Download("default", dbName, dbhub.Identifier{})
	if err != nil {
		log.Fatal(err)
	}

	// Save the standard database file in the current directory
	buf, err := io.ReadAll(dbStream)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(dbName, buf, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Saved database file as '%s'\n", dbName)
	fmt.Println()
	if y, err := os.Stat(dbName); err != nil || y.Size() != 16384 {
		log.Fatal(err)
	}
	if err = os.Remove(dbName); err != nil {
		log.Fatal(err)
	}
	log.Println("Deleted downloaded file")

	// Run Execute() on the remote live database
	var rows int
	rows, err = db.Execute("default", liveName, "UPDATE table1 SET Name = 'Testing 1' WHERE id = 1")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Number of rows changed by the Execute statement: %d", rows)

	// Remove the uploaded files
	err = db.Delete(dbName)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Delete(liveName)
	if err != nil {
		log.Fatal(err)
	}
}
