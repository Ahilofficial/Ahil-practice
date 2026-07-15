package repository

import (
	"backend_institutions/logger_service/internals/model"
	"encoding/json"
	"fmt"

	// "logger-service/internal/model"
	"os"
	"time"
)
type LoggerRepo struct{}

func NewLoggerRepo ()*LoggerRepo{
	return &LoggerRepo{}
}
func(l *LoggerRepo) WriteFile(log model.Log)error{
	open_create,err:=os.OpenFile("log/app.log",os.O_APPEND | os.O_CREATE |os.O_WRONLY,0644)
	fmt.Println("Successfully Opened files")
	if err!=nil{
		return err
	}
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
	_, err = open_create.WriteString(logContent)
	if err != nil {
		return err
	}

	return nil

}