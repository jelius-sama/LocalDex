package logger

import (
	"fmt"
	"os"
	"time"
)

var LoggerStyle string = "brackets"

func SetStyle(s string) {
	switch s {
	case "brackets":
		LoggerStyle = "brackets"
		Okay("Logger style set to `" + LoggerStyle + "`.")
		return

	case "colon":
		LoggerStyle = "colon"
		Okay("Logger style set to `" + LoggerStyle + "`.")
		return

	default:
		LoggerStyle = "brackets"
		Warning("Logger style " + s + " does not exists, setting to default instead!")
		return
	}
}

func applyStyle(format string, label string) string {
	switch LoggerStyle {
	case "brackets":
		return fmt.Sprintf(format, "["+label+"]")

	case "colon":
		return fmt.Sprintf(format, label+":")

	default:
		Error("Unreachable code reached!")
		return fmt.Sprintf(format, "["+label+"]")
	}
}

func Error(a ...any) {
	fmt.Fprintln(os.Stderr, append(append([]any{applyStyle("\n\033[31m%s", "ERROR")}, a...), []any{"\033[0m"}...)...)
}

func Debug(a ...any) {
	fmt.Println(append(append([]any{applyStyle("\n\033[34m%s", "DEBUG")}, a...), []any{"\033[0m"}...)...)
}

func Panic(a ...any) {
	fmt.Println(append(append([]any{applyStyle("\n\033[31m%s", "PANIC")}, a...), []any{"\033[0m"}...)...)
	os.Exit(-1)
}

func Info(a ...any) {
	fmt.Println(append(append([]any{applyStyle("\n\033[0;36m%s", "INFO")}, a...), []any{"\033[0m"}...)...)
}

func Okay(a ...any) {
	fmt.Println(append(append([]any{applyStyle("\n\033[32m%s", "OK")}, a...), []any{"\033[0m"}...)...)
}

func Warning(a ...any) {
	fmt.Println(append(append([]any{applyStyle("\n\033[33m%s", "WARN")}, a...), []any{"\033[0m"}...)...)
}

func TimedError(a ...any) {
	Error(append([]any{time.Now().Format("2006/01/02 15:04:05")}, a...)...)
}

func TimedDebug(a ...any) {
	Debug(append([]any{time.Now().Format("2006/01/02 15:04:05")}, a...)...)
}

func TimedPanic(a ...any) {
	Panic(append([]any{time.Now().Format("2006/01/02 15:04:05")}, a...)...)
}

func TimedInfo(a ...any) {
	Info(append([]any{time.Now().Format("2006/01/02 15:04:05")}, a...)...)
}

func TimedOkay(a ...any) {
	Okay(append([]any{time.Now().Format("2006/01/02 15:04:05")}, a...)...)
}

func TimedWarning(a ...any) {
	Warning(append([]any{time.Now().Format("2006/01/02 15:04:05")}, a...)...)
}
