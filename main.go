package main

import (
	"log"

	"gopkg.in/alecthomas/kingpin.v2"
)

var version = "dev"

// ComicMetadataMeta defines the "meta" block in the comic metadata file
type ComicMetadataMeta struct {
	Color     bool     `toml:"color"`
	Con       bool     `toml:"con"`
	Pairings  []string `toml:"pairings"`
	Themes    []string `toml:"themes"`
	Locations []string `toml:"locations"`
}

// NewComicMetadataMeta creates a ComicMetadataMeta with non-nil empty values
func NewComicMetadataMeta() ComicMetadataMeta {
	ss := []string{}
	return ComicMetadataMeta{Pairings: ss, Themes: ss, Locations: ss}
}

// ComicMetadata defines the top-level structure of the metadata file
type ComicMetadata struct {
	Meta   ComicMetadataMeta     `toml:"meta"`
	Chars  map[string]string     `toml:"chars"`
	Scenes map[string][][]string `toml:"scenes"`
}

var (
	// init subcommand
	initCmd             = kingpin.Command("init", "initialize Markdown and TOML files for a new series").Alias("i")
	initCmdLabel        = initCmd.Arg("label", "Label/shortname of the comic").Required().String()
	initCmdPageCount    = initCmd.Arg("num_pages", "Number of pages in the comic").Required().Int()
	initCmdPageLabel    = initCmd.Flag("page", "Use 'page' as page label (instead of series label)").Short('p').Bool()
	initCmdPageStart    = initCmd.Flag("start", "Starting filename number (does not affect page numbers)").Default("2").Short('s').Int()
	initCmdPageStartAll = initCmd.Flag("start-all", "Starting number for pages and filenames").Default("-1").Int()
	initCmdTextMode     = initCmd.Flag("text", "Enable text mode, where the page sections include top and bottom text placeholders").Short('t').Bool()

	// check subcommand
	checkCmd         = kingpin.Command("check", "check a TOML metadata file for correctness").Alias("c")
	checkCmdFiles    = checkCmd.Arg("file", "TOML file to be checked").Required().Strings()
	checkCmdTagsFile = checkCmd.Flag("tagfile", "TOML file containing list of allowed tags").Short('t').String()

	// quote subcommand
	quoteCmd       = kingpin.Command("quote", "generate quote group templates to enter into Markdown file").Alias("q")
	quoteCmdGroups = quoteCmd.Arg("groups", "quote groups per page").Required().String()
	quoteCmdBorder = quoteCmd.Flag("border", "which border(s) to put separator(s) on").Short('b').Default("top").Enum("top", "bottom", "both", "none")
	quoteCmdPrint  = quoteCmd.Flag("print", "print to stdout instead of copying to clipboard").Short('p').Bool()
)

func main() {
	kingpin.Version(version)
	kingpin.CommandLine.HelpFlag.Short('h')

	var err error
	switch kingpin.Parse() {
	case "init":
		pageLabel := *initCmdLabel
		if *initCmdPageLabel {
			pageLabel = "page"
		}
		fnStart := *initCmdPageStart
		pageOffset := 1 - fnStart
		if *initCmdPageStartAll != -1 {
			pageOffset = 0
		}
		err = initAction(*initCmdLabel, pageLabel, *initCmdPageCount, fnStart, pageOffset, *initCmdTextMode)
	case "check":
		checkAction(*checkCmdFiles, *checkCmdTagsFile)
	case "quote":
		err = quoteAction(*quoteCmdGroups, *quoteCmdBorder, *quoteCmdPrint)
	}

	if err != nil {
		log.Fatalln(err)
	}
}
