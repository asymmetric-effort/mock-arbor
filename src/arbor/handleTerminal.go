package arbor

import (
	"fmt"
	"strconv"
	"strings"
)

func handleTerminal(cmd string) string {
	// terminal length <n>
	parts := strings.Fields(cmd)
	if len(parts) != 3 {
		return "Usage: terminal length <n>\n"
	}
	n, err := strconv.Atoi(parts[2])
	if err != nil || n < 0 {
		return "Invalid terminal length\n"
	}
	terminalLength = n
	return getFixture(
		"terminal/length",
		map[string]int{"N": terminalLength},
		fmt.Sprintf("Terminal length set to %d (emulated).\n", terminalLength))
}
