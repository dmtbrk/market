package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/ortymid/market/http/handler"
	"github.com/ortymid/market/market"
)

type Server struct {
	Market    market.Interface
	JWTAlg    string
	JWTSecret interface{}

	httpSrv *http.Server
}

func NewServer(port int, market market.Interface) *Server {
	httpSrv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	s := Server{
		Market: market,
	}
	s.setupHTTP(httpSrv)

	return &s
}

func (s *Server) setupHTTP(httpSrv *http.Server) {
	s.setupHandler(httpSrv)
	s.httpSrv = httpSrv
}

func (s *Server) setupHandler(httpSrv *http.Server) {
	r := mux.NewRouter()

	sr := r.PathPrefix("/products").Subrouter()
	productHandler := &handler.ProductHandler{
		Market: s.Market,
	}
	productHandler.Setup(sr)

	httpSrv.Handler = handler.JWTMiddleware(r, s.JWTAlg, s.JWTSecret)
}

// Run is a convenient function to start an http server with graceful shotdown.
func (s *Server) Run() {
	idle := make(chan struct{})
	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-done
		log.Println("Gracefully stopping...")
		if err := s.httpSrv.Shutdown(context.Background()); err != nil {
			log.Println("server Shutdown:", err)
		}
		close(idle)
		log.Println("Server stopped")
	}()

	go func() {
		if err := s.httpSrv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln("server ListenAndServe:", err)
		}
	}()
	log.Print("Server started at ", s.httpSrv.Addr)

	<-idle
}
