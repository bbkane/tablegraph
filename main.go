package main

import (
	"os"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value"
)

func main() {

	htmlTitleFlag := flag.New(
		"HTML title. Flag ignored when --format != 'html'",
		value.String,
		flag.Default("tablegraph output"),
	)

	csvParseFlags := flag.FlagMap{
		"--fieldnames": flag.New(
			"Pass list of field names",
			value.StringSlice,
		),
		"--fieldsep": flag.New(
			flag.HelpShort("Field separator for input table. TAB by default"),
			value.String,
			flag.Default(","), // changed from TAB
			flag.Required(),
		),
		"--firstline": flag.New(
			flag.HelpShort("Use the first line of the table as fieldnames"), // TODO: make exclusive to --fieldnames
			value.Bool,
		),
	}

	// flags for making both charts and tables!
	ioFlags := flag.FlagMap{
		"--div-id": flag.New(
			flag.HelpShort("ID of div when --format is 'div'"),
			value.String,
		),
		"--format": flag.New(
			flag.HelpShort("Output format"),
			value.StringEnum("div", "html", "json"),
			flag.Default("html"),
			flag.Required(),
		),
		"--html-title": htmlTitleFlag,
		"--input": flag.New(
			flag.HelpShort("Input file"),
			value.Path,
		),
	}

	graphFlags := flag.FlagMap{
		"--chart-title": flag.New(
			flag.HelpShort("Chart title"),
			value.String,
		),
		"--x-axis-title": flag.New(
			flag.HelpShort("X-Axis Title"),
			value.String,
		),
		"--x-type": flag.New(
			"X type. See https://vega.github.io/vega-lite/docs/type.html",
			value.StringEnum("nominal", "quantitative", "temporal"),
		),
		"--y-type": flag.New(
			"Y type. See https://vega.github.io/vega-lite/docs/type.html",
			value.StringEnum("nominal", "quantitative", "temporal"),
		),
		"--y-axis-title": flag.New(
			flag.HelpShort("Y-Axis Title"),
			value.String,
		),
	}

	app := warg.New(
		"tablegraph",
		section.New(
			section.HelpShort("Turn CSVs into graphs! NOTE: this is an experiment at this stage. Use chart.py"),
			section.Flag(
				"--output",
				flag.HelpShort("Path to output file. Use DATEME as an alias for 'chart.<timestamp>'."),
				value.Path,
				flag.Default("DATEME.html"),
				flag.Required(),
			),
			section.Command(
				"3-col",
				// "Graph from a 3-column CSV. First column is x-axis, and should be datetime. Second column is a string whose values are used to 'group'. Third column is a numeric column for the y-axis",
				"Various graphs from 3-column CSVs (x,category,y). Point, line, grouped-bar, stacked-bar ",
				command.DoNothing,
				command.ExistingFlags(csvParseFlags),
				command.ExistingFlags(ioFlags),
				command.ExistingFlags(graphFlags),
				command.Flag(
					"--type",
					"Type of graph to generate",
					value.StringEnum("point", "line", "grouped-bar", "stacked-bar"),
				),
				command.HelpLong("First column is x-axis, and should be datetime. Second column is a string whose values are used to 'group'. Third column is a numeric column for the y-axis"),
			),
			section.Command(
				"table",
				command.HelpShort("Make HTML table"),
				command.DoNothing,
				command.ExistingFlags(csvParseFlags),
				command.ExistingFlags(ioFlags),
				command.Flag(
					"--page-length",
					flag.HelpShort("Entries in table before needing to click next page"),
					value.Int,
				),
				command.HelpLong("NOTE: columns should not have a `.` in the title. See https://datatables.net/forums/discussion/69257/data-with-a-in-the-name-makes-table-creation-fail#latest"),
			),
			section.Section(
				"html",
				"HTML snippets",
				section.Command(
					"top",
					"HTML top",
					command.DoNothing,
					command.ExistingFlag("--html-title", htmlTitleFlag),
				),
				section.Command(
					"bottom",
					"HTML bottom",
					command.DoNothing,
				),
			),
		),
	)
	app.MustRun(os.Args, os.LookupEnv)
}
