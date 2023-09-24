package api

import (
	"os"
	db "techschool/db/sqlc"
	"techschool/util"
	"testing"

	"github.com/gin-gonic/gin"
)

func NewTestServer(t *testing.T, store db.Store) *Server {

	config := util.Config{
		TokenKey:           util.RandomString(32),
		AccesTokenDuration: 1,
	}

	return NewServer(store, &config)

}

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())

}
