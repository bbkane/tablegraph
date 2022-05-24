package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.bbkane.com/warg/command"

	vl "go.bbkane.com/tablegraph/vegalite"
)

type graphCSV struct {

	// CSVContents contains the CSV plus any header rows
	CSVContents string

	XField     string
	ColorField string
	YField     string
}

func graph(ctx command.Context) error {
	// Parse flags
	input, inputExists := ctx.Flags["--input"].(string)
	fieldNames := ctx.Flags["--fieldnames"].(string) // TODO: this has changed to string
	fieldSep := ctx.Flags["--fieldsep"].(string)     // TODO: use rune for this...

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

	// Read file pointer into CSV
	csvReader := csv.NewReader(inputFp)
	csvReader.Comma = fieldSepRune
	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("unable to parse CSV: %s: %w", input, err)
	}
	if len(records) == 0 {
		return fmt.Errorf("the CSV appears to have no rows: %s", input)
	}

	if len(records[0]) != 3 {
		return fmt.Errorf("the CSV should have 3 columns: %s", input)

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
			return fmt.Errorf("--fieldnames should be a list of length 3: %s", fieldNames)
		}
	}

	// encode CSV into string
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf) // TODO: put this into a string
	if fieldNames != "firstline" {
		if err := w.Write(fieldNamesSlice); err != nil {
			return fmt.Errorf("error writing fieldnames to csv: %s: %w", fieldNames, err)
		}
	}

	w.WriteAll(records)
	if err := w.Error(); err != nil {
		return fmt.Errorf("error writing csv: %s: %w", input, err)
	}

	csvStr := buf.String()

	gCSV := graphCSV{
		CSVContents: csvStr,
		XField:      fieldNamesSlice[0],
		ColorField:  fieldNamesSlice[1],
		YField:      fieldNamesSlice[2],
	}

	// Now, make the fucking JSON

	vlj := vl.JSON{
		Schema:      "https://vega.github.io/schema/vega-lite/v5.json",
		Description: "TODO: Description",
		Data: vl.Data{
			Values: csvStr,
			Format: vl.Format{
				Type: "csv",
			},
		},
		Mark: vl.Mark{
			Type:    "line", // TODO: param
			Tooltip: true,
			Point:   true,
		},
		Height: "container",
		Width:  "container",
		Encoding: vl.Encoding{
			X: vl.XY{
				Field:    gCSV.XField,
				Type:     "temporal",         // TODO: param
				TimeUnit: "utcyearmonthdate", // TODO: param
				Scale: &vl.Scale{
					Type: "utc",
				},
			},
			Y: vl.XY{
				Field: gCSV.YField,
				Type:  "quantitative", // TODO: param
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
			Text: "TODO: title", // TODO: param
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

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(vlj); err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	return nil
}
