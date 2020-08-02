package main

import (
	"fmt"
	"log"

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

	// Query the remote server
	// FIXME: Add a (reasonable) test for blob values
	r, err := db.Query("justinclift", "Join Testing.sqlite", false,
		`SELECT table1.Name, table2.value
			FROM table1 JOIN table2
			USING (id)
			ORDER BY table1.id`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Query results: %v\n", r) // TODO: Format this a bit better
	fmt.Println()

	// Retrieve the list of tables in the remote database
	tables, err := db.Tables("justinclift", "Join Testing.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tables:")
	for _, j := range tables {
		fmt.Printf("  * %s\n", j)
	}
	fmt.Println()

	// Retrieve the list of views in the remote database
	views, err := db.Views("justinclift", "Join Testing.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Views:")
	for _, j := range views {
		fmt.Printf("  * %s\n", j)
	}
	fmt.Println()

	// Retrieve the list of indexes in the remote database
	indexes, err := db.Indexes("justinclift", "Join Testing.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Indexes:")
	for i, j := range indexes {
		fmt.Printf("  * '%s' on table '%s'\n", i, j)
	}
	fmt.Println()

	// Retrieve the column info for a table or view in the remote database
	table := "table1"
	columns, err := db.Columns("justinclift", "Join Testing.sqlite", table)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Columns on table or view '%s':\n", table)
	for _, j := range columns {
		fmt.Printf("  * '%v':\n", j.Name)
		fmt.Printf("      Autoincrement: %v\n", j.Autoinc)
		fmt.Printf("      Cid: %v\n", j.Cid)
		fmt.Printf("      Collation Sequence: %v\n", j.CollSeq)
		fmt.Printf("      Data Type: %v\n", j.DataType)
		fmt.Printf("      Default Value: %v\n", j.DfltValue)
		fmt.Printf("      Not Null: %v\n", j.NotNull)
		fmt.Printf("      Primary Key: %v\n", j.Pk)
	}
	fmt.Println()

	// Generate and display diff for two commits of the remote database
	diffs, err := db.Diff("justinclift", "Join Testing.sqlite", "c82ba65add364427e9af3f540be8bf98e8cd6bdb825b07c334858e816c983db0", "", "", "adf78104254ece17ff40dab80ae800574fa5d429a4869792a64dcf2027cd9cd9", dbhub.PreservePkMerge)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("SQL statements for turning the first commit into the second:\n")
	for _, i := range diffs.Diff {
		if i.Schema != nil {
			fmt.Printf("%s\n", i.Schema.Sql)
		}
		for _, j := range i.Data {
			fmt.Printf("%s\n", j.Sql)
		}
	}
	fmt.Println()
}
