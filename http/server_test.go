package http

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/ortymid/market/market/product"
	"github.com/ortymid/market/market/user"
	"github.com/ortymid/market/mock"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type setupMocks func(as *mock.HTTPAuthService, ps *mock.ProductService)

func TestServer(t *testing.T) {
	tests := []struct {
		name       string
		req        *http.Request
		setupMocks setupMocks
		wantStatus int
		wantBody   []byte
	}{
		// GET /products/
		{
			name: "Should return products for offset=0 and limit=2",
			req:  httptest.NewRequest(http.MethodGet, "/products/?offset=0&limit=2", nil),
			setupMocks: func(as *mock.HTTPAuthService, ps *mock.ProductService) {
				as.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(nil, nil)

				ps.EXPECT().Find(
					gomock.Any(),
					product.FindRequest{Offset: 0, Limit: 2},
				).Return(
					[]*product.Product{
						{ID: "1", Name: "p1", Price: 100, Seller: "1"},
						{ID: "2", Name: "p2", Price: 200, Seller: "2"},
					},
					nil,
				)
			},
			wantStatus: http.StatusOK,
			wantBody: testBody([]*product.Product{
				{ID: "1", Name: "p1", Price: 100, Seller: "1"},
				{ID: "2", Name: "p2", Price: 200, Seller: "2"},
			}),
		},

		// GET /products/{id}
		{
			name: "Should return product",
			req:  httptest.NewRequest(http.MethodGet, "/products/1", nil),
			setupMocks: func(as *mock.HTTPAuthService, ps *mock.ProductService) {
				as.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(nil, nil)

				ps.EXPECT().FindOne(
					gomock.Any(),
					"1",
				).Return(
					&product.Product{ID: "1", Name: "p1", Price: 100, Seller: "1"},
					nil,
				)
			},
			wantStatus: http.StatusOK,
			wantBody:   testBody(&product.Product{ID: "1", Name: "p1", Price: 100, Seller: "1"}),
		},

		// POST /products/
		{
			name: "Should create product",
			req: httptest.NewRequest(
				http.MethodPost,
				"/products/",
				bytes.NewReader(testBody(product.CreateRequest{
					Name:  "p1",
					Price: 100,
				})),
			),
			setupMocks: func(as *mock.HTTPAuthService, ps *mock.ProductService) {
				as.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(&user.User{ID: "1"}, nil)

				ps.EXPECT().Create(
					gomock.Any(),
					product.CreateRequest{
						Name:  "p1",
						Price: 100,
					},
				).Return(
					&product.Product{ID: "1", Name: "p1", Price: 100, Seller: "1"},
					nil,
				)
			},
			wantStatus: http.StatusOK,
			wantBody:   testBody(&product.Product{ID: "1", Name: "p1", Price: 100, Seller: "1"}),
		},

		// PATCH /products/{id}
		{
			name: "Should update product name",
			req: httptest.NewRequest(
				http.MethodPatch,
				"/products/1",
				bytes.NewReader(testBody(product.UpdateRequest{
					Name: testStringPtr("p2"),
				})),
			),
			setupMocks: func(as *mock.HTTPAuthService, ps *mock.ProductService) {
				as.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(&user.User{ID: "1"}, nil)

				ps.EXPECT().Update(
					gomock.Any(),
					product.UpdateRequest{
						ID:   "1",
						Name: testStringPtr("p2"),
					},
				).Return(
					&product.Product{ID: "1", Name: "p2", Price: 100, Seller: "1"},
					nil,
				)
			},
			wantStatus: http.StatusOK,
			wantBody:   testBody(&product.Product{ID: "1", Name: "p2", Price: 100, Seller: "1"}),
		},

		// DELETE /products/{id}
		{
			name: "Should delete product",
			req: httptest.NewRequest(
				http.MethodDelete,
				"/products/1",
				nil,
			),
			setupMocks: func(as *mock.HTTPAuthService, ps *mock.ProductService) {
				as.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(&user.User{ID: "1"}, nil)

				ps.EXPECT().Delete(
					gomock.Any(),
					"1",
				).Return(
					&product.Product{ID: "1", Name: "p1", Price: 100, Seller: "1"},
					nil,
				)
			},
			wantStatus: http.StatusOK,
			wantBody:   testBody(&product.Product{ID: "1", Name: "p1", Price: 100, Seller: "1"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			as := mock.NewHTTPAuthService(ctrl)
			ps := mock.NewProductService(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(as, ps)
			}

			s := &Server{
				AuthService:    as,
				ProductService: ps,
			}

			w := httptest.NewRecorder()
			s.Handler().ServeHTTP(w, tt.req)
			res := w.Result()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("got status %v, want %v", res.StatusCode, tt.wantStatus)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("error reading response body: %v", err)
			}
			if !reflect.DeepEqual(body, tt.wantBody) {
				t.Errorf("got body %q, want %q", body, tt.wantBody)
			}
		})
	}
}

func testBody(v interface{}) []byte {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(v)
	if err != nil {
		log.Fatalf("error encoding test body: %v", err)
	}
	return b.Bytes()
}

func testStringPtr(s string) *string {
	return &s
}

func testInt64Ptr(i int64) *int64 {
	return &i
}
