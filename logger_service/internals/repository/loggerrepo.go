package repository

import (
	"backend_institutions/logger_service/internals/model"
	"encoding/json"
	"fmt"
	"path/filepath"

	// "logger-service/internal/model"
	"os"
	"time"
)
type LoggerRepo struct{}

func NewLoggerRepo ()*LoggerRepo{
	return &LoggerRepo{}
}
func(l *LoggerRepo) WriteFile(log model.Log)error{
	
	path := "backend_institutions/logger_service/internals/logs/app.log"


	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	open_create, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	fmt.Println("Successfully Opened files")
	requestbody,err:=json.MarshalIndent(log.Request,"","")
	if err!=nil{
		return err
	}
	responsebody,err:=json.MarshalIndent(log.Response,"","")
	if err!=nil{
		return err
	}

	logContent := fmt.Sprintf(`
**
Time      : %s
Service   : %s
Method    : %s
Endpoint  : %s
Status    : %d

Request:
%s

Response:
%s

**

`,
time.Now(),
log.Service,
log.Method,
log.Endpoint,
log.Status,
string(requestbody),
string(responsebody),
)


fmt.Println(logContent)
	n, err := open_create.WriteString(logContent)
	fmt.Println("Bytes written:", n)


	if err != nil {
		return err
	}

	return nil
}

