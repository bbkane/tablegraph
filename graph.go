package main

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	// NOTE: we're generating JSON and we don't want it escaped, so just use text/template, NOT html/template
	"text/template"

	"go.bbkane.com/warg/command"

	vl "go.bbkane.com/tablegraph/vegalite"
)

// -- graphCSV

type graphCSV struct {

	// CSVContents contains the CSV plus any header rows
	CSVContents string

	XField     string
	ColorField string
	YField     string
}

func newGraphCSV(r io.Reader, readerName string, fieldNames string, fieldSep rune) (*graphCSV, error) {
	// Read file pointer into CSV
	csvReader := csv.NewReader(r)
	csvReader.Comma = fieldSep
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to parse CSV: %s: %w", readerName, err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("the CSV appears to have no rows: %s", readerName)
	}

	if len(records[0]) != 3 {
		return nil, fmt.Errorf("the CSV should have 3 columns: %s", readerName)

	}

	// get fieldnames
	var fieldNamesSlice []string
	// if firstline passed, use it
	if fieldNames == "firstline" {
		fieldNamesSlice = records[0]
	} else {
		// othewise, use passed fieldnames, add them to csv
		fieldNamesSlice = strings.Split(fieldNames, ",")
		if len(fieldNamesSlice) != 3 {
			return nil, fmt.Errorf("--fieldnames should be a list of length 3: %s", fieldNames)
		}
	}

	// encode CSV into string
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf) // TODO: put this into a string
	if fieldNames != "firstline" {
		if err := w.Write(fieldNamesSlice); err != nil {
			return nil, fmt.Errorf("error writing fieldnames to csv: %s: %w", fieldNames, err)
		}
	}

	w.WriteAll(records)
	if err := w.Error(); err != nil {
		return nil, fmt.Errorf("error saving CSV to string: %w", err)
	}

	csvStr := buf.String()

	gCSV := graphCSV{
		CSVContents: csvStr,
		XField:      fieldNamesSlice[0],
		ColorField:  fieldNamesSlice[1],
		YField:      fieldNamesSlice[2],
	}
	return &gCSV, nil
}

// -- graphDiv

//go:embed embedded/graphDiv.html
var graphDivTmpl string

type graphDivParams struct {
	DivId        string
	VegaLiteJSON string
}

func buildGraphDiv(divId string, vegaLiteJSON string) (string, error) {
	params := graphDivParams{DivId: divId, VegaLiteJSON: vegaLiteJSON}
	tmpl, err := template.New("tmpl").Parse(graphDivTmpl)
	if err != nil {
		return "", fmt.Errorf("internal graph div template creation error: %w", err)
	}
	sb := strings.Builder{}
	err = tmpl.Execute(&sb, params)
	if err != nil {
		return "", fmt.Errorf("internal graph div template execution error: %w", err)
	}
	return sb.String(), nil
}

// -- graph command

func randomHexString(length int) string {
	// https://stackoverflow.com/a/65607935/2958070
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("div_%x", b)[:length+4]
}

func graph(ctx command.Context) error {
	// I/O flags
	format := ctx.Flags["--format"].(string)
	input, inputExists := ctx.Flags["--input"].(string)
	divID, divIDExists := ctx.Flags["--div-id"].(string)

	// CSV flags
	fieldNames := ctx.Flags["--fieldnames"].(string)
	fieldSep := ctx.Flags["--fieldsep"].(string)

	// Graph Flags
	gTitle := ctx.Flags["--graph-title"].(string)
	gType := ctx.Flags["--type"].(string)
	gXType, _ := ctx.Flags["--x-type"].(string)
	gYType, _ := ctx.Flags["--y-type"].(string)
	gXTimeUnit, _ := ctx.Flags["--x-time-unit"].(string)

	if !divIDExists {
		divID = randomHexString(10)
	}

	switch gType {
	case "grouped-bar":
		gType = "bar"
		// TODO: add the xOffset thingie
	case "stacked-bar":
		gType = "bar"
	default:
		// pass
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

	gCSV, err := newGraphCSV(inputFp, input, fieldNames, fieldSepRune)
	if err != nil {
		return err
	}

	vlj := vl.JSON{
		Schema:      "https://vega.github.io/schema/vega-lite/v5.json",
		Description: "TODO: Description",
		Data: vl.Data{
			Values: gCSV.CSVContents,
			Format: vl.Format{
				Type: "csv",
			},
		},
		Mark: vl.Mark{
			Type:    gType,
			Tooltip: true,
			Point:   true,
		},
		Height: "container",
		Width:  "container",
		Encoding: vl.Encoding{
			X: vl.XY{
				Field:    gCSV.XField,
				Type:     gXType,
				TimeUnit: gXTimeUnit,
				Scale: &vl.Scale{
					Type: "utc",
				}, // TODO: param
			},
			Y: vl.XY{
				Field: gCSV.YField,
				Type:  gYType,
			},
			Color: vl.Color{
				Field: gCSV.ColorField,
				Type:  "nominal",
			},
			Opacity: vl.Opacity{
				Condition: vl.Condition{
					Param: "hover",
					Value: 1,
				},
				Value: 0.1,
			},
		},
		Title: vl.Title{
			Text: gTitle,
		},
		Params: []vl.Params{
			{
				Name: "hover",
				Bind: "legend",
				Select: vl.Select{
					Type:   "point",
					Fields: []string{"symbol"},
				},
			},
		},
	}

	switch format {
	case "div":
		// I'm getting the prefix from the html template
		jsonBytes, err := json.MarshalIndent(vlj, "        ", "  ")
		if err != nil {
			return fmt.Errorf("error marshalling JSON: %w", err)
		}
		jsonStr := string(jsonBytes)

		graphDiv, err := buildGraphDiv(divID, jsonStr)
		if err != nil {
			return fmt.Errorf("error formatting div: %w", err)
		}
		fmt.Println(graphDiv)
		return nil

	case "html":
		return errors.New("TODO")

	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(vlj); err != nil {
			return fmt.Errorf("error encoding JSON: %w", err)
		}
		return nil

	default:
		return fmt.Errorf("invalid --format flag value: %v", format)
	}

}
