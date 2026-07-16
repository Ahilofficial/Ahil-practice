package grpc

import (
	"backend_institutions/internal/loggerpb"
	"backend_institutions/logger_service/internals/model"
	"backend_institutions/logger_service/internals/services"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type LoggerServer struct {
	loggerpb.UnimplementedLoggerServiceServer
	loggerService *services.LoggerService
}

func NewLoggerServer(loggerService *services.LoggerService) *LoggerServer {
	return &LoggerServer{
		loggerService: loggerService,
	}
}

func (s *LoggerServer) SaveLog(ctx context.Context, req *loggerpb.LogRequest) (*loggerpb.LogResponse, error) {

	logEntry := model.Log{
		Service:  req.Service,
		Method:   req.Method,
		Endpoint: req.Endpoint,
		Request:  req.Request,
		Response: req.Response,
		Status:   int(req.Status),
	}

	err := s.loggerService.SaveLog(logEntry)
	if err != nil {
		log.Printf("Failed to save log: %v", err)
		return &loggerpb.LogResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to save log: %v", err),
		}, nil
	}

	return &loggerpb.LogResponse{
		Success: true,
		Message: "Log Saved",
	}, nil
}

func StartGRPCServer(loggerService *services.LoggerService, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	loggerpb.RegisterLoggerServiceServer(grpcServer, NewLoggerServer(loggerService))
	reflection.Register(grpcServer)

	log.Printf("gRPC Logger Server starting on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}