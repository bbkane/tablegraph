package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.bbkane.com/warg/command"
)

type graphCSV struct {
	FieldNames []string

	// CSVContents contains the CSV plus any header rows
	CSVContents string
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
		FieldNames:  fieldNamesSlice,
		CSVContents: csvStr,
	}

	fmt.Println(gCSV)
	return nil
}
