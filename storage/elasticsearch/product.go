package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/ortymid/market/market/product"
)

type indexResponse struct {
	ID     string `json:"_id"`
	Result string `json:"result"`
}

type updateResponse struct {
	ID     string      `json:"_id"`
	Result string      `json:"result"`
	Get    getResponse `json:"get"`
}

type getResponse struct {
	Found  bool   `json:"found"`
	Source source `json:"_source"`
}

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
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
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
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("parsing elasticsearch response body: %w", err)
		}

		reason, _ := e["error"].(map[string]interface{})["reason"] // TODO: handle no reason situation

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
	req := esapi.GetRequest{
		Index:      s.index,
		DocumentID: id,
	}

	res, err := req.Do(ctx, s.es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, product.ErrNotFound
		}
		return nil, fmt.Errorf("elasticsearch: %s", res.Status())
	}

	var gr getResponse
	if err := json.NewDecoder(res.Body).Decode(&gr); err != nil {
		return nil, fmt.Errorf("parsing elasticseach response body: %w", err)
	}

	if !gr.Found {
		return nil, product.ErrNotFound
	}

	p := &product.Product{
		ID:     id,
		Name:   gr.Source.Name,
		Price:  gr.Source.Price,
		Seller: gr.Source.Seller,
	}
	return p, nil
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
		return nil, fmt.Errorf("making elasticsearch request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch: %s", res.Status())
	}

	var ir indexResponse
	if err := json.NewDecoder(res.Body).Decode(&ir); err != nil {
		return nil, fmt.Errorf("parsing elasticseach response body: %w", err)
	}

	if ir.Result != "created" {
		return nil, fmt.Errorf("not created, result: %s", ir.Result)
	}

	p := &product.Product{
		ID:     ir.ID,
		Name:   r.Name,
		Price:  r.Price,
		Seller: r.Seller,
	}
	return p, nil
}

func (s *ProductStorage) Update(ctx context.Context, r product.UpdateRequest) (*product.Product, error) {
	var buf bytes.Buffer
	b := map[string]interface{}{
		"doc": r,
	}
	err := json.NewEncoder(&buf).Encode(b)
	if err != nil {
		return nil, fmt.Errorf("encoding elasticsearch request: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:          s.index,
		DocumentID:     r.ID,
		Body:           &buf,
		SourceIncludes: []string{"name", "price", "seller"},
	}

	res, err := req.Do(ctx, s.es)
	if err != nil {
		return nil, fmt.Errorf("making elasticsearch request: %w", err)
	}

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, product.ErrNotFound
		}
		return nil, fmt.Errorf("elasticsearch: %s", res.Status())
	}

	var ur updateResponse
	if err := json.NewDecoder(res.Body).Decode(&ur); err != nil {
		return nil, fmt.Errorf("parsing elasticseach response body: %w", err)
	}

	p := &product.Product{
		ID:     ur.ID,
		Name:   ur.Get.Source.Name,
		Price:  ur.Get.Source.Price,
		Seller: ur.Get.Source.Seller,
	}
	return p, nil
}

func (s *ProductStorage) Delete(ctx context.Context, id string) (*product.Product, error) {
	p, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	req := esapi.DeleteRequest{
		Index:      s.index,
		DocumentID: id,
	}

	res, err := req.Do(ctx, s.es)
	if err != nil {
		return nil, fmt.Errorf("making elasticsearch request: %w", err)
	}

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, product.ErrNotFound
		}
		return nil, fmt.Errorf("elasticsearch: %s", res.Status())
	}

	return p, nil
}
