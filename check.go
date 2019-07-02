package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

func contains(haystack []string, needle string) bool {
	for _, hay := range haystack {
		if hay == needle {
			return true
		}
	}
	return false
}

func hasDupes(haystack []string) bool {
	for i, needle := range haystack {
		for _, hay := range haystack[i+1:] {
			if needle == hay {
				return true
			}
		}
	}
	return false
}

func loadSmudgeToml(tomlPath string) (SmudgeToml, error) {
	// TODO: Decode to a generic map so we can validate existence of keys separately from content
	res := SmudgeToml{}
	_, err := toml.DecodeFile(tomlPath, &res)
	return res, err
}

func loadAllowedTagsToml(tomlPath string) ([]string, error) {
	res := struct{ Tags []string }{}
	_, err := toml.DecodeFile(tomlPath, &res)
	return res.Tags, err
}

func validateMetaBlock(data SmudgeToml, allowedTags []string) []string {
	results := []string{}

	// validate misc tag lists
	if res := validateTagList(data.Meta.Pairings, allowedTags, "pairings"); len(res) > 0 {
		results = append(results, res...)
	}
	if res := validateTagList(data.Meta.Themes, allowedTags, "themes"); len(res) > 0 {
		results = append(results, res...)
	}
	if res := validateTagList(data.Meta.Locations, allowedTags, "locations"); len(res) > 0 {
		results = append(results, res...)
	}

	return results
}

func validateTagList(tagList, allowed []string, name string) []string {
	results := []string{}

	// check that it’s not empty
	if len(tagList) == 0 {
		results = append(results, fmt.Sprintf("error: %s tag list is empty", name))
	}

	// check if they’re sorted
	if !sort.StringsAreSorted(tagList) {
		results = append(results, fmt.Sprintf("warn: %s tag list is not in alphabetical order", name))
		sort.Strings(tagList)
		sorted := strings.Join(tagList, "\", \"")
		results = append(results, fmt.Sprintf("info: %s tag list sorted - \"%s\"", name, sorted))
	}

	// check for duplicates
	if hasDupes(tagList) {
		results = append(results, fmt.Sprintf("error: %s tag list contains duplicate entries", name))
	}

	// check that it only includes tags from master tag list
	if len(allowed) > 0 {
		for _, tag := range tagList {
			if !contains(allowed, tag) {
				results = append(results, fmt.Sprintf("error: %s tag “%s” not found in allowed tag list", name, tag))
			}
		}
	}

	return results
}

func validateCharList(data SmudgeToml) []string {
	results := []string{}

	// check that it’s not empty
	if len(data.Chars) == 0 {
		results = append(results, "error: [chars] char list is empty")
	}

	// check that all chars have a description
	for char, desc := range data.Chars {
		if len(desc) == 0 {
			results = append(results, fmt.Sprintf("error: [chars] description for char “%s” is empty", char))
		}
	}

	return results
}

func validateSceneList(data SmudgeToml) []string {
	results := []string{}

	// check that it’s not empty
	if len(data.Scenes) == 0 {
		results = append(results, "error: [scenes] scene list is empty")
	} else {
		for pagekey, page := range data.Scenes {
			if len(page) == 0 {
				results = append(results, fmt.Sprintf("error: [scenes] page “%s” is empty", pagekey))
			} else {
				for cellnum, cell := range page {
					if len(cell) == 0 {
						results = append(results, fmt.Sprintf("error: [scenes] cell %d of “%s” is empty", cellnum+1, pagekey))
					} else {
						for descnum, desc := range cell {
							if len(desc) == 0 {
								results = append(results, fmt.Sprintf("error: [scenes] desc %d of page “%s”/cell %d is empty", descnum+1, pagekey, cellnum+1))
							}
						}
					}
				}
			}
		}
	}

	return results
}

func validateSmudgeToml(data SmudgeToml, allowedTags []string) []string {
	results := []string{}

	results = append(results, validateMetaBlock(data, allowedTags)...)
	results = append(results, validateCharList(data)...)
	results = append(results, validateSceneList(data)...)

	return results
}

func validateFile(filepath string, allowedTags []string) []string {
	data, err := loadSmudgeToml(filepath)
	if err != nil {
		panic(fmt.Sprintf("error: unable to load TOML file: %s", err))
	}
	return validateSmudgeToml(data, allowedTags)
}

func checkAction(tomlFiles []string, allowedTagsFile string) {
	allowedTags := []string{}
	var err error
	if len(allowedTagsFile) > 0 {
		allowedTags, err = loadAllowedTagsToml(allowedTagsFile)
		if err != nil {
			fmt.Printf("error: unable to load allowed tags file: %s\n", err)
		}
	}
	if len(allowedTags) == 0 {
		fmt.Println("warn: [tags] allowed tag list empty; skipping tag content validation")
	}

	hasError := false
	for _, tomlFile := range tomlFiles {
		// dumb skip of AllowedTagsFile
		if tomlFile == allowedTagsFile {
			continue
		}
		fmt.Printf("checking %s...\n", tomlFile)
		results := validateFile(tomlFile, allowedTags)
		for _, result := range results {
			fmt.Println(result)
			if strings.HasPrefix(result, "error:") {
				hasError = true
			}
		}
	}

	if hasError {
		os.Exit(1)
	}
}
