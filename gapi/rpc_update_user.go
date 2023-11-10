package gapi

import (
	"context"
	"database/sql"
	db "techschool/db/sqlc"
	"techschool/pb"
	"techschool/util"
	"techschool/val"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	auth, err := s.authorizeUser(ctx)
	if err != nil {
		return nil, unauthentication(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if auth.Username != req.GetUsername() {
		return nil, status.Error(codes.PermissionDenied, "cannot update another users")
	}

	payload := db.UpdateUserParams{
		Username: req.GetUsername(),

		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {

		password, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password %s", err.Error())
		}

		payload.HashedPassword = sql.NullString{
			String: password,
			Valid:  true,
		}

		payload.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

	}

	user, err := s.store.UpdateUser(ctx, payload)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "error : %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "error : %s", err.Error())
	}

	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUserName(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if req.FullName != nil {
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("fullname", err))
		}
	}

	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	return violations
}
