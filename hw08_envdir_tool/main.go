package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("usage: %s /path/to/env/dir command arg1 arg2", filepath.Base(os.Args[0]))
	}
	envDir := os.Args[1]
	command := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(RunCmd(command, env))
}
