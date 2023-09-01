package api

import (
	"log"
	db "techschool/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getAllAccounts)

	server.router = router
	return server

}

func (s *Server) Start(address string) {

	err := s.router.Run(address)
	if err != nil {
		log.Fatal(err)
	}

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
