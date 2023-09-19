package controller

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"service/internal/repository"
	"strconv"
	"sync"
	"text/template"
	"time"
)

type Controller struct {
	db              *sql.DB
	repo            *repository.Repository
	templates       map[string]*template.Template
	cache           []ProductResponse
	cacheProduct    map[int64]ProductResponse
	cacheProductMut *sync.RWMutex
}

func invalidateCache(cache *[]ProductResponse, cacheProductMut *sync.RWMutex, cacheProduct *map[int64]ProductResponse) {
	timer := time.NewTicker(time.Second * 15)

	go func() {
		for {
			select {
			case <-timer.C:
				log.Println("invalidate cache")
				*cache = []ProductResponse{}

				cacheProductMut.Lock()
				*cacheProduct = make(map[int64]ProductResponse)
				cacheProductMut.Unlock()
			}
		}
	}()
}

func NewController(db *sql.DB, templates map[string]*template.Template) *Controller {
	ctrl := &Controller{
		db:              db,
		templates:       templates,
		repo:            repository.NewRepository(),
		cacheProduct:    make(map[int64]ProductResponse),
		cacheProductMut: &sync.RWMutex{},
	}

	invalidateCache(&ctrl.cache, ctrl.cacheProductMut, &ctrl.cacheProduct)

	return ctrl
}

type ProductResponse struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Price        int64  `json:"price"`
	Manufacturer string `json:"manufacturer"`
}

func (ctrl *Controller) ProductPage(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	productIDVal, err := strconv.ParseInt(productID, 10, 64)

	if err != nil || productIDVal == 0 {
		errorMessage := `
            <html>
                <head>
                    <script type="text/javascript">
                        alert('Пожалуйста, введите в адресе id конкретного товара, либо выберите товар из общего списка, нажав на id.');
                        window.location.href="/";
                    </script>
                </head>
                <body></body>
            </html>`

		w.Write([]byte(errorMessage))
		return
	}

	product, err := ctrl.repo.GetByID(ctrl.db, productIDVal)
	if err != nil {
		log.Println("get product error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tpl := ctrl.templates["product"]
	tpl.Execute(w, product)
}

func (ctrl *Controller) ProductsPage(w http.ResponseWriter, r *http.Request) {
	tpl := ctrl.templates["products"]

	var response []ProductResponse
	if len(ctrl.cache) > 0 {
		log.Println("use cache")
		response = ctrl.cache
	} else {
		log.Println("use db")
		result, err := ctrl.repo.GetProducts(ctrl.db)
		if err != nil {
			log.Println("get products error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, m := range result {
			response = append(response, ProductResponse{
				ID:           m.ID,
				Name:         m.Name,
				Price:        m.Price,
				Manufacturer: m.Manufacturer,
			})
		}

		ctrl.cache = response
	}

	tpl.Execute(w, struct {
		Products []ProductResponse
	}{
		response,
	})
}

func (ctrl *Controller) GetProducts(w http.ResponseWriter, r *http.Request) {
	result, err := ctrl.repo.GetProducts(ctrl.db)
	if err != nil {
		log.Println("get products error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response []ProductResponse

	for _, m := range result {
		response = append(response, ProductResponse{
			Name:         m.Name,
			Price:        m.Price,
			Manufacturer: m.Manufacturer,
		})
	}

	json.NewEncoder(w).Encode(response)
}
