package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	agentName := flag.String("agent", "", "Name of the agent to run")
	taskDescription := flag.String("task", "", "Task description for the agent")
	flag.Parse()

	if *agentName != "" {
		fmt.Printf("Running agent: %s\n", *agentName)
		if *taskDescription != "" {
			fmt.Printf("Task: %s\n", *taskDescription)
		}
		return
	}

	flag.Usage()
	os.Exit(0)
}
