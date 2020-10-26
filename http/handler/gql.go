package handler

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/ortymid/market/gql"
	"github.com/ortymid/market/gql/gen"
	"github.com/ortymid/market/market/product"
)

type GraphQL struct {
	ProductService product.Interface
}

// Setup registers all available routes under the provided *mux.Router.
func (g *GraphQL) Setup(r *mux.Router) {
	gqlSrv := handler.NewDefaultServer(gen.NewExecutableSchema(gen.Config{Resolvers: &gql.Resolver{
		ProductService: g.ProductService,
	}}))

	rt := r.Handle("/gql", gqlSrv)
	url, err := rt.URL()
	if err != nil {
		panic(fmt.Errorf("obtaining GraphQL handler url: %w", err))
	}

	r.Handle("/gql/play", playground.Handler("GQLP", url.Path))
}
