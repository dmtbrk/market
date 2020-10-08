package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ortymid/market/market/product"
	"strconv"
)

type ProductStorage struct {
	rdb *redis.Client

	baseKey string
	idsKey  string
}

func NewProductStorage(rdb *redis.Client, key string) *ProductStorage {
	idsKey := fmt.Sprintf("%s:ids", key)

	return &ProductStorage{rdb: rdb, baseKey: key, idsKey: idsKey}
}

func (s *ProductStorage) List(ctx context.Context, r product.ListRequest) ([]*product.Product, error) {
	start, stop := r.Offset, r.Offset+r.Limit-1

	ids, err := s.rdb.ZRange(ctx, s.idsKey, start, stop).Result()
	if err != nil {
		return nil, err
	}

	var products []*product.Product
	for _, id := range ids {
		p, err := s.getProductFromHash(ctx, id)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func (s *ProductStorage) Get(ctx context.Context, id string) (*product.Product, error) {
	return s.getProductFromHash(ctx, id)
}

func (s *ProductStorage) Create(ctx context.Context, r product.CreateRequest) (*product.Product, error) {
	// Get new id.
	id, err := s.rdb.Incr(ctx, fmt.Sprintf("%s:id", s.baseKey)).Result()
	if err != nil {
		return nil, fmt.Errorf("getting new id: %w", err)
	}

	// Prepare new product.
	p := &product.Product{
		ID:     strconv.FormatInt(id, 10),
		Name:   r.Name,
		Price:  r.Price,
		Seller: r.Seller,
	}

	// Store new product.
	err = s.setProductToHash(ctx, p)
	if err != nil {
		return nil, err
	}

	// Update ids sorted set.
	zID := &redis.Z{
		Score:  float64(id),
		Member: id,
	}
	err = s.rdb.ZAdd(ctx, fmt.Sprintf("%s:ids", s.baseKey), zID).Err()
	if err != nil {
		// Rollback.
		err = s.rdb.Del(ctx, s.hashKey(p.ID)).Err()
		if err != nil {
			return nil, fmt.Errorf("rolling back ids sorted set update: %w", err)
		}
		return nil, fmt.Errorf("updating ids sorted set: %w", err)
	}

	return p, nil
}

func (s *ProductStorage) Update(ctx context.Context, r product.UpdateRequest) (*product.Product, error) {
	// Get product checking for existence.
	p, err := s.getProductFromHash(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	// Update not-nil fields.
	if r.Name != nil {
		p.Name = *r.Name
	}
	if r.Price != nil {
		p.Price = *r.Price
	}

	// Store updated product.
	err = s.setProductToHash(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *ProductStorage) Delete(ctx context.Context, id string) (*product.Product, error) {
	// Get product checking for existence.
	p, err := s.getProductFromHash(ctx, id)
	if err != nil {
		return nil, err
	}

	// Remove id from sorted set.
	err = s.rdb.ZRem(ctx, s.idsKey, id).Err()
	if err != nil {
		return nil, err
	}

	// Remove product hash.
	err = s.rdb.Del(ctx, s.hashKey(id)).Err()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *ProductStorage) setProductToHash(ctx context.Context, p *product.Product) error {
	err := s.rdb.HSet(
		ctx, s.hashKey(p.ID),
		"name", p.Name,
		"price", strconv.FormatInt(p.Price, 10),
		"seller", p.Seller,
	).Err()
	if err != nil {
		return fmt.Errorf("setting product hash: %w", err)
	}
	return nil
}

func (s *ProductStorage) getProductFromHash(ctx context.Context, id string) (*product.Product, error) {
	num, err := s.rdb.Exists(ctx, s.hashKey(id)).Result()
	if num == 0 {
		return nil, product.ErrNotFound
	}

	val, err := s.rdb.HMGet(ctx, s.hashKey(id), "name", "price", "seller").Result()
	if err != nil {
		return nil, err
	}

	name, ok := val[0].(string)
	if !ok {
		return nil, errors.New("nil name field in redis")
	}

	priceString, ok := val[1].(string)
	if !ok {
		return nil, errors.New("nil price field in redis")
	}
	price, err := strconv.ParseInt(priceString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing price: %w", err)
	}

	seller, ok := val[2].(string)
	if !ok {
		return nil, errors.New("nil seller field in redis")
	}

	p := &product.Product{
		ID:     id,
		Name:   name,
		Price:  price,
		Seller: seller,
	}
	return p, nil
}

func (s *ProductStorage) hashKey(id string) string {
	return fmt.Sprintf("%s:%s", s.baseKey, id)
}
