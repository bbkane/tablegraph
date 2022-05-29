package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"go.bbkane.com/tablegraph/datatables"
	"go.bbkane.com/warg/command"
)

// -- tableCSV

type tableCSV struct {
	Headers []string
	Rows    [][]string
}

// -- tableDiv

//go:embed embedded/table_div.html
var tableDivTmpl string

type tableDivParams struct {
	DivId     string
	TableJSON string
}

func buildTableDiv(divId string, tableJSON string) (string, error) {
	params := tableDivParams{DivId: divId, TableJSON: tableJSON}
	tmpl, err := template.New("tmpl").Parse(tableDivTmpl)
	if err != nil {
		return "", fmt.Errorf("internal table div template creation error: %w", err)
	}
	sb := strings.Builder{}
	err = tmpl.Execute(&sb, params)
	if err != nil {
		return "", fmt.Errorf("internal table div template execution error: %w", err)
	}
	return sb.String(), nil
}

// -- table command

func table(ctx command.Context) error {
	// parse flags: TODO
	// I/O flags
	format := ctx.Flags["--format"].(string)
	divID, divIDExists := ctx.Flags["--div-id"].(string)

	if !divIDExists {
		divID = randomHexString(10)
	}

	// build csv: TODO

	tCSV := tableCSV{
		Headers: []string{"time", "x", "y"},
		Rows: [][]string{
			{"1", "2", "3"},
			{"3", "4", "5"},
		},
	}

	// build JSON
	dtJSON := datatables.DataTable{
		Data:    tCSV.Rows,
		Columns: make([]datatables.Column, 0),
		ColumnDefs: []datatables.ColumnDef{
			{
				ClassName: "dt-center",
				Targets:   "_all",
			},
		},
		PageLength: 25, // TODO: param
	}

	for _, t := range tCSV.Headers {
		dtJSON.Columns = append(dtJSON.Columns, datatables.Column{
			Title: t,
		})
	}

	switch format {
	case "div":
		// I'm getting the prefix from the html template
		jsonBytes, err := json.MarshalIndent(dtJSON, "        ", "  ")
		if err != nil {
			return fmt.Errorf("error marshalling JSON: %w", err)
		}
		jsonStr := string(jsonBytes)

		div, err := buildTableDiv(divID, jsonStr)
		if err != nil {
			return fmt.Errorf("error formatting div: %w", err)
		}
		fmt.Println(div)
		return nil

	case "html":
		return errors.New("TODO")

	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(dtJSON); err != nil {
			return fmt.Errorf("error encoding JSON: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("invalid --format flag value: %v", format)
	}

}
