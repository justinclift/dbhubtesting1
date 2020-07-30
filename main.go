package main

import (
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
	r, err := db.Query("justinclift", "Join Testing.sqlite",
		`SELECT table1.Name, table2.value
			FROM table1 JOIN table2
			USING (id)
			ORDER BY table1.id`)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Results: %v\n", r)
}
