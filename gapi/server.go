package gapi

import (
	db "techschool/db/sqlc"
	"techschool/pb"
	"techschool/token"
	"techschool/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     *util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(store db.Store, config *util.Config) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil

}
