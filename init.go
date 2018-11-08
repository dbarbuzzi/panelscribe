package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

func createBaseToml(numPages int) SmudgeToml {
	var res SmudgeToml

	res.Tags = []string{}
	res.Chars = make(map[string]string)
	res.Scenes = make(map[string][][]string)

	pagekeyTemplate := "p%d"
	if numPages >= 10 {
		pagekeyTemplate = "p%02d"
	}
	for i := 1; i <= numPages; i++ {
		pagekey := fmt.Sprintf(pagekeyTemplate, i)
		res.Scenes[pagekey] = [][]string{}
	}

	return res
}

func writeInitialToml(basename string, data SmudgeToml) error {
	f, err := os.Create(fmt.Sprintf("%s.toml", basename))
	if err != nil {
		return err
	}
	defer f.Close()

	e := toml.NewEncoder(f)
	err = e.Encode(data)
	return err
}

func createBaseMd(title string, numPages int) string {
	lines := []string{fmt.Sprintf("# %s\n", title)}

	for i := 1; i <= numPages; i++ {
		line := fmt.Sprintf("## Page %d <page%d.jpg>\n", i, i)
		lines = append(lines, []string{"", "upper placeholder", "", "---", "", "lower placeholder", ""}...)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func writeInitialMd(basename, data string) error {
	f, err := os.Create(fmt.Sprintf("%s.md", basename))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	return err
}

// initAction initializes meta (TOML) and transcript (Markdown) files
func initAction(basename string, numPages int) error {
	// Create the initial TOML file
	tomlData := createBaseToml(numPages)
	err := writeInitialToml(basename, tomlData)
	if err != nil {
		return err
	}

	// Create the initial Markdown file
	md := createBaseMd(basename, numPages)
	err = writeInitialMd(basename, md)
	if err != nil {
		return err
	}

	return nil
}
