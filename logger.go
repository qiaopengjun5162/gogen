package main

import (
	"fmt"
	"io"
)

func logInfo(output io.Writer, msg string) {
	_, _ = fmt.Fprintf(output, "%s[INFO] %s%s\n", colorCyan, msg, colorReset)
}

func logInput(output io.Writer, msg string) {
	_, _ = fmt.Fprintf(output, "%s[INPUT] %s%s", colorCyan, msg, colorReset)
}

func logProgress(output io.Writer, msg string) {
	_, _ = fmt.Fprintf(output, "%s[PROGRESS] %s%s\n", colorCyan, msg, colorReset)
}

func logSuccess(output io.Writer, msg string) {
	_, _ = fmt.Fprintf(output, "%s[SUCCESS] %s%s\n", colorGreen, msg, colorReset)
}

func logError(output io.Writer, err error) {
	_, _ = fmt.Fprintf(output, "%s[ERROR] %v%s\n", colorRed, err, colorReset)
}
