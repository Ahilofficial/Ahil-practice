package services

import (
	"backend_institutions/internal/loggerpb"
	"backend_institutions/logger_service/internals/model"
	"backend_institutions/logger_service/internals/repository"
	"context"
	"fmt"
	"log"
)

type LoggerService struct {
	loggerpb.UnimplementedLoggerServiceServer
	loggerRepo *repository.LoggerRepo
}

func NewLoggerService(loggerRepo *repository.LoggerRepo) *LoggerService {
	return &LoggerService{
		loggerRepo: loggerRepo,
	}
}

func (s *LoggerService) SaveLog(ctx context.Context, req *loggerpb.LogRequest) (*loggerpb.LogResponse, error) {
	logEntry := model.Log{
		Service:  req.Service,
		Method:   req.Method,
		Endpoint: req.Endpoint,
		Request:  req.Request,
		Response: req.Response,
		Status:   int(req.Status),
	}

	err := s.loggerRepo.WriteFile(logEntry)
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