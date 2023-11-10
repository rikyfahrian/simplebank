package gapi

import (
	"context"
	"fmt"
	"strings"
	"techschool/token"

	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (s *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	authorization := md.Get(authorizationHeader)
	if len(authorization) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	value := authorization[0]

	fields := strings.Fields(value)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("type authorization is not match")
	}

	token := fields[1]
	payload, err := s.tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("token is not valid")
	}

	return payload, nil
}
