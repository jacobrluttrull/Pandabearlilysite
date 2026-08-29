package main

import (
	"fmt"
	"log"
	"os"
)

// commands maps each subcommand name to its handler and the one-line description shown
// by `cli` with no arguments. Adding a command means adding one entry here.
var commands = []struct {
	name    string
	summary string
	run     func([]string) error
}{
	{"upload", "add one clip (length measured, name derived)", runUpload},
	{"import", "bulk-import a folder of .mp3 clips", runImport},
	{"names", "seed/refresh the display-name overrides file", runNames},
	{"apply-names", "push edited names into the database", runApplyNames},
	{"rename", "relabel one clip (updates names.json too)", runRename},
	{"set-date", "record when a clip was made", runSetDate},
	{"dedupe", "find and merge byte-identical duplicate clips", runDedupe},
	{"remove", "delete clips by filename", runRemove},
	{"reset-plays", "zero play tallies (-all, or by filename)", runResetPlays},
	{"sync-clips", "upload local clips the store does not have yet", runSyncClips},
	{"check", "audit database, audio files, and names for drift", runCheck},
	{"list", "show all stored soundbites", runList},
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	for _, cmd := range commands {
		if cmd.name != os.Args[1] {
			continue
		}
		if err := cmd.run(os.Args[2:]); err != nil {
			log.Fatalf("%s: %v", cmd.name, err)
		}
		return
	}

	fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
	usage()
	os.Exit(1)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: cli <command> [flags]")
	fmt.Fprintln(os.Stderr, "\ncommands:")
	for _, cmd := range commands {
		fmt.Fprintf(os.Stderr, "  %-12s %s\n", cmd.name, cmd.summary)
	}
	fmt.Fprintln(os.Stderr, "\nrun `cli <command> -h` for that command's flags")
}
