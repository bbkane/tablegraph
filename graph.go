package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
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
	fieldNames, fieldNamesExists := ctx.Flags["--fieldnames"].([]string)
	firstLine, firstLineExists := ctx.Flags["--firstline"].(bool)
	fieldSep := ctx.Flags["--fieldsep"].(string) // TODO: use rune for this...

	// TODO: This doesn't put the field names back on the file!! I think the only way to do this is to always parse the CSV... then extract field names if needed, and add them back on if needed

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

	// Read CSV into file pointer
	var inputFp *os.File
	if inputExists {
		fp, err := os.Open(input)
		if err != nil {
			return fmt.Errorf("Error opening input CSV: %s: %w", input, err)
		}
		defer fp.Close()
		inputFp = fp
	} else {
		inputFp = os.Stdin
		input = "stdin"
	}

	// read file pointer into string
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(inputFp)
	if err != nil {
		return fmt.Errorf("Error reading input CSV: %s: %w", input, err)
	}
	csvContents := buf.String()

	// read first line of CSV and assert len == 3
	r := csv.NewReader(strings.NewReader(csvContents))
	r.Comma = fieldSepRune
	firstLineContents, err := r.Read()
	if err == io.EOF {
		return fmt.Errorf("input CSV appears to be empty: %s: %w", input, err)
	}
	if err != nil {
		return fmt.Errorf("error parsing CSV: %s: %w", input, err)
	}
	if len(firstLineContents) != 3 {
		return fmt.Errorf("the CSV should have 3 columns")
	}

	// Get field names
	if fieldNamesExists && firstLineExists {
		return errors.New("both --fieldnames and --firstline flags passed. Pass only one or don't pass either to generate the fieldnames")
	}

	if fieldNamesExists {
		// do nothing, the user passed fieldnames
		if len(fieldNames) != 3 {
			return fmt.Errorf("--fieldnames should be a list of length 3: %s", fieldNames)
		}
	} else if firstLineExists {
		// grab field names from the first line
		if firstLine == false {
			return errors.New("--firstline should not be false...")
		}
		fieldNames = firstLineContents
	} else {
		// generate columns
		fieldNames = []string{"field_1", "field_2", "field_3"}
	}

	gCSV := graphCSV{
		FieldNames:  fieldNames,
		CSVContents: csvContents,
	}

	fmt.Printf("%#v\n", gCSV)

	// TODO: Turn the above into a separate function

	// print JSON
	return nil
}
