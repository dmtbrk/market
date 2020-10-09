package http

import (
	"context"
	"github.com/ortymid/market/http/route"
	"github.com/ortymid/market/market/product"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

type Server struct {
	AuthService    AuthService
	ProductService product.Interface
}

func (s *Server) Handler() http.Handler {
	r := mux.NewRouter()

	// Product
	products := route.Product{ProductService: s.ProductService}
	products.Setup(r)

	// GraphQL
	gql := route.GraphQL{ProductService: s.ProductService}
	gql.Setup(r.PathPrefix("/gql/").Subrouter())

	// CORS
	h := cors.Default().Handler(r)

	// Auth
	h = AuthMiddleware(s.AuthService, h)

	return h
}

func (s *Server) Run(addr string) {
	httpServer := http.Server{
		Addr:    addr,
		Handler: s.Handler(),
	}

	wait := make(chan struct{})
	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-done
		log.Println("Gracefully stopping...")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Println("shutting down server:", err)
		}

		close(wait)
		log.Println("Server stopped.")
	}()

	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln("listening:", err)
		}
	}()
	log.Printf("Server started at %s.", httpServer.Addr)

	<-wait
}
