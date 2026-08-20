package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "upload":
		if err := runUpload(os.Args[2:]); err != nil {
			log.Fatalf("upload: %v", err)
		}
	case "list":
		if err := runList(os.Args[2:]); err != nil {
			log.Fatalf("list: %v", err)
		}
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: cli <command> [flags]")
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  upload  add a new soundbite clip")
	fmt.Fprintln(os.Stderr, "  list    show all stored soundbites")
}
