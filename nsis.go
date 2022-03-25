package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func NSISCompat() error {
	// check if no flags(args beginning with '-') were passed
	if len(os.Args) == 1 {
		return nil
	}
	for _, arg := range os.Args[1:] {
		if arg[0] == '-' {
			return nil
		}
	}
	var convertArgs = []string{}
	// check if any args beginning with /S or /D
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "/S") {
			log.Println("/S flag was passed, we're operating in NSIS installer compatibility mode")
		}
		if strings.HasPrefix(arg, "/D") {
			if len(arg) > 3 {
				return fmt.Errorf("/D flag was passed with a directory argument in NSIS compatibility mode")
			}
			log.Println("/D flag was passed, we're operating in NSIS installer compatibility mode")
			convertArgs = append(convertArgs, "--directory="+arg[2:])
		}
	}
	if len(convertArgs) == 0 {
		return nil
	}
	// if we're here, we're operating in NSIS compatibility mode
	// re-run ourselves with the converted args
	// start by getting the path to our executable
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	// forumulate our new command
	cmd := exec.Command(exePath, convertArgs...)
	// set the current working directory to the same as our executable
	cmd.Dir = filepath.Dir(exePath)
	// run the command
	err = cmd.Run()
	if err != nil {
		return err
	}
	// if we're here, we're done
	return nil
}
