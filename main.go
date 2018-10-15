package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var version = "dev"

// SmudgeToml defines the top-level structure of the TOML file
type SmudgeToml struct {
	Tags   []string              `toml:"tags"`
	Chars  map[string]string     `toml:"chars"`
	Scenes map[string][][]string `toml:"scenes"`
}

var (
	// init subcommand
	initCmd         = kingpin.Command("init", "initalize")
	initCmdBasename = initCmd.Arg("basename", "Basename/shortname of the comic").Required().String()
	initCmdPages    = initCmd.Arg("num_pages", "Number of pages in the comic").Required().Int()

	// check subcommand
	checkCmd         = kingpin.Command("check", "check")
	checkCmdFiles    = checkCmd.Arg("file", "TOML file to be checked").Required().Strings()
	checkCmdTagsFile = checkCmd.Flag("tagfile", "TOML file containing list of allowed tags").Short('t').String()

	// quote subcommand
	quoteCmd       = kingpin.Command("quote", "generate quote group templates")
	quoteCmdGroups = quoteCmd.Arg("groups", "quote groups per page").Required().String()
	quoteCmdFlip   = quoteCmd.Flag("flip", "change external separator location to the bottom").Short('f').Bool()
	quoteCmdPrint  = quoteCmd.Flag("print", "print to stdout instead of copying to clipboard").Short('p').Bool()
)

func main() {
	kingpin.Version(version)
	switch kingpin.Parse() {
	case "init":
		initAction(*initCmdBasename, *initCmdPages)
	case "check":
		checkAction(*checkCmdFiles, *checkCmdTagsFile)
	case "quote":
		quoteAction(*quoteCmdGroups, *quoteCmdFlip, *quoteCmdPrint)
	}
}
