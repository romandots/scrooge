package utils

import (
	"errors"
	"fmt"
)

func log(level string, args ...any) string {
	format := args[0].(string)
	message := fmt.Sprintf(format, args[1:]...)
	fmt.Printf("[%s] %s\n", level, message)
	return message
}

func Info(args ...any) string {
	return log("INFO", args...)
}

func Debug(args ...any) string {
	return log("DEBUG", args...)
}

func Error(args ...any) error {
	return errors.New(log("ERROR", args...))
}
