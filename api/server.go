package api

import (
	"log"
	db "techschool/db/sqlc"
	"techschool/token"
	"techschool/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	token  token.Maker
	router *gin.Engine
}

func NewServer(store db.Store) *Server {

	tokenMaker, err := token.NewPasetoMaker(util.RandomString(32))
	if err != nil {
		return nil
	}

	server := &Server{
		store: store,
		token: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.router = RouterSetup(server)
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
