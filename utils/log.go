package utils

import "fmt"

func log(level string, args ...any) {
	format := args[0].(string)
	message := fmt.Sprintf(format, args[1:]...)

	fmt.Printf("[%s] %s\n", level, message)
}

func Info(args ...any) {
	log("INFO", args...)
}

func Debug(args ...any) {
	log("DEBUG", args...)
}

func Error(args ...any) {
	log("ERROR", args...)
}
