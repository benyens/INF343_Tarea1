package db

import (
	"database/sql"
	"os"
	"log"
)

func MakeMigrate(db *sql.DB, schemaFile string) {
	schema, err := os.ReadFile(schemaFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Base de datos y tablas creadas exitosamente")
}