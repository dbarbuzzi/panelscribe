package main

import (
	"errors"
	"fmt"
	"text/template"
)

type Page struct {
	Label   string
	PageNum int
	FileNum int
}

type Series struct {
	Label           string
	Pages           []Page
	PlaceholderText string
}

const (
	markdown_template_text string = `# {{ .Label }}
{{ range .Pages }}
## Page {{ .PageNum }} <{{ .Label }}{{ .FileNum }}.jpg>

{{ $.PlaceholderText }}
{{ end }}`

	toml_template_text string = `[meta]
color = false
con = false
pairings = []
themes = []
locations = []

[chars]

[scenes]
{{ range .Pages }}{{ .Label }}{{ .FileNum }} = []
{{ end }}`
)

func getTemplate(key string) (*template.Template, error) {
	templateMap := map[string]string{
		"markdown": markdown_template_text,
		"toml":     toml_template_text,
	}

	if template_text, ok := templateMap[key]; ok {
		return template.Must(template.New("").Parse(template_text)), nil
	}

	return nil, errors.New(fmt.Sprintf("error: key '%s' is not an eligible template", key))
}
