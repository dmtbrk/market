package product_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/ortymid/market/market/auth"
	"github.com/ortymid/market/market/product"
	"github.com/ortymid/market/market/user"
	"github.com/ortymid/market/mock"
	"reflect"
	"testing"
)

type setupMocks func(m *mock.ProductStorage)

func TestService_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		r   product.CreateRequest
	}
	tests := []struct {
		name                    string
		args                    args
		setupMockProductStorage setupMocks
		wantP                   *product.Product
		wantErr                 bool
	}{
		{
			name: "Should create product",
			args: args{
				ctx: auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
				r:   product.CreateRequest{Name: "name", Price: 100},
			},
			setupMockProductStorage: func(m *mock.ProductStorage) {
				m.EXPECT().Create(
					auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
					product.CreateRequest{Name: "name", Price: 100, Seller: "1"},
				).Return(
					&product.Product{ID: "1", Name: "name", Price: 100, Seller: "1"},
					nil,
				)
			},
			wantP: &product.Product{ID: "1", Name: "name", Price: 100, Seller: "1"},
		},
		{
			name: "Should error when context without user",
			args: args{
				ctx: context.Background(),
				r:   product.CreateRequest{Name: "name", Price: 100},
			},
			wantErr: true,
		},
		{
			name: "Should error when storage returns error",
			args: args{
				ctx: auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
				r:   product.CreateRequest{Name: "name", Price: 100},
			},
			setupMockProductStorage: func(m *mock.ProductStorage) {
				m.EXPECT().Create(
					auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
					product.CreateRequest{Name: "name", Price: 100, Seller: "1"},
				).Return(
					nil,
					errors.New("test error"),
				)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mock.NewProductStorage(ctrl)
			if tt.setupMockProductStorage != nil {
				tt.setupMockProductStorage(storage)
			}

			s := &product.Service{
				Storage: storage,
			}
			gotP, err := s.Create(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotP, tt.wantP) {
				t.Errorf("Create() gotP = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name                    string
		args                    args
		setupMockProductStorage setupMocks
		want                    *product.Product
		wantErr                 bool
	}{
		{
			name: "Should delete product",
			args: args{
				ctx: auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
				id:  "1",
			},
			setupMockProductStorage: func(m *mock.ProductStorage) {
				m.EXPECT().Get(
					auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
					"1",
				).Return(
					&product.Product{ID: "1", Name: "name", Price: 100, Seller: "1"},
					nil,
				)

				m.EXPECT().Delete(
					auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
					"1",
				).Return(
					&product.Product{ID: "1", Name: "name", Price: 100, Seller: "1"},
					nil,
				)
			},
			want: &product.Product{ID: "1", Name: "name", Price: 100, Seller: "1"},
		},
		{
			name: "Should error when product not found",
			args: args{
				ctx: auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
				id:  "1",
			},
			setupMockProductStorage: func(m *mock.ProductStorage) {
				m.EXPECT().Get(
					auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
					"1",
				).Return(
					nil,
					product.ErrNotFound,
				)
			},
			wantErr: true,
		},
		{
			name: "Should error when context without user",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			wantErr: true,
		},
		{
			name: "Should error when user is not seller",
			args: args{
				ctx: auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
				id:  "1",
			},
			setupMockProductStorage: func(m *mock.ProductStorage) {
				m.EXPECT().Get(
					auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
					"1",
				).Return(
					&product.Product{ID: "1", Name: "name", Price: 100, Seller: "2"},
					nil,
				)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mock.NewProductStorage(ctrl)
			if tt.setupMockProductStorage != nil {
				tt.setupMockProductStorage(storage)
			}

			s := &product.Service{
				Storage: storage,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name                    string
		args                    args
		setupMockProductStorage setupMocks
		want                    *product.Product
		wantErr                 bool
	}{
		{
			name: "Should get product",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			setupMockProductStorage: func(m *mock.ProductStorage) {
				m.EXPECT().Get(
					context.Background(),
					"1",
				).Return(
					&product.Product{ID: "1", Name: "name", Price: 100, Seller: "1"},
					nil,
				)
			},
			want: &product.Product{ID: "1", Name: "name", Price: 100, Seller: "1"},
		},
		{
			name: "Should error when product not found",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			setupMockProductStorage: func(m *mock.ProductStorage) {
				m.EXPECT().Get(
					context.Background(),
					"1",
				).Return(
					nil,
					product.ErrNotFound,
				)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mock.NewProductStorage(ctrl)
			if tt.setupMockProductStorage != nil {
				tt.setupMockProductStorage(storage)
			}

			s := &product.Service{
				Storage: storage,
			}
			got, err := s.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	type args struct {
		ctx context.Context
		r   product.ListRequest
	}
	tests := []struct {
		name                    string
		args                    args
		setupMockProductStorage setupMocks
		want                    []*product.Product
		wantErr                 bool
	}{
		{
			name: "Should list products",
			args: args{
				ctx: context.Background(),
				r:   product.ListRequest{Offset: 2, Limit: 2},
			},
			setupMockProductStorage: func(m *mock.ProductStorage) {
				m.EXPECT().List(
					context.Background(),
					product.ListRequest{Offset: 2, Limit: 2},
				).Return(
					[]*product.Product{
						{ID: "1", Name: "name1", Price: 100, Seller: "1"},
						{ID: "2", Name: "name2", Price: 200, Seller: "2"},
					},
					nil,
				)
			},
			want: []*product.Product{
				{ID: "1", Name: "name1", Price: 100, Seller: "1"},
				{ID: "2", Name: "name2", Price: 200, Seller: "2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mock.NewProductStorage(ctrl)
			if tt.setupMockProductStorage != nil {
				tt.setupMockProductStorage(storage)
			}

			s := &product.Service{
				Storage: storage,
			}
			got, err := s.List(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		r   product.UpdateRequest
	}
	tests := []struct {
		name       string
		args       args
		setupMocks setupMocks
		want       *product.Product
		wantErr    bool
	}{
		{
			name: "Should update product",
			args: args{
				ctx: auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
				r: product.UpdateRequest{
					ID:    "1",
					Name:  testStringPtr("new name"),
					Price: testInt64Ptr(100),
				},
			},
			setupMocks: func(m *mock.ProductStorage) {
				m.EXPECT().Get(
					auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
					"1",
				).Return(
					&product.Product{ID: "1", Name: "name", Price: 10, Seller: "1"},
					nil,
				)

				m.EXPECT().Update(
					auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
					product.UpdateRequest{
						ID:    "1",
						Name:  testStringPtr("new name"),
						Price: testInt64Ptr(100),
					},
				).Return(
					&product.Product{ID: "1", Name: "new name", Price: 100, Seller: "1"},
					nil,
				)
			},
			want: &product.Product{ID: "1", Name: "new name", Price: 100, Seller: "1"},
		},
		{
			name: "Should error when user is not seller",
			args: args{
				ctx: auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
				r: product.UpdateRequest{
					ID:    "1",
					Name:  testStringPtr("new name"),
					Price: testInt64Ptr(100),
				},
			},
			setupMocks: func(m *mock.ProductStorage) {
				m.EXPECT().Get(
					auth.NewContextWithUser(context.Background(), &user.User{ID: "1"}),
					"1",
				).Return(
					&product.Product{ID: "1", Name: "name", Price: 10, Seller: "2"},
					nil,
				)
			},
			wantErr: true,
		},
		{
			name: "Should error when context without user",
			args: args{
				ctx: context.Background(),
				r: product.UpdateRequest{
					ID:    "1",
					Name:  testStringPtr("new name"),
					Price: testInt64Ptr(100),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mock.NewProductStorage(ctrl)
			if tt.setupMocks != nil {
				tt.setupMocks(storage)
			}

			s := &product.Service{
				Storage: storage,
			}
			got, err := s.Update(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func testStringPtr(s string) *string {
	return &s
}

func testInt64Ptr(i int64) *int64 {
	return &i
}
