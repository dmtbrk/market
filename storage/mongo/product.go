package mongo

import (
	"context"
	"github.com/ortymid/market/market/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductStorage struct {
	col *mongo.Collection
}

func NewProductStorage(col *mongo.Collection) *ProductStorage {
	return &ProductStorage{col: col}
}

func (s *ProductStorage) Find(ctx context.Context, r product.FindRequest) ([]*product.Product, error) {
	opts := options.Find().SetSkip(r.Offset).SetLimit(r.Limit)
	cur, err := s.col.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var ps []*product.Product
	for cur.Next(ctx) {
		p := &product.Product{}

		if err := cur.Decode(&p); err != nil {
			return nil, err
		}

		ps = append(ps, p)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return ps, nil
}

func (s *ProductStorage) FindOne(ctx context.Context, id string) (*product.Product, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	p := &product.Product{}

	err = s.col.FindOne(ctx, bson.D{{"_id", oid}}).Decode(p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, product.ErrNotFound
		}
		return nil, err
	}

	return p, nil
}

func (s *ProductStorage) Create(ctx context.Context, r product.CreateRequest) (*product.Product, error) {
	res, err := s.col.InsertOne(ctx, r)
	if err != nil {
		return nil, err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		// TODO: Something to do with id.
	}
	p := &product.Product{
		ID:     oid.Hex(),
		Name:   r.Name,
		Price:  r.Price,
		Seller: r.Seller,
	}
	return p, nil
}

func (s *ProductStorage) Update(ctx context.Context, r product.UpdateRequest) (*product.Product, error) {
	oid, err := primitive.ObjectIDFromHex(r.ID)
	if err != nil {
		return nil, err
	}

	p := &product.Product{}
	f := bson.D{{"_id", oid}}
	u := bson.D{{"$set", r}}
	o := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err = s.col.FindOneAndUpdate(ctx, f, u, o).Decode(p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, product.ErrNotFound
		}
		return nil, err
	}

	return p, nil
}

func (s *ProductStorage) Delete(ctx context.Context, id string) (*product.Product, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	p := &product.Product{}

	err = s.col.FindOneAndDelete(ctx, bson.D{{"_id", oid}}).Decode(p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, product.ErrNotFound
		}
		return nil, err
	}

	return p, nil
}
