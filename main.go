package main

import (
	"log"

	"gopkg.in/alecthomas/kingpin.v2"
)

var version = "dev"

// SmudgeTomlMeta defines the "meta" block in the smudge toml files
type SmudgeTomlMeta struct {
	Color     bool     `toml:"color"`
	Con       bool     `toml:"con"`
	Pairings  []string `toml:"pairings"`
	Themes    []string `toml:"themes"`
	Locations []string `toml:"locations"`
}

// NewSmudgeMeta creates a SmudgeTomlMeta with non-nil empty values
func NewSmudgeMeta() SmudgeTomlMeta {
	ss := []string{}
	return SmudgeTomlMeta{Pairings: ss, Themes: ss, Locations: ss}
}

// SmudgeToml defines the top-level structure of the TOML file
type SmudgeToml struct {
	Meta   SmudgeTomlMeta        `toml:"meta"`
	Chars  map[string]string     `toml:"chars"`
	Scenes map[string][][]string `toml:"scenes"`
}

var (
	// init subcommand
	initCmd             = kingpin.Command("init", "initialize").Alias("i")
	initCmdLabel        = initCmd.Arg("label", "Label/shortname of the comic").Required().String()
	initCmdPageCount    = initCmd.Arg("num_pages", "Number of pages in the comic").Required().Int()
	initCmdPageLabel    = initCmd.Flag("page", "Use 'page' as page label (instead of series label)").Short('p').Bool()
	initCmdPageStart    = initCmd.Flag("start", "Starting filename number (does not affect page numbers)").Default("2").Short('s').Int()
	initCmdPageStartAll = initCmd.Flag("start-all", "Starting number for pages and filenames").Default("-1").Int()

	// check subcommand
	checkCmd         = kingpin.Command("check", "check").Alias("c")
	checkCmdFiles    = checkCmd.Arg("file", "TOML file to be checked").Required().Strings()
	checkCmdTagsFile = checkCmd.Flag("tagfile", "TOML file containing list of allowed tags").Short('t').String()

	// quote subcommand
	quoteCmd       = kingpin.Command("quote", "generate quote group templates").Alias("q")
	quoteCmdGroups = quoteCmd.Arg("groups", "quote groups per page").Required().String()
	quoteCmdFlip   = quoteCmd.Flag("flip", "change external separator location to the bottom").Short('f').Bool()
	quoteCmdPrint  = quoteCmd.Flag("print", "print to stdout instead of copying to clipboard").Short('p').Bool()
)

func main() {
	kingpin.Version(version)
	kingpin.CommandLine.HelpFlag.Short('h')

	var err error
	switch kingpin.Parse() {
	case "init":
		pageLabel := *initCmdLabel
		if *initCmdPageLabel == true {
			pageLabel = "page"
		}
		fnStart := *initCmdPageStart
		pageOffset := 1 - fnStart
		if *initCmdPageStartAll != -1 {
			pageOffset = 0
		}
		err = initAction(*initCmdLabel, pageLabel, *initCmdPageCount, fnStart, pageOffset)
	case "check":
		checkAction(*checkCmdFiles, *checkCmdTagsFile)
	case "quote":
		err = quoteAction(*quoteCmdGroups, *quoteCmdFlip, *quoteCmdPrint)
	}

	if err != nil {
		log.Fatalln(err)
	}
}
