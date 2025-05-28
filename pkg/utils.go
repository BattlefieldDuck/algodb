package pkg

import (
	"fmt"
	"time"
)

// Printf prints a timestamped message.
func Printf(format string, args ...any) {
	ts := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] %s", ts, fmt.Sprintf(format, args...))
}
