package market_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ortymid/market/market"
	"github.com/ortymid/market/market/mock"
)

type MockFuncUser struct {
	expect     bool
	argID      string
	returnUser *market.User
	returnErr  error
}
type MockUserService struct {
	User MockFuncUser
}

func (opt *MockUserService) Setup(m *mock.MockUserService) {
	if opt.User.expect {
		m.EXPECT().User(gomock.Eq(opt.User.argID)).Return(opt.User.returnUser, opt.User.returnErr)
	} else {
		m.EXPECT().User(nil).MaxTimes(0)
	}
}

type MockFuncProducts struct {
	expect         bool
	argOffset      int
	argLimit       int
	returnProducts []*market.Product
	returnErr      error
}
type MockFuncProduct struct {
	expect        bool
	argID         int
	returnProduct *market.Product
	returnErr     error
}
type MockFuncAddProduct struct {
	expect        bool
	argRequest    market.AddProductRequest
	returnProduct *market.Product
	returnErr     error
}
type MockFuncEditProduct struct {
	expect        bool
	argRequest    market.EditProductRequest
	returnProduct *market.Product
	returnErr     error
}
type MockFuncDeleteProduct struct {
	expect    bool
	argID     int
	returnErr error
}
type MockProductService struct {
	Products      MockFuncProducts
	Product       MockFuncProduct
	AddProduct    MockFuncAddProduct
	EditProduct   MockFuncEditProduct
	DeleteProduct MockFuncDeleteProduct
}

func (opt *MockProductService) Setup(m *mock.MockProductService) {
	if opt.Products.expect {
		m.EXPECT().Products(context.TODO(), opt.Products.argOffset, opt.Products.argLimit).Return(opt.Products.returnProducts, opt.Products.returnErr)
	} else {
		m.EXPECT().Products(nil, nil, nil).MaxTimes(0)
	}
	if opt.Product.expect {
		m.EXPECT().Product(context.TODO(), opt.Product.argID).Return(opt.Product.returnProduct, opt.Product.returnErr)
	} else {
		m.EXPECT().Product(nil, nil).MaxTimes(0)
	}
	if opt.AddProduct.expect {
		m.EXPECT().AddProduct(context.TODO(), opt.AddProduct.argRequest).Return(opt.AddProduct.returnProduct, opt.AddProduct.returnErr)
	} else {
		m.EXPECT().AddProduct(nil, nil).MaxTimes(0)
	}
	if opt.EditProduct.expect {
		m.EXPECT().EditProduct(context.TODO(), opt.EditProduct.argRequest).Return(opt.EditProduct.returnProduct, opt.EditProduct.returnErr)
	} else {
		m.EXPECT().EditProduct(nil, nil).MaxTimes(0)
	}
	if opt.DeleteProduct.expect {
		m.EXPECT().DeleteProduct(context.TODO(), opt.DeleteProduct.argID).Return(opt.DeleteProduct.returnErr)
	} else {
		m.EXPECT().DeleteProduct(nil, nil).MaxTimes(0)
	}
}

func TestMarket_Products(t *testing.T) {
	type mocks struct {
		UserService    MockUserService
		ProductService MockProductService
	}
	type args struct {
		ctx    context.Context
		offset int
		limit  int
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    []*market.Product
		wantErr bool
	}{
		{
			name: "Returns products",
			mocks: mocks{
				ProductService: MockProductService{
					Products: MockFuncProducts{
						expect:    true,
						argOffset: 0,
						argLimit:  2,
						returnProducts: []*market.Product{
							{ID: 1, Name: "p1", Price: 100, Seller: "1"},
							{ID: 2, Name: "p2", Price: 200, Seller: "2"},
						},
					},
				},
			},
			args: args{
				ctx:    context.TODO(),
				offset: 0,
				limit:  2,
			},
			want: []*market.Product{
				{ID: 1, Name: "p1", Price: 100, Seller: "1"},
				{ID: 2, Name: "p2", Price: 200, Seller: "2"},
			},
		},
		{
			name: "Returns empty products",
			mocks: mocks{
				ProductService: MockProductService{
					Products: MockFuncProducts{
						argOffset:      0,
						argLimit:       2,
						expect:         true,
						returnProducts: []*market.Product{},
					},
				},
			},
			args: args{
				ctx:    context.TODO(),
				offset: 0,
				limit:  2,
			},
			want: []*market.Product{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			us := mock.NewMockUserService(ctrl)
			tt.mocks.UserService.Setup(us)

			ps := mock.NewMockProductService(ctrl)
			tt.mocks.ProductService.Setup(ps)

			m := &market.Market{
				UserService:    us,
				ProductService: ps,
			}
			got, err := m.Products(tt.args.ctx, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Market.Products() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Market.Products() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarket_Product(t *testing.T) {
	type mocks struct {
		UserService    MockUserService
		ProductService MockProductService
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    *market.Product
		wantErr bool
	}{
		{
			name: "Returns a product",
			mocks: mocks{
				ProductService: MockProductService{
					Product: MockFuncProduct{
						expect:        true,
						argID:         1,
						returnProduct: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			want: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
		},
		{
			name: "Returns an error for not existing product",
			mocks: mocks{
				ProductService: MockProductService{
					Product: MockFuncProduct{
						expect:    true,
						argID:     1,
						returnErr: market.ErrProductNotFound,
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				id: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			us := mock.NewMockUserService(ctrl)
			tt.mocks.UserService.Setup(us)

			ps := mock.NewMockProductService(ctrl)
			tt.mocks.ProductService.Setup(ps)

			m := &market.Market{
				UserService:    us,
				ProductService: ps,
			}
			got, err := m.Product(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Market.Products() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Market.Products() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarket_AddProduct(t *testing.T) {
	type mocks struct {
		UserService    MockUserService
		ProductService MockProductService
	}
	type args struct {
		ctx    context.Context
		req    market.AddProductRequest
		userID string
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    *market.Product
		wantErr bool
	}{
		{
			name: "Adds product",
			mocks: mocks{
				UserService: MockUserService{
					User: MockFuncUser{
						expect:     true,
						argID:      "1",
						returnUser: &market.User{ID: "1", Name: "u1"},
					},
				},
				ProductService: MockProductService{
					AddProduct: MockFuncAddProduct{
						expect:        true,
						argRequest:    market.AddProductRequest{Name: "p1", Price: 100, Seller: "1"},
						returnProduct: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
					},
				},
			},
			args: args{
				ctx:    context.TODO(),
				req:    market.AddProductRequest{Name: "p1", Price: 100, Seller: "1"},
				userID: "1",
			},
			want: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
		},
		{
			name: "Returns an error for not existing user",
			mocks: mocks{
				UserService: MockUserService{
					User: MockFuncUser{
						expect:    true,
						argID:     "1",
						returnErr: &market.ErrUserNotFound{},
					},
				},
			},
			args: args{
				ctx:    context.TODO(),
				req:    market.AddProductRequest{Name: "p1", Price: 100, Seller: "1"},
				userID: "1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			us := mock.NewMockUserService(ctrl)
			tt.mocks.UserService.Setup(us)

			ps := mock.NewMockProductService(ctrl)
			tt.mocks.ProductService.Setup(ps)

			m := &market.Market{
				UserService:    us,
				ProductService: ps,
			}
			got, err := m.AddProduct(tt.args.ctx, tt.args.req, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Market.Products() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Market.Products() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarket_EditProduct(t *testing.T) {
	type mocks struct {
		UserService    MockUserService
		ProductService MockProductService
	}
	type args struct {
		ctx    context.Context
		req    market.EditProductRequest
		userID string
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    *market.Product
		wantErr bool
	}{
		{
			name: "Should edit a product",
			mocks: mocks{
				UserService: MockUserService{
					User: MockFuncUser{
						expect:     true,
						argID:      "1",
						returnUser: &market.User{ID: "1", Name: "u1"},
					},
				},
				ProductService: MockProductService{
					EditProduct: MockFuncEditProduct{
						expect:        true,
						argRequest:    market.EditProductRequest{ID: 1, Name: "p2", Price: testIntPtr(200), Seller: "1"},
						returnProduct: &market.Product{ID: 1, Name: "p2", Price: 200, Seller: "1"},
					},
				},
			},
			args: args{
				ctx:    context.TODO(),
				req:    market.EditProductRequest{ID: 1, Name: "p2", Price: testIntPtr(200), Seller: "1"},
				userID: "1",
			},
			want: &market.Product{ID: 1, Name: "p2", Price: 200, Seller: "1"},
		},
		{
			name: "Returns an error for not existing user",
			mocks: mocks{
				UserService: MockUserService{
					User: MockFuncUser{
						expect:    true,
						argID:     "1",
						returnErr: &market.ErrUserNotFound{},
					},
				},
			},
			args: args{
				ctx:    context.TODO(),
				req:    market.EditProductRequest{Name: "p1", Price: testIntPtr(100), Seller: "1"},
				userID: "1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			us := mock.NewMockUserService(ctrl)
			tt.mocks.UserService.Setup(us)

			ps := mock.NewMockProductService(ctrl)
			tt.mocks.ProductService.Setup(ps)

			m := &market.Market{
				UserService:    us,
				ProductService: ps,
			}
			got, err := m.EditProduct(tt.args.ctx, tt.args.req, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Market.Products() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Market.Products() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarket_DeleteProduct(t *testing.T) {
	type mocks struct {
		UserService    MockUserService
		ProductService MockProductService
	}
	type args struct {
		ctx    context.Context
		id     int
		userID string
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    *market.Product
		wantErr bool
	}{
		{
			name: "Deletes product",
			mocks: mocks{
				ProductService: MockProductService{
					Product: MockFuncProduct{
						expect:        true,
						argID:         1,
						returnProduct: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
					},
					DeleteProduct: MockFuncDeleteProduct{
						expect: true,
						argID:  1,
					},
				},
			},
			args: args{
				ctx:    context.TODO(),
				id:     1,
				userID: "1",
			},
			want: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
		},
		{
			name: "Returns an error for not existing product",
			mocks: mocks{
				ProductService: MockProductService{
					Product: MockFuncProduct{
						expect:    true,
						argID:     1,
						returnErr: market.ErrProductNotFound,
					},
				},
			},
			args: args{
				ctx:    context.TODO(),
				id:     1,
				userID: "1",
			},
			wantErr: true,
		},
		{
			name: "Returns an error for user mismatch",
			mocks: mocks{
				ProductService: MockProductService{
					Product: MockFuncProduct{
						expect:        true,
						argID:         1,
						returnProduct: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
					},
				},
			},
			args: args{
				ctx:    context.TODO(),
				id:     1,
				userID: "2",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			us := mock.NewMockUserService(ctrl)
			tt.mocks.UserService.Setup(us)

			ps := mock.NewMockProductService(ctrl)
			tt.mocks.ProductService.Setup(ps)

			m := &market.Market{
				UserService:    us,
				ProductService: ps,
			}
			err := m.DeleteProduct(tt.args.ctx, tt.args.id, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Market.Products() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func testIntPtr(n int) *int {
	return &n
}
