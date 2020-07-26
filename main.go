package main

import (
	"log"

	"github.com/sqlitebrowser/go-dbhub"
)

func main() {
	// Create a new DBHub.io API object
	db, err := dbhub.New("YOUR_API_KEY_HERE")
	if err != nil {
		log.Fatal(err)
	}
	//db.ChangeServer("https://jctesting1.dbhub.io:8443") // Local testing address

	// Query the remote server
	r, err := db.Query("justinclift", "Join Testing.sqlite", "SELECT sqlite_version()")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Results: %v\n", r)
}
