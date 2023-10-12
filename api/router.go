package api

import "github.com/gin-gonic/gin"

func (server *Server) RouterSetup() {

	router := gin.Default()

	router.POST("/users/login", server.LoginUser)
	router.POST("/users", server.CreateUser)
	router.POST("/tokens/renew", server.renewAccessToken)

	auth := router.Group("/").Use(authMiddleware(server.tokenMaker))

	auth.POST("/accounts", server.createAccount)
	auth.GET("/accounts/:id", server.getAccount)
	auth.GET("/accounts", server.getAllAccounts)

	auth.POST("/transfers", server.createTransfer)

	server.router = router

}
