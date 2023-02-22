package main

import (
	"fmt"
	"os"
)

// TODO: make all errors explicit/wrapped via fmt.Errorf

func generateToml(series Series) error {
	f, err := os.Create(fmt.Sprintf("%s.toml", series.Label))
	if err != nil {
		return err
	}
	defer f.Close()

	tpl, err := getTemplate("toml")
	if err != nil {
		return err
	}

	return tpl.Execute(f, series)
}

func generateMarkdown(series Series) error {
	f, err := os.Create(fmt.Sprintf("%s.md", series.Label))
	if err != nil {
		return err
	}
	defer f.Close()

	tpl, err := getTemplate("markdown")
	if err != nil {
		return err
	}

	return tpl.Execute(f, series)
}

func initAction(label, pageLabel string, pageCount, fnStart, pageOffset int) error {
	pages := make([]Page, 0, pageCount)
	end := fnStart + pageCount
	for i := fnStart; i < end; i++ {
		page := Page{pageLabel, i + pageOffset, i}
		pages = append(pages, page)
	}
	series := Series{label, pages}

	if err := generateMarkdown(series); err != nil {
		return fmt.Errorf("error generating Markdown file: %w", err)
	}
	if err := generateToml(series); err != nil {
		return fmt.Errorf("error generating TOML file: %w", err)
	}

	return nil
}
