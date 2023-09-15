package api

import "github.com/gin-gonic/gin"

func RouterSetup(server *Server) *gin.Engine {

	router := gin.Default()

	router.POST("/users/login", server.LoginUser)
	router.POST("/users", server.CreateUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getAllAccounts)
	router.POST("/transfers", server.createTransfer)

	return router

}
