package main

import (
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"go.bbkane.com/tablegraph/datatables"
	"go.bbkane.com/warg/command"
)

// -- tableCSV

type tableCSV struct {
	ColumnNames []string
	Rows        [][]string
}

func newTableCSV(r io.Reader, fieldNames string, fieldSep rune) (*tableCSV, error) {
	// Read file pointer into CSV
	csvReader := csv.NewReader(r)
	csvReader.Comma = fieldSep
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to parse CSV: %w", err)
	}

	var tCSV tableCSV

	// if firstline passed, use it
	if fieldNames == "firstline" {
		tCSV.ColumnNames = records[0]
		tCSV.Rows = records[1:]
	} else if fieldNames == "" {
		// no fieldnames passed, let's generate them now that we can know how many columns are in the CSV
		numCols := len(records[0])
		for i := 0; i < numCols; i++ {
			tCSV.ColumnNames = append(tCSV.ColumnNames, fmt.Sprintf("col_%d", i))
		}
		tCSV.Rows = records
	} else {
		// othewise, use passed fieldnames, add them to csv
		tCSV.ColumnNames = strings.Split(fieldNames, ",")
		tCSV.Rows = records
	}
	return &tCSV, nil
}

// -- tableDiv

//go:embed embedded/tableDiv.html
var tableDivTmpl string

type tableDivParams struct {
	DivId     string
	TableJSON string
}

// buildTemplateString reads a template from a string, fills in the params,
// and returns the resulting string
func buildTemplateString(tmpl string, params interface{}) (string, error) {

	parsed, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("internal template creation error: %w", err)
	}
	sb := strings.Builder{}
	err = parsed.Execute(&sb, params)
	if err != nil {
		return "", fmt.Errorf("internal template execution error: %w", err)
	}
	return sb.String(), nil
}

func buildTableDiv(divId string, dtJSON datatables.DataTable) (string, error) {

	// I'm getting the prefix from the html template
	jsonBytes, err := json.MarshalIndent(dtJSON, "        ", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %w", err)
	}
	jsonStr := string(jsonBytes)

	params := tableDivParams{DivId: divId, TableJSON: jsonStr}

	return buildTemplateString(tableDivTmpl, params)
}

// -- table command

func table(ctx command.Context) error {
	// I/O flags
	format := ctx.Flags["--format"].(string)
	input, inputExists := ctx.Flags["--input"].(string)
	divID, divIDExists := ctx.Flags["--div-id"].(string)

	// CSV flags
	// for the table, we want an empty string so we can generate the fieldnames
	fieldNames, _ := ctx.Flags["--fieldnames"].(string)
	fieldSep := ctx.Flags["--fieldsep"].(string)

	htmlTitle := ctx.Flags["--html-title"].(string)
	pageLength := ctx.Flags["--page-length"].(int)

	if !divIDExists {
		divID = randomHexString(10)
	}

	// Get a fieldSepRune
	if fieldSep == "" {
		return errors.New("--fieldsep should not be an empty string")
	}
	var fieldSepRune rune
	if fieldSepRunes := []rune(fieldSep); len(fieldSepRunes) != 1 {
		return fmt.Errorf("--fieldsep should only be one character")
	} else {
		fieldSepRune = fieldSepRunes[0]
	}

	// Get file pointer
	var inputFp *os.File
	if inputExists {
		fp, err := os.Open(input)
		if err != nil {
			return fmt.Errorf("error opening input CSV: %s: %w", input, err)
		}
		defer fp.Close()
		inputFp = fp
	} else {
		inputFp = os.Stdin
		input = "stdin"
	}

	// build CSV
	tCSV, err := newTableCSV(inputFp, fieldNames, fieldSepRune)
	if err != nil {
		return fmt.Errorf("Error building tableCSV: %s: %w", input, err)
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
		PageLength: pageLength,
	}

	for _, t := range tCSV.ColumnNames {
		dtJSON.Columns = append(dtJSON.Columns, datatables.Column{
			Title: t,
		})
	}

	switch format {
	case "div":
		div, err := buildTableDiv(divID, dtJSON)
		if err != nil {
			return fmt.Errorf("error formatting div: %w", err)
		}
		fmt.Println(div)
		return nil

	case "html":
		div, err := buildTableDiv(divID, dtJSON)
		if err != nil {
			return fmt.Errorf("error formatting div: %w", err)
		}
		htmlTop, err := buildHtmlTop(htmlTitle)
		if err != nil {
			return fmt.Errorf("error formatting html top: %w", err)
		}

		fmt.Println(htmlTop)
		fmt.Println(div)
		fmt.Println(buildHtmlBottom())
		return nil

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
