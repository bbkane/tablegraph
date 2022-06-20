package main

import (
	"os"

	_ "embed"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value"
)

//go:embed embedded/graphFooter.txt
var graphFooter string

func main() {

	htmlTitleFlag := flag.New(
		"HTML title. Flag ignored when --format != 'html'",
		value.String,
		flag.Default("tablegraph output"),
	)

	csvParseFlags := flag.FlagMap{
		"--fieldsep": flag.New(
			flag.HelpShort("Field separator for input table"),
			value.String,
			flag.Default(","), // changed from TAB
			flag.Required(),
		),
	}

	// flags for making both graphs and tables!
	ioFlags := flag.FlagMap{
		"--div-id": flag.New(
			flag.HelpShort("ID of div when --format is 'div'. If not passed, a random string will be used"),
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
			flag.HelpShort("Input file. tablegraph will use stdin if not passed"),
			value.Path,
		),
	}

	graphFlags := flag.FlagMap{
		"--graph-title": flag.New(
			flag.HelpShort("Graph title"),
			value.String,
			flag.Default("Graph Title"),
		),
		"--mark-size": flag.New(
			"Size of mark (useful for making bar charts wider",
			value.Int,
		),
		"--type": flag.New(
			"Type of graph to generate",
			value.StringEnum("point", "line", "grouped-bar", "stacked-bar"),
			flag.Default("line"),
			flag.Required(),
		),
		"--x-time-unit": flag.New(
			"X Time Unit only used for temporal types - see https://vega.github.io/vega-lite/docs/timeunit.html . It's advised to use ones prefixed with utc - utcyear , utcyearmonthday",
			value.String,
		),
		"--x-scale-type": flag.New(
			"X scale type. See https://vega.github.io/vega-lite/docs/scale.html . Particularly useful to set to 'utc' for time-valued charts",
			value.String,
		),
		"--x-type": flag.New(
			"X type. See https://vega.github.io/vega-lite/docs/type.html",
			value.StringEnum("nominal", "quantitative", "temporal"),
		),
		"--y-scale-type": flag.New(
			"Y scale type. See https://vega.github.io/vega-lite/docs/scale.html . Particularly useful to set to 'utc' for time-valued charts",
			value.String,
		),
		"--y-type": flag.New(
			"Y type. See https://vega.github.io/vega-lite/docs/type.html",
			value.StringEnum("nominal", "quantitative", "temporal"),
		),
	}

	app := warg.New(
		"tablegraph",
		section.New(
			section.HelpShort("Turn CSVs into graphs! NOTE: this is an experiment at this stage. Use tablegraph.py"),
			section.Section(
				"examples",
				"Print some more complicated real-world examples",
				section.Command(
					"repolines",
					"Stacked bar chart of lines changed over time in a git repo",
					printString(examplesLinesText),
				),
				// section.Command(
				// 	"langlines",
				// 	"Grouped bar chart of lines of code by language in a folder",
				// 	printString(examplesLangLinesText),
				// ),
				// Add something from osquery, something from github stars? Shell history chart?
			),
			section.Command(
				"graph",
				// "Graph from a 3-column CSV. First column is x-axis, and should be datetime. Second column is a string whose values are used to 'group'. Third column is a numeric column for the y-axis",
				"Various graphs from 3-column CSVs (x,category,y). Point, line, grouped-bar, stacked-bar ",
				graph,
				command.ExistingFlags(csvParseFlags),
				command.ExistingFlags(ioFlags),
				command.ExistingFlags(graphFlags),
				command.Flag(
					"--fieldnames",
					"Pass comma separated fieldnames (ex: 'date,type,lines') or the string 'firstline' to use the first line of the CSV",
					value.String,
					flag.Default("x,category,y"),
				),
				command.HelpLong("First column is x-axis, and should be datetime. Second column is a string whose values are used to 'group'. Third column is a numeric column for the y-axis"),
				command.Footer(
					graphFooter,
				),
			),
			section.Command(
				"table",
				command.HelpShort("Make HTML table"),
				table,
				command.ExistingFlags(csvParseFlags),
				command.ExistingFlags(ioFlags),
				command.Flag(
					"--fieldnames",
					"Pass comma separated fieldnames (ex: 'date,type,lines') or the string 'firstline' to use the first line of the CSV. Will be generated if not passed.",
					value.String,
				),
				command.Flag(
					"--page-length",
					"Entries in table before needing to click next page",
					value.Int,
					flag.Default("30"),
				),
				command.HelpLong("NOTE: columns should not have a `.` in the title. See https://datatables.net/forums/discussion/69257/data-with-a-in-the-name-makes-table-creation-fail#latest"),
			),
			section.Section(
				"html",
				"HTML snippets",
				section.Command(
					"top",
					"HTML top",
					runHtmlTop,
					command.ExistingFlag("--html-title", htmlTitleFlag),
				),
				section.Command(
					"bottom",
					"HTML bottom",
					runHtmlBottom,
				),
			),
			section.Command(
				"version",
				"Print version",
				printVersion,
			),
		),
	)
	app.MustRun(os.Args, os.LookupEnv)
}
