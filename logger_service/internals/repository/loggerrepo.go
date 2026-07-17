package repository

import (
	"backend_institutions/logger_service/internals/model"
	"backend_institutions/utilities"
)

type LoggerRepo struct{}

func NewLoggerRepo()*LoggerRepo{
	return &LoggerRepo{}
}

func(l *LoggerRepo) WriteFile(log model.Log)error{
	return utilities.WriteAppLog(log.Service, log.Method, log.Endpoint, int(log.Status), log.Request, log.Response)
}

