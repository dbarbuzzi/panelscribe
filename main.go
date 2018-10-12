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
	initCommand     = kingpin.Command("init", "initalize")
	initBasenameArg = initCommand.Arg("basename", "Basename/shortname of the comic").Required().String()
	initPagesArg    = initCommand.Arg("num_pages", "Number of pages in the comic").Required().Int()

	// check subcommand
	checkCommand         = kingpin.Command("check", "check")
	checkTomlFiles       = checkCommand.Arg("file", "TOML file to be checked").Required().Strings()
	checkAllowedTagsFile = checkCommand.Flag("tagfile", "TOML file containing list of allowed tags").Short('t').String()
)

func main() {
	kingpin.Version(version)
	switch kingpin.Parse() {
	case "init":
		initAction(*initBasenameArg, *initPagesArg)
	case "check":
		checkAction(*checkTomlFiles, *checkAllowedTagsFile)
	}
}
