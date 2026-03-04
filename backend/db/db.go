package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", "./aidoctor.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	schema, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatal("Failed to read schema:", err)
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		log.Fatal("Failed to initialize schema:", err)
	}

	log.Println("✅ Database initialized")
}
