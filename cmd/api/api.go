package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AhmedRabea0302/ecom/services/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// Register user routes here
	userHandler := user.NewHandler()
	userHandler.RegisterRoutes(subrouter)

	log.Println("Listening on ", s.addr)

	return http.ListenAndServe(s.addr, router)
}
