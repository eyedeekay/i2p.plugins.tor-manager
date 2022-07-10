package main

import (
	"log"
	"os"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
)

func CleanupArgs() (args []string, trailers []string) {
	// get a list of all possible flags from the flag package
	for i, arg := range os.Args[1:] {
		log.Printf("arg %d: %s", i, arg)
		trailer := true
		trimmed := strings.Split(strings.TrimLeft(arg, "-"), "=")[0]
		flag.VisitAll(func(f *flag.Flag) {
			//log.Printf("testing f.Name: %s against %s", f.Name, trimmed)
			if f.Name == trimmed {
				log.Printf("found flag: %s", f.Name)
				// if the next arg is not a flag(does not start with '-'), add it to the end of arg before
				// appending arg to args
				if i+1 < len(os.Args) && os.Args[i+1][0] != '-' {
					args = append(args, arg+"="+os.Args[i+1])
					i++
					trailer = false
				} else {
					args = append(args, arg)
					trailer = false
				}
				return
			}
		})
		if trailer {
			if arg != "--profile" && arg != "-P" && arg != "-profile" {
				log.Printf("found possible firefox flag: %s", trimmed)
				trailers = append(trailers, arg)
			} else {
				i++
			}
		}
		time.Sleep(time.Second)
	}
	return
}
