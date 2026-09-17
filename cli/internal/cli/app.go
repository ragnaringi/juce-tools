package cli

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func run(cmd *exec.Cmd) (bool, error) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return false, err
	}
	return true, nil
}

func Run(args []string) int {
	flags := flag.NewFlagSet("juce-tools", flag.ExitOnError)
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
