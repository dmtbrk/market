package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/ortymid/market/market/product"
)

type jsonObject map[string]interface{}

type searchResponse struct {
	Hits hits `json:"hits"`
}

type hits struct {
	Hits []hit `json:"hits"`
}

type hit struct {
	ID     string `json:"_id"`
	Source source `json:"_source"`
}

type source struct {
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	Seller string `json:"seller"`
}

type ProductStorage struct {
	es    *elasticsearch.Client
	index string
}

func NewProductStorage(es *elasticsearch.Client, index string) *ProductStorage {
	return &ProductStorage{es: es, index: index}
}

func (s *ProductStorage) List(ctx context.Context, r product.ListRequest) ([]*product.Product, error) {
	var b bytes.Buffer
	query := jsonObject{
		"query": jsonObject{
			"match_all": jsonObject{},
		},
		"from": r.Offset,
		"size": r.Limit,
	}
	if err := json.NewEncoder(&b).Encode(query); err != nil {
		return nil, fmt.Errorf("encoding elasticsearch query: %w", err)
	}

	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(s.index),
		s.es.Search.WithBody(&b),
	)
	if err != nil {
		return nil, fmt.Errorf("searching: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e jsonObject
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("parsing elasticsearch response body: %w", err)
		}

		reason, _ := e["error"].(jsonObject)["reason"] // TODO: handle no reason situation

		return nil, fmt.Errorf("searching: %s", reason)
	}

	var sr searchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("parsing elasticseach response body: %w", err)
	}

	var ps []*product.Product
	for _, hit := range sr.Hits.Hits {
		p := &product.Product{
			ID:     hit.ID,
			Name:   hit.Source.Name,
			Price:  hit.Source.Price,
			Seller: hit.Source.Seller,
		}

		ps = append(ps, p)
	}

	return ps, nil
}

func (s *ProductStorage) Get(ctx context.Context, id string) (*product.Product, error) {
	panic("implement me")
}

func (s *ProductStorage) Create(ctx context.Context, r product.CreateRequest) (*product.Product, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req := esapi.IndexRequest{
		Index: s.index,
		Body:  bytes.NewReader(b),
	}

	res, err := req.Do(ctx, s.es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New("cannot create product") // TODO: add elasticseach error description
	}

	var rm map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("parsing elasticseach response body: %w", err)
	}

	if rm["result"] != "created" {
		return nil, errors.New("not created")
	}

	id, _ := rm["_id"].(string) // TODO: handle no id situation

	p := &product.Product{
		ID:     id,
		Name:   r.Name,
		Price:  r.Price,
		Seller: r.Seller,
	}
	return p, nil
}

func (s *ProductStorage) Update(ctx context.Context, r product.UpdateRequest) (*product.Product, error) {
	panic("implement me")
}

func (s *ProductStorage) Delete(ctx context.Context, id string) (*product.Product, error) {
	panic("implement me")
}
