package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

const (
	grpcGatewayUserAgentHandler = "grpcgateway-user-agent"
	userAgentHeader             = "user-agent"
	xForwadedFor                = "x-forwarded-for"
)

func (s *Server) extractMetaData(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {

		if userAgent := md.Get(grpcGatewayUserAgentHandler); len(userAgent) > 0 {
			mtdt.UserAgent = userAgent[0]
		}

		if userAgentH := md.Get(userAgentHeader); len(userAgentH) > 0 {
			mtdt.UserAgent = userAgentH[0]

		}

		if clientIP := md.Get(xForwadedFor); len(clientIP) > 0 {
			mtdt.ClientIP = clientIP[0]
		}

	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
