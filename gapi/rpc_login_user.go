package gapi

import (
	"context"
	"database/sql"
	db "techschool/db/sqlc"
	"techschool/pb"
	"techschool/util"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "error : %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "error : %s", err.Error())
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error : %s", "incorrect password")
	}

	token, accessPayload, err := s.tokenMaker.CreateToken(user.Username, time.Duration(s.config.AccesTokenDuration*int(time.Minute)))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error : %s", "failed create token")
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.Username, time.Duration(s.config.RefreshTokenDuration*int(time.Hour)))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error : %s", "failed create refresh token")

	}

	metadata := s.extractMetaData(ctx)

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "error : %s", err.Error())
	}

	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           token,
		RefreshToken:          refreshToken,
		AccessTokenExpiredAt:  timestamppb.New(accessPayload.ExpiresAt),
		RefreshTokenExpiredAt: timestamppb.New(refreshPayload.ExpiresAt),
	}
	return rsp, nil

}
