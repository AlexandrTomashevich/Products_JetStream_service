package test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"service/internal/models"
	"service/internal/repository"
)

var testDB *sql.DB
var repo *repository.Repository

func TestMain(m *testing.M) {
	// Инициализация тестовой базы данных
	connStr := "user=postgres dbname=postgreSQL sslmode=disable password=" + os.Getenv("DB_PASSWORD")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	testDB = db
	repo = repository.NewRepository()

	code := m.Run()

	testDB.Close()

	os.Exit(code)
}

func TestAddProduct(t *testing.T) {
	product := models.Product{
		Name:         "TestProduct",
		Price:        1000,
		Manufacturer: "TestManufacturer",
	}

	err := repo.AddProduct(testDB, product)
	if err != nil {
		t.Fatalf("Failed to add product: %v", err)
	}

}
