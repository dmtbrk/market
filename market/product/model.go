package product

type Product struct {
	ID     string `json:"id" bson:"_id"`
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	Seller string `json:"seller"`
}

type ListRequest struct {
	Offset int64
	Limit  int64
}

type CreateRequest struct {
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	Seller string `json:"seller"`
}

type UpdateRequest struct {
	ID    string  `bson:"-"`               // Required to find the product.
	Name  *string `bson:"name,omitempty"`  // Optional.
	Price *int64  `bson:"price,omitempty"` // Optional.
}
