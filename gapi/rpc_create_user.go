package gapi

import (
	"context"
	db "techschool/db/sqlc"
	"techschool/pb"
	"techschool/util"
	"techschool/val"
	"techschool/worker"
	"time"

	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	password, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password %s", err.Error())
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: password,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			payloadDist := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			opts := []asynq.Option{
				asynq.Queue(worker.QueueCritical),
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
			}
			return s.taskDistributor.DistributeTaskSendVerifyEmail(ctx, payloadDist, opts...)

		},
	}

	txResult, err := s.store.CreateUserTx(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "error : %s", err.Error())
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(txResult.User),
	}
	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUserName(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("fullname", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
