package grpc

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/ortymid/market/grpc/grpctest"
	"github.com/ortymid/market/grpc/pb"
	"github.com/ortymid/market/market/product"
	"github.com/ortymid/market/mock"
	"reflect"
	"testing"
)

type setupMocks func(as *mock.GRPCAuthService, ps *mock.ProductService)

func TestServer_List(t *testing.T) {
	tests := []struct {
		name       string
		req        *pb.FindRequest
		setupMocks setupMocks
		wantStream []*pb.ProductReply
		wantErr    bool
	}{
		{
			name: "Should stream products for offset=0 and limit=2",
			req: &pb.FindRequest{
				Offset: 0,
				Limit:  2,
			},
			setupMocks: func(as *mock.GRPCAuthService, ps *mock.ProductService) {
				ps.EXPECT().Find(
					gomock.Any(), product.FindRequest{Offset: 0, Limit: 2},
				).Return(
					[]*product.Product{
						{ID: "1", Name: "p1", Price: 100, Seller: "1"},
						{ID: "2", Name: "p2", Price: 200, Seller: "2"},
					},
					nil,
				)
			},
			wantStream: []*pb.ProductReply{
				{Id: "1", Name: "p1", Price: 100, Seller: "1"},
				{Id: "2", Name: "p2", Price: 200, Seller: "2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			as := mock.NewGRPCAuthService(ctrl)
			ps := mock.NewProductService(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(as, ps)
			}

			s := &Server{
				AuthService:    as,
				ProductService: ps,
			}

			stream := grpctest.NewProductService_ListRecorder()

			if err := s.Find(tt.req, stream); (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(stream.Stream, tt.wantStream) {
				t.Errorf("Find() stream = %v, wantStream %v", stream.Stream, tt.wantStream)
			}
		})
	}
}

func TestServer_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *pb.FindOneRequest
	}
	tests := []struct {
		name       string
		args       args
		setupMocks setupMocks
		want       *pb.ProductReply
		wantErr    bool
	}{
		{
			name: "Should reply with product",
			args: args{
				ctx: context.Background(),
				r:   &pb.FindOneRequest{Id: "1"},
			},
			setupMocks: func(as *mock.GRPCAuthService, ps *mock.ProductService) {
				ps.EXPECT().FindOne(gomock.Any(), "1").
					Return(&product.Product{ID: "1", Name: "p1", Price: 100, Seller: "1"}, nil)
			},
			want: &pb.ProductReply{Id: "1", Name: "p1", Price: 100, Seller: "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			as := mock.NewGRPCAuthService(ctrl)
			ps := mock.NewProductService(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(as, ps)
			}

			s := &Server{
				AuthService:    as,
				ProductService: ps,
			}
			got, err := s.FindOne(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *pb.CreateRequest
	}
	tests := []struct {
		name       string
		args       args
		setupMocks setupMocks
		want       *pb.ProductReply
		wantErr    bool
	}{
		{
			name: "Should create product",
			args: args{
				ctx: context.Background(),
				r:   &pb.CreateRequest{Name: "p1", Price: 100},
			},
			setupMocks: func(as *mock.GRPCAuthService, ps *mock.ProductService) {
				ps.EXPECT().Create(
					gomock.Any(),
					product.CreateRequest{
						Name:  "p1",
						Price: 100,
					}).
					Return(&product.Product{ID: "1", Name: "p1", Price: 100, Seller: "1"}, nil)
			},
			want: &pb.ProductReply{Id: "1", Name: "p1", Price: 100, Seller: "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			as := mock.NewGRPCAuthService(ctrl)
			ps := mock.NewProductService(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(as, ps)
			}

			s := &Server{
				AuthService:    as,
				ProductService: ps,
			}
			got, err := s.Create(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *pb.UpdateRequest
	}
	tests := []struct {
		name       string
		args       args
		setupMocks setupMocks
		want       *pb.ProductReply
		wantErr    bool
	}{
		{
			name: "Should update product name",
			args: args{
				ctx: context.Background(),
				r:   &pb.UpdateRequest{Id: "1", Name: testStringPtr("p2")},
			},
			setupMocks: func(as *mock.GRPCAuthService, ps *mock.ProductService) {
				ps.EXPECT().Update(
					gomock.Any(),
					product.UpdateRequest{
						ID:   "1",
						Name: testStringPtr("p2"),
					},
				).Return(&product.Product{ID: "1", Name: "p2", Price: 100, Seller: "1"}, nil)
			},
			want: &pb.ProductReply{Id: "1", Name: "p2", Price: 100, Seller: "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			as := mock.NewGRPCAuthService(ctrl)
			ps := mock.NewProductService(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(as, ps)
			}

			s := &Server{
				AuthService:    as,
				ProductService: ps,
			}
			got, err := s.Update(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *pb.DeleteRequest
	}
	tests := []struct {
		name       string
		args       args
		setupMocks setupMocks
		want       *pb.ProductReply
		wantErr    bool
	}{
		{
			name: "Should delete product",
			args: args{
				ctx: context.Background(),
				r:   &pb.DeleteRequest{Id: "1"},
			},
			setupMocks: func(as *mock.GRPCAuthService, ps *mock.ProductService) {
				ps.EXPECT().Delete(gomock.Any(), "1").
					Return(&product.Product{ID: "1", Name: "p2", Price: 100, Seller: "1"}, nil)
			},
			want: &pb.ProductReply{Id: "1", Name: "p2", Price: 100, Seller: "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			as := mock.NewGRPCAuthService(ctrl)
			ps := mock.NewProductService(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(as, ps)
			}

			s := &Server{
				AuthService:    as,
				ProductService: ps,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
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
