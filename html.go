package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"strings"

	"go.bbkane.com/warg/command"
)

//go:embed embedded/htmlTop.html
var htmlTop string

//go:embed embedded/htmlBottom.html
var htmlBottom string

type htmlTopParams struct {
	HtmlTitle string
}

func buildHtmlTop(htmlTitle string) (string, error) {
	htmlInfo := htmlTopParams{HtmlTitle: htmlTitle}
	tmpl, err := template.New("tmpl").Parse(htmlTop)
	if err != nil {
		return "", fmt.Errorf("internal html template creation error: %w", err)
	}
	sb := strings.Builder{}
	err = tmpl.Execute(&sb, htmlInfo)
	if err != nil {
		return "", fmt.Errorf("internal html template execution error: %w", err)
	}
	return sb.String(), nil
}

func buildHtmlBottom() string {
	return htmlBottom
}

func runHtmlTop(ctx command.Context) error {
	htmlTitle := ctx.Flags["--html-title"].(string)

	htmlTop, err := buildHtmlTop(htmlTitle)
	if err != nil {
		return err
	}
	fmt.Println(htmlTop)
	return nil
}

func runHtmlBottom(_ command.Context) error {
	fmt.Println(htmlBottom)
	return nil
}
