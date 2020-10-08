package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ortymid/market/market/product"
)

type ProductStorage struct {
	db    *sql.DB
	table string
}

func NewProductStorage(db *sql.DB, table string) *ProductStorage {
	return &ProductStorage{db: db, table: table}
}

func (s *ProductStorage) List(ctx context.Context, r product.ListRequest) ([]*product.Product, error) {
	query := fmt.Sprintf(
		`SELECT id, name, price, seller FROM %s LIMIT $1 OFFSET $2`,
		s.table,
	)

	rows, err := s.db.QueryContext(ctx, query, r.Limit, r.Offset)
	if err != nil {
		return nil, err
	}

	ps := make([]*product.Product, 0, r.Limit)
	for rows.Next() {
		var id, name, seller string
		var price int64
		if err := rows.Scan(&id, &name, &price, &seller); err != nil {
			return ps, err
		}

		ps = append(ps, &product.Product{
			ID:     id,
			Name:   name,
			Price:  price,
			Seller: seller,
		})
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ps, nil
}

func (s *ProductStorage) Get(ctx context.Context, id string) (p *product.Product, err error) {
	query := fmt.Sprintf(
		`SELECT id, name, price, seller FROM %s WHERE id = $1`,
		s.table,
	)

	var name, seller string
	var price int64
	err = s.db.QueryRowContext(ctx, query, id).Scan(&id, &name, &price, &seller)
	if err != nil {
		return p, err
	}

	return &product.Product{
		ID:     id,
		Name:   name,
		Price:  price,
		Seller: seller,
	}, nil
}

func (s *ProductStorage) Create(ctx context.Context, r product.CreateRequest) (p *product.Product, err error) {
	query := fmt.Sprintf(
		`INSERT INTO %s (name, price, seller) VALUES ($1, $2, $3) RETURNING id, name, price, seller`,
		s.table,
	)

	var id, name, seller string
	var price int64
	err = s.db.QueryRowContext(ctx, query, r.Name, r.Price, r.Seller).Scan(&id, &name, &price, &seller)
	if err != nil {
		return p, err
	}

	return &product.Product{
		ID:     id,
		Name:   name,
		Price:  price,
		Seller: seller,
	}, nil
}

func (s *ProductStorage) Update(ctx context.Context, r product.UpdateRequest) (*product.Product, error) {
	p, err := s.Get(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	if r.Name != nil {
		p.Name = *r.Name
	}
	if r.Price != nil {
		p.Price = *r.Price
	}

	query := fmt.Sprintf(
		`UPDATE %s SET name = $2, price = $3 WHERE id = $1 RETURNING id, name, price, seller`,
		s.table,
	)
	err = s.db.QueryRowContext(ctx, query, p.ID, p.Name, p.Price).Scan(&p.ID, &p.Name, &p.Price, &p.Seller)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s *ProductStorage) Delete(ctx context.Context, id string) (p *product.Product, err error) {
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE id = $1 RETURNING id, name, price, seller`,
		s.table,
	)

	var name, seller string
	var price int64
	err = s.db.QueryRowContext(ctx, query, id).Scan(&id, &name, &price, &seller)
	if err != nil {
		return p, err
	}

	return &product.Product{
		ID:     id,
		Name:   name,
		Price:  price,
		Seller: seller,
	}, nil
}
