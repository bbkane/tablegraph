package main

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	// NOTE: we're generating JSON and we don't want it escaped, so just use text/template, NOT html/template

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

	err = w.WriteAll(records)
	if err != nil {
		return nil, fmt.Errorf("error writing row: %w", err)
	}
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
	DivWidth     string
	DivHeight    string
}

func buildGraphDiv2(divId string, vgl vl.JSON, divHeight string, divWidth string) (string, error) {
	// I'm getting the prefix from the html template
	jsonBytes, err := json.MarshalIndent(vgl, "        ", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %w", err)
	}
	jsonStr := string(jsonBytes)

	params := graphDivParams{DivId: divId, VegaLiteJSON: jsonStr, DivWidth: divWidth, DivHeight: divHeight}
	return buildTemplateString(graphDivTmpl, params)

}

// -- graph command

func randomHexString(length int) string {
	// https://stackoverflow.com/a/54491783/2958070
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Sprintf("error generating random number: %s", err))
	}
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

	// HTML flags
	htmlTitle := ctx.Flags["--html-title"].(string)
	divWidth := ctx.Flags["--div-width"].(string)
	divHeight := ctx.Flags["--div-height"].(string)

	// Graph Flags
	gMarkSize, _ := ctx.Flags["--mark-size"].(int)
	gTitle := ctx.Flags["--graph-title"].(string)
	gType := ctx.Flags["--type"].(string)
	gXScaleType, gXScaleTypeExists := ctx.Flags["--x-scale-type"].(string)
	gXType, _ := ctx.Flags["--x-type"].(string)
	gYScaleType, gYScaleTypeExists := ctx.Flags["--y-scale-type"].(string)
	gYType, _ := ctx.Flags["--y-type"].(string)
	gXTimeUnit, _ := ctx.Flags["--x-time-unit"].(string)

	if !divIDExists {
		divID = randomHexString(10)
	}

	var xScale *vl.Scale
	if gXScaleTypeExists {
		xScale = &vl.Scale{Type: gXScaleType}
	}
	var yScale *vl.Scale
	if gYScaleTypeExists {
		yScale = &vl.Scale{Type: gYScaleType}
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

	var xOffset *vl.XOffset
	switch gType {
	case "grouped-bar":
		gType = "bar"
		xOffset = &vl.XOffset{Field: gCSV.ColorField}
	case "stacked-bar":
		gType = "bar"
	default:
		// pass
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
			Size:    gMarkSize,
		},
		Height: "container",
		Width:  "container",
		Encoding: vl.Encoding{
			X: vl.XY{
				Field:    gCSV.XField,
				Type:     gXType,
				TimeUnit: gXTimeUnit,
				Scale:    xScale,
			},
			Y: vl.XY{
				Field:    gCSV.YField,
				Type:     gYType,
				Scale:    yScale,
				TimeUnit: "",
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
			XOffset: xOffset,
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
					Fields: []string{gCSV.ColorField},
				},
			},
		},
	}

	switch format {
	case "div":
		div, err := buildGraphDiv2(divID, vlj, divHeight, divWidth)
		if err != nil {
			return fmt.Errorf("error formatting div: %w", err)
		}
		fmt.Println(div)

		return nil

	case "html":
		div, err := buildGraphDiv2(divID, vlj, divHeight, divWidth)
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
		if err := enc.Encode(vlj); err != nil {
			return fmt.Errorf("error encoding JSON: %w", err)
		}
		return nil

	default:
		return fmt.Errorf("invalid --format flag value: %v", format)
	}

}
