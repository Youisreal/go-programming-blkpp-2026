package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"gin-auth/internal/database"

	"github.com/joho/godotenv"
)

func run(confirm bool) error {
	if !confirm {
		return errors.New("penghapusan dibatalkan; gunakan -confirm untuk menghapus tabel " + tableName + " beserta seluruh datanya")
	}
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}
	db, err := database.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	conn, err := db.DB()
	if err != nil {
		return err
	}
	defer conn.Close()

	// Driver PostgreSQL GORM memakai CASCADE: dependency seperti foreign key ikut dihapus.
	if err := db.Migrator().DropTable("categories"); err != nil {
		return fmt.Errorf("hapus tabel categories gagal: %w", err)
	}
	log.Printf("Tabel %s sudah tidak ada", tableName)
	return nil
}

func main() {
	confirm := flag.Bool("confirm", false, "hapus tabel categories beserta seluruh data dan dependency (CASCADE)")
	flag.Parse()
	if err := run(*confirm); err != nil {
		log.Fatal(err)
	}
}
