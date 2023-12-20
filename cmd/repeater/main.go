package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	n = flag.Int("n", -1, "max iterations, -1 (default) for infinite")
	v = flag.Bool("v", false, "verbose")
	s = flag.String("s", "", "time to wait between iterations, default is no sleep")
)

func main() {
	flag.Parse()

	args := flag.Args()
	cmd := args[0]
	args = args[1:]

	var sleep time.Duration
	if *s != "" {
		var err error
		sleep, err = time.ParseDuration(*s)
		if err != nil {
			fmt.Printf("Failed to parse duration: %s\n", err)
			return
		}
	}

	// Repeat logic
	repeat(cmd, args, *n, sleep)
}

func repeat(command string, args []string, n int, s time.Duration) {
	c := command + " " + strings.Join(args, " ")
	i := 0

	for n < 0 || i < n {
		if *v {
			fmt.Printf("Executing command (i=%d): %s\n", i, c)
		}

		cmd := exec.Command(command, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to execute command: %s\n", err)
			break
		}

		i++

		if s != 0 {
			time.Sleep(s)
		}
	}
}
