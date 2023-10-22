package gapi

import (
	"context"
	db "techschool/db/sqlc"
	"techschool/pb"
	"techschool/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	password, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password %s", err.Error())
	}

	payload := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: password,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := s.store.CreateUser(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "error : %s", err.Error())
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}
