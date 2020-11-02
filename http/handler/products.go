package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ortymid/market/market/product"
	"net/http"
	"net/url"
	"strconv"
)

type Products struct {
	ProductService product.Interface
}

func (h *Products) Setup(r *mux.Router) {
	// Find
	r.HandleFunc("/products", h.Find).Methods(http.MethodGet)
	r.HandleFunc("/products/", h.Find).Methods(http.MethodGet)
	// FindOne
	r.HandleFunc("/products/{id}", h.FindOne).Methods(http.MethodGet)
	r.HandleFunc("/products/{id}/", h.FindOne).Methods(http.MethodGet)
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

func (h *Products) Find(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	findReq, err := makeFindRequestFromQuery(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := h.ProductService.Find(r.Context(), findReq)
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

func makeFindRequestFromQuery(query url.Values) (r product.FindRequest, err error) {
	offset, err := strconv.ParseInt(query.Get("offset"), 10, 64)
	if err != nil {
		return r, errors.New("valid offset query parameter required")
	}

	limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
	if err != nil {
		return r, errors.New("valid limit query parameter required")
	}

	var name *string
	if names, ok := query["name"]; ok && len(names) > 0 {
		name = &names[0]
	}

	var priceFrom *int64
	if pfs, ok := query["price_from"]; ok && len(pfs) > 0 {
		p, err := strconv.ParseInt(pfs[0], 10, 64)
		if err != nil {
			return r, fmt.Errorf("invalid price_from: %w", err)
		}
		priceFrom = &p
	}

	var priceTo *int64
	if pfs, ok := query["price_to"]; ok && len(pfs) > 0 {
		p, err := strconv.ParseInt(pfs[0], 10, 64)
		if err != nil {
			return r, fmt.Errorf("invalid price_to: %w", err)
		}
		priceTo = &p
	}

	var priceRange *product.PriceRange
	if priceFrom != nil || priceTo != nil {
		priceRange = &product.PriceRange{
			From: priceFrom,
			To:   priceTo,
		}
	}

	return product.FindRequest{
		Offset:     offset,
		Limit:      limit,
		Name:       name,
		PriceRange: priceRange,
	}, nil
}

func (h *Products) FindOne(w http.ResponseWriter, r *http.Request) {
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
