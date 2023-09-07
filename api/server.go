package api

import (
	"log"
	db "techschool/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getAllAccounts)
	router.POST("/transfers", server.createTransfer)

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
