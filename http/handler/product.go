package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ortymid/market/market"
)

// ProductHandler forwards product requests to the business logic.
type ProductHandler struct {
	Market market.Interface
}

func (h *ProductHandler) Setup(r *mux.Router) {
	r.HandleFunc("/", h.List).Methods(http.MethodGet)
	r.HandleFunc("/", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/{id}", h.Detail).Methods(http.MethodGet)
	r.HandleFunc("/{id}", h.Edit).Methods(http.MethodPut)
	r.HandleFunc("/{id}", h.Delete).Methods(http.MethodDelete)
}

// List handles requests for all products.
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.Market.Products()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := productListResponse(products)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Detail handles requests for the specific product detail.
func (h *ProductHandler) Detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString, ok := vars["id"]
	if !ok {
		http.Error(w, "id not specified", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "id is not an integer", http.StatusBadRequest)
		return
	}

	product, err := h.Market.Product(id)
	if err != nil {
		err = fmt.Errorf("getting product: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if product == nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	resp := productDetailResponse(*product)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Create handles requests for creation of new products.
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(KeyUserID).(string)
	if !ok {
		http.Error(w, "authorization required", http.StatusForbidden)
		return
	}

	data := struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		err = fmt.Errorf("decoding create product request: %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product := &market.Product{Name: data.Name, Price: data.Price, Seller: userID}
	product, err = h.Market.AddProduct(product, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, &market.ErrPermission{}) {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}
	if product == nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	resp := productCreateResponse(*product)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Edit handles product edit requests.
func (h *ProductHandler) Edit(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(KeyUserID).(string)
	if !ok {
		http.Error(w, "authorization required", http.StatusForbidden)
		return
	}

	id, err := getVarProductID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product := &market.Product{ID: id}

	data := struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		err = fmt.Errorf("decoding create product request: %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product.Name = data.Name
	product.Price = data.Price
	product.Seller = userID

	product, err = h.Market.ReplaceProduct(product, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, &market.ErrPermission{}) || errors.Is(err, market.ErrProductNotFound) {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}
	if product == nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	resp := productEditResponse(*product)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Delete handles product delete requests.
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(KeyUserID).(string)
	if !ok {
		http.Error(w, "authorization: token required", http.StatusForbidden)
		return
	}

	id, err := getVarProductID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// product, err := h.market.Product(id)
	// if err != nil {
	// 	http.Error(w, http.StatusInternalServerError, err)
	// 	return
	// }
	// if product == nil {
	// 	http.Error(w, http.StatusInternalServerError, errors.New("something went wrong"))
	// 	return
	// }

	err = h.Market.DeleteProduct(id, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, &market.ErrPermission{}) || errors.Is(err, market.ErrProductNotFound) {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func getVarProductID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idString, ok := vars["id"]
	if !ok {
		return 0, errors.New("id not specified")
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, errors.New("id is not an integer")
	}
	return id, nil
}

type productListResponse []*market.Product

func (r productListResponse) MarshalJSON() ([]byte, error) {
	type respProduct struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Price  int    `json:"price"`
		Seller string `json:"seller"`
	}

	respProducts := make([]respProduct, len(r))
	for i, p := range r {
		respProducts[i] = respProduct{
			ID:     p.ID,
			Name:   p.Name,
			Price:  p.Price,
			Seller: p.Seller,
		}
	}

	return json.Marshal(respProducts)
}

type productDetailResponse market.Product

func (r productDetailResponse) MarshalJSON() ([]byte, error) {
	type respProduct struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Price  int    `json:"price"`
		Seller string `json:"seller"`
	}

	return json.Marshal(respProduct(r))
}

type productCreateResponse market.Product

func (r productCreateResponse) MarshalJSON() ([]byte, error) {
	type respProduct struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Price  int    `json:"price"`
		Seller string `json:"seller"`
	}

	return json.Marshal(respProduct(r))
}

type productEditResponse market.Product

func (r productEditResponse) MarshalJSON() ([]byte, error) {
	type respProduct struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Price  int    `json:"price"`
		Seller string `json:"seller"`
	}

	return json.Marshal(respProduct(r))
}
