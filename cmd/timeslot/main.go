package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Melpic13/timeslot"
)

func main() {
	version := flag.Bool("version", false, "print timeslot version")
	flag.Parse()

	if *version {
		fmt.Println(timeslot.Version)
		return
	}

	fmt.Fprintln(os.Stderr, "timeslot CLI")
	fmt.Fprintln(os.Stderr, "Use --version to print the current version.")
	os.Exit(2)
}
