package logger

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func Info(args ...interface{}) {
	var msg string
	for _, arg := range args {
		msg += fmt.Sprintf("%v ", arg)
	}
	fmt.Printf("%s [INFO] %s\n", time.Now().Format(time.RFC3339), msg)
}

func InfoContext(ctx context.Context, args ...interface{}) {
	var msg string
	for _, arg := range args {
		msg += fmt.Sprintf("%v ", arg)
	}
	val := ctx.Value("request_id")
	if val == nil {
		val = ""
	}
	fmt.Printf("%s [INFO - %v] %s\n", time.Now().Format(time.RFC3339), val, msg)
}

func Error(err error) {
	fmt.Printf("%s [ERROR] %v\n", time.Now().Format(time.RFC3339), err)
}

func ErrorContext(ctx context.Context, err error) {
	val := ctx.Value("request_id")
	if val == nil {
		val = ""
	}
	fmt.Printf("%s [ERROR - %v] %v\n", time.Now().Format(time.RFC3339), val, err)
}

func Panic(args ...interface{}) {
	var msg string
	for _, arg := range args {
		msg += fmt.Sprintf(" %v", arg)
	}
	fmt.Printf("%s [PANIC] %v\n", time.Now().Format(time.RFC3339), msg)
}

func LogRequest(ctx context.Context, r *http.Request, args ...interface{}) {
	var msg string
	for _, arg := range args {
		msg += fmt.Sprintf("%v ", arg)
	}
	val := ctx.Value("request_id")
	if val == nil {
		val = ""
	}
	fmt.Printf("%s [LOG - %v] %s %s\n", time.Now().Format(time.RFC3339), val, r.URL.Path, msg)
}
