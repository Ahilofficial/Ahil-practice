package grpc

import (
	"backend_institutions/internal/loggerpb"
	"context"
)

type LoggerServer struct {
	loggerpb.UnimplementedLoggerServiceServer
}

func (s *LoggerServer) SaveLog(
	ctx context.Context,
	req *loggerpb.LogRequest,
) (*loggerpb.LogResponse, error) {

	return &loggerpb.LogResponse{
		Success: true,
		Message: "Log Saved",
	}, nil
}