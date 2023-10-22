package gapi

import (
	db "techschool/db/sqlc"
	"techschool/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		PasswordChangeAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:        timestamppb.New(user.CreatedAt),
	}
}
