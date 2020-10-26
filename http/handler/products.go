package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ortymid/market/market/product"
	"net/http"
	"strconv"
)

type Products struct {
	ProductService product.Interface
}

func (h *Products) Setup(r *mux.Router) {
	// List
	r.HandleFunc("/products", h.List).Methods(http.MethodGet)
	r.HandleFunc("/products/", h.List).Methods(http.MethodGet)
	// Detail
	r.HandleFunc("/products/{id}", h.Detail).Methods(http.MethodGet)
	r.HandleFunc("/products/{id}/", h.Detail).Methods(http.MethodGet)
	// Create
	r.HandleFunc("/products", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/products/", h.Create).Methods(http.MethodPost)
	// Update
	r.HandleFunc("/products/{id}", h.Update).Methods(http.MethodPatch)
	r.HandleFunc("/products/{id}/", h.Update).Methods(http.MethodPatch)
	// Delete
	r.HandleFunc("/products/{id}", h.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/products/{id}/", h.Delete).Methods(http.MethodDelete)

}

func (h *Products) List(w http.ResponseWriter, r *http.Request) {
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

	lr := product.FindRequest{
		Offset: offset,
		Limit:  limit,
	}

	p, err := h.ProductService.Find(r.Context(), lr)
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

func (h *Products) Detail(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	p, err := h.ProductService.FindOne(r.Context(), id)
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

func (h *Products) Create(w http.ResponseWriter, r *http.Request) {
	var cr product.CreateRequest

	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := h.ProductService.Create(r.Context(), cr)
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

func (h *Products) Update(w http.ResponseWriter, r *http.Request) {
	var ur product.UpdateRequest

	err := json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := mux.Vars(r)["id"]
	ur.ID = id

	p, err := h.ProductService.Update(r.Context(), ur)
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

func (h *Products) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	p, err := h.ProductService.Delete(r.Context(), id)
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
