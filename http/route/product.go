package route

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ortymid/market/market/product"
	"net/http"
	"strconv"
)

type Product struct {
	ProductService product.Interface
}

func (pr *Product) Setup(r *mux.Router) {
	// List
	r.HandleFunc("/products", pr.List).Methods(http.MethodGet)
	r.HandleFunc("/products/", pr.List).Methods(http.MethodGet)
	// Detail
	r.HandleFunc("/products/{id}", pr.Detail).Methods(http.MethodGet)
	r.HandleFunc("/products/{id}/", pr.Detail).Methods(http.MethodGet)
	// Create
	r.HandleFunc("/products", pr.Create).Methods(http.MethodPost)
	r.HandleFunc("/products/", pr.Create).Methods(http.MethodPost)
	// Update
	r.HandleFunc("/products/{id}", pr.Update).Methods(http.MethodPatch)
	r.HandleFunc("/products/{id}/", pr.Update).Methods(http.MethodPatch)
	// Delete
	r.HandleFunc("/products/{id}", pr.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/products/{id}/", pr.Delete).Methods(http.MethodDelete)

}

func (pr *Product) List(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset, err := strconv.ParseInt(params.Get("offset"), 10, 64)
	if err != nil {
		http.Error(w, "valid offset query parameter required", http.StatusBadRequest)
		return
	}
	limit, err := strconv.ParseInt(params.Get("limit"), 10, 64)
	if err != nil {
		http.Error(w, "valid limit query parameter required", http.StatusBadRequest)
		return
	}

	lr := product.ListRequest{
		Offset: offset,
		Limit:  limit,
	}

	p, err := pr.ProductService.List(r.Context(), lr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (pr *Product) Detail(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	p, err := pr.ProductService.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (pr *Product) Create(w http.ResponseWriter, r *http.Request) {
	var cr product.CreateRequest

	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := pr.ProductService.Create(r.Context(), cr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (pr *Product) Update(w http.ResponseWriter, r *http.Request) {
	var ur product.UpdateRequest

	err := json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := mux.Vars(r)["id"]
	ur.ID = id

	p, err := pr.ProductService.Update(r.Context(), ur)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (pr *Product) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	p, err := pr.ProductService.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
