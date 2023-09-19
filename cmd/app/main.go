package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"text/template"

	"service/internal/controller"
)

func migrateDB(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id serial primary key,
			name text not null,
			price integer not null,
			manufacturer text not null
		)
	`)
	if err != nil {
		panic(err)
	}
}

var templates = map[string]*template.Template{
	"products": nil,
	"product":  nil,
}

func main() {
	dbConnectionString := "user=postgres dbname=postgreSQL sslmode=disable password=" + os.Getenv("DB_PASSWORD")
	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		panic(err)
	}

	migrateDB(db)

	for tplName := range templates {
		tpl, err := template.ParseGlob("./static/" + tplName + ".html")
		if err != nil {
			panic(err)
		}
		templates[tplName] = tpl
	}

	ctrl := controller.NewController(db, templates)

	http.HandleFunc("/", ctrl.ProductsPage)
	http.HandleFunc("/product", ctrl.ProductPage)
	http.HandleFunc("/products", ctrl.GetProducts)

	log.Println("start service")
	http.ListenAndServe(":8000", nil)
}
