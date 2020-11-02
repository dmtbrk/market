package elasticsearch

import (
	"github.com/ortymid/market/market/product"
	"reflect"
	"testing"
)

func Test_makeSearchQuery(t *testing.T) {
	type args struct {
		r product.FindRequest
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "Should make query with name",
			args: args{r: product.FindRequest{
				Offset:     0,
				Limit:      10,
				Name:       testPtrString("name"),
				PriceRange: nil,
				Seller:     nil,
			}},
			want: map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []interface{}{
						map[string]interface{}{
							"match": map[string]interface{}{
								"name": map[string]interface{}{
									"query":     "name",
									"fuzziness": "AUTO",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Should make query with price range",
			args: args{r: product.FindRequest{
				Offset: 0,
				Limit:  10,
				Name:   nil,
				PriceRange: &product.PriceRange{
					From: testPtrInt64(10),
					To:   testPtrInt64(100),
				},
				Seller: nil,
			}},
			want: map[string]interface{}{
				"bool": map[string]interface{}{
					"filter": []interface{}{
						map[string]interface{}{
							"range": map[string]interface{}{
								"price": map[string]interface{}{
									"gte": int64(10),
									"lte": int64(100),
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeSearchQuery(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeSearchQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testPtrString(v string) *string {
	return &v
}

func testPtrInt64(v int64) *int64 {
	return &v
}
