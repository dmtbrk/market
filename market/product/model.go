package product

type Product struct {
	ID     string `json:"id" bson:"_id"`
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	Seller string `json:"seller"`
}

type FindRequest struct {
	Offset int64
	Limit  int64

	// Optional filters.
	Name       *string
	PriceRange *PriceRange
	Seller     *string
}

type PriceRange struct {
	From *int64 // nil means no lower limit
	To   *int64 // nil means no upper limit
}

type CreateRequest struct {
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	Seller string `json:"seller"`
}

type UpdateRequest struct {
	ID    string  `json:"-" bson:"-"`                             // Required to find the product.
	Name  *string `json:"name,omitempty" bson:"name,omitempty"`   // Optional.
	Price *int64  `json:"price,omitempty" bson:"price,omitempty"` // Optional.
}
