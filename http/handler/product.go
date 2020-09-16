package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	listHandler := PaginationMiddleware(http.HandlerFunc(h.List))
	r.Handle("/", listHandler).Methods(http.MethodGet)
	r.HandleFunc("/", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/{id}", h.Detail).Methods(http.MethodGet)
	r.HandleFunc("/{id}", h.Edit).Methods(http.MethodPut)
	r.HandleFunc("/{id}", h.Delete).Methods(http.MethodDelete)
}

// List handles requests for all products.
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page, ok := ctx.Value(KeyPage).(Page)
	if !ok {
		log.Printf("ERROR: context value for KeyPage is %v", page)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	log.Printf("%#v", page)

	products, err := h.Market.Products(ctx, page.Offset, page.Limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("%#v", products)

	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Detail handles requests for the specific product detail.
func (h *ProductHandler) Detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	product, err := h.Market.Product(ctx, id)
	if err != nil {
		err = fmt.Errorf("getting product: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if product == nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Create handles requests for creation of new products.
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(KeyUserID).(string)
	if !ok {
		http.Error(w, "authorization required", http.StatusForbidden)
		return
	}

	var pr market.AddProductRequest
	err := json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		err = fmt.Errorf("decoding create product request: %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pr.Seller = userID

	product, err := h.Market.AddProduct(ctx, pr, userID)
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

	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Edit handles product edit requests.
func (h *ProductHandler) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(KeyUserID).(string)
	if !ok {
		http.Error(w, "authorization required", http.StatusForbidden)
		return
	}

	id, err := getVarProductID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var pr market.EditProductRequest
	err = json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		err = fmt.Errorf("decoding create product request: %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pr.ID = id

	product, err := h.Market.EditProduct(ctx, pr, userID)
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

	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Delete handles product delete requests.
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(KeyUserID).(string)
	if !ok {
		http.Error(w, "authorization: token required", http.StatusForbidden)
		return
	}

	id, err := getVarProductID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Market.DeleteProduct(ctx, id, userID)
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
		return 0, errors.New("product id not specified")
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, errors.New("product id is not an integer")
	}
	return id, nil
}
