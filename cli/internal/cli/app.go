package cli

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var verbose bool

func formatCommand(cmd *exec.Cmd) string {
	var parts []string

	if cmd.Dir != "" {
		parts = append(parts, "cd", shellQuote(cmd.Dir), "&&")
	}

	for _, arg := range cmd.Args {
		parts = append(parts, shellQuote(arg))
	}

	return strings.Join(parts, " ")
}

func shellQuote(value string) string {
	if value == "" {
		return `""`
	}

	if !strings.ContainsAny(value, " \t\n\"'\\$&;()[]{}<>|*?!") {
		return value
	}

	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}

func run(cmd *exec.Cmd) (bool, error) {
	notice("$ %s", formatCommand(cmd))

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return false, err
		}

		return true, nil
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return false, err
	}
	return true, nil
}

func Run(args []string) int {
	flags := flag.NewFlagSet("juce-tools", flag.ExitOnError)
	flags.BoolVar(&verbose, "verbose", false, "stream command output")
	flags.Parse(args)

	if flags.NArg() == 0 {
		fmt.Println("No arguments found")
		return 1
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := runCommand(workingDirectory, flags.Args()); err != nil {
		fail("%v", err)
		return 1
	}

	return 0
}
