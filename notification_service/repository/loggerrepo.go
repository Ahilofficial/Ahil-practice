package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LoggerRepo struct{}

func NewLoggerRepo() *LoggerRepo {
	return &LoggerRepo{}
}

func (r *LoggerRepo) WriteLog(to, subject string, success bool, errorMsg string) error {
	
	logPath := "log/notification.log"
	if _, err := os.Stat("notification_service"); err == nil {
		fmt.Println("Cant able to find the path")
	}

	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directories: %w", err)
	}

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}


	status := "SUCCESS"
	detail := ""
	if !success {
		status = "FAILED"
		detail = fmt.Sprintf(" | Error: %s", errorMsg)
	}

	logEntry := fmt.Sprintf("[%s] Status: %s | To: %s | Subject: %s%s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		status,
		to,
		subject,
		detail,
	)

	fmt.Print(logEntry)
	if _, err := file.WriteString(logEntry); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return nil
}
