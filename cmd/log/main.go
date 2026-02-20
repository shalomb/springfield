package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shalomb/springfield/pkg/logger"
)

func main() {
	agent := flag.String("agent", "unknown", "Agent name")
	level := flag.String("level", "INFO", "Log level")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: log [--agent name] [--level level] message")
		os.Exit(1)
	}

	message := flag.Arg(0)
	if err := logger.Log(message, *level, *agent, "", "", nil, 0, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error logging message: %v\n", err)
		os.Exit(1)
	}
}
