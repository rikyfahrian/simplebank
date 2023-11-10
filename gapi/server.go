package gapi

import (
	db "techschool/db/sqlc"
	"techschool/pb"
	"techschool/token"
	"techschool/util"
	"techschool/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config          *util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(store db.Store, config *util.Config, distributor worker.TaskDistributor) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: distributor,
	}

	return server, nil

}
