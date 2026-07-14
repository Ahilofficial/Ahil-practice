package services

import (
	"backend_institutions/logger_service/internals/model"
	"backend_institutions/logger_service/internals/repository"
)

type LoggerService struct {
	loggerRepo *repository.LoggerRepo
}

func NewLoggerService(loggerRepo *repository.LoggerRepo) *LoggerService {
	return &LoggerService{
		loggerRepo: loggerRepo,
	}
}

func (s *LoggerService) SaveLog(log model.Log) error {
	return s.loggerRepo.WriteFile(log)
}