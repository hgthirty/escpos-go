package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"os/exec"
)

func test() {
	command := flag.String("cmd", "pwd", "Set the command.")
	args := flag.String("args", "", "Set the args. (separated by spaces)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-cmd <command>] [-args <the arguments (separated by spaces)>]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	fmt.Println("Command: ", *command)
	fmt.Println("Arguments: ", *args)
	var argArray []string
	if *args != "" {
		argArray = strings.Split(*args, " ")
	} else {
		argArray = make([]string, 0)
	}
	cmd := exec.Command(*command, argArray...)
	buf, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "The command failed to perform: %s (Command: %s, Arguments: %s)", err, *command, *args)
		return
	}
	fmt.Fprintf(os.Stdout, "Result: %s", buf)
}
