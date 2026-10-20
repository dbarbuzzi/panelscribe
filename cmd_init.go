package main

import (
	"fmt"
	"os"
)

func generateToml(series Series) error {
	tpl, err := getTemplate("toml")
	if err != nil {
		return fmt.Errorf("failed to get 'toml' template: %w", err)
	}

	filename := fmt.Sprintf("%s.toml", series.Label)
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create toml file '%s': %w", filename, err)
	}
	defer f.Close()

	return tpl.Execute(f, series)
}

func generateMarkdown(series Series) error {
	tpl, err := getTemplate("markdown")
	if err != nil {
		return fmt.Errorf("failed to get 'markdown' template: %w", err)
	}

	filename := fmt.Sprintf("%s.md", series.Label)
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create markdown file '%s': %w", filename, err)
	}
	defer f.Close()

	return tpl.Execute(f, series)
}

func initAction(label, pageLabel string, pageCount, fnStart, pageOffset int, textMode bool) error {
	pages := make([]Page, 0, pageCount)
	end := fnStart + pageCount
	for i := fnStart; i < end; i++ {
		page := Page{pageLabel, i + pageOffset, i}
		pages = append(pages, page)
	}
	placeholderText := "placeholder_text"
	if textMode {
		placeholderText = "top_placeholder\n\n---\n\nbottom_placeholder"
	}
	series := Series{label, pages, placeholderText}

	if err := generateMarkdown(series); err != nil {
		return fmt.Errorf("error generating Markdown file: %w", err)
	}
	if err := generateToml(series); err != nil {
		return fmt.Errorf("error generating TOML file: %w", err)
	}

	return nil
}
