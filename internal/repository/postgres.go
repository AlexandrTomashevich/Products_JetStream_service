package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"service/internal/models"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetByID(db *sql.DB, id int64) (*models.Product, error) {
	row := db.QueryRow("select id, name, price, manufacturer from products where id = $1", id)

	var m models.Product
	if err := row.Scan(&m.ID, &m.Name, &m.Price, &m.Manufacturer); err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *Repository) GetProducts(db *sql.DB) ([]models.Product, error) {
	rows, err := db.Query("select id, name, price, manufacturer from products order by id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Product
	for rows.Next() {
		var m models.Product
		if err := rows.Scan(&m.ID, &m.Name, &m.Price, &m.Manufacturer); err != nil {
			return nil, err
		}

		result = append(result, m)
	}

	return result, nil
}

func (r *Repository) AddProduct(db *sql.DB, p models.Product) error {
	_, err := db.Exec("insert into products (name, price, manufacturer) values ($1, $2, $3)", p.Name, p.Price, p.Manufacturer)
	return err
}
