package utils

import (
	"fmt"
	"strings"
)

func HexToANSI(hex string) string {
	hex = strings.TrimPrefix(hex, "#")

	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)

	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

func PrintMessage(color, nickname, message string) {
	fmt.Printf("%s%s%s: %s", color, nickname, "\033[0m", message)
}
