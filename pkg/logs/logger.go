package logs

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	file *os.File
}

func NewLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{file: file}, nil
}

func (l *Logger) LogRequest(url string, method string, statusCode int, payload interface{}, response interface{}) error {
	timestamp := time.Now().Format("2006/01/02 - 15:04:05")
	logEntry := fmt.Sprintf("[%s] (%s) %s - %d \nPayload: %s \nResponse: %s\n\n",
		timestamp,
		method,
		url,
		statusCode,
		payload,
		fmt.Sprintf("%+v", response),
	)

	_, err := l.file.WriteString(logEntry)
	return err
}

func (l *Logger) LogError(message string, err error) error {
	timestamp := time.Now().Format("2006/01/02 - 15:04:05")
	logEntry := fmt.Sprintf("[%s] ERROR: %s - %v\n\n\n", timestamp, message, err)
	_, err = l.file.WriteString(logEntry)
	return err
}

func (l *Logger) Close() error {
	return l.file.Close()
}
