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
			"Pass comma separated fieldnames (ex: 'date,type,lines') or the string 'firstline' to use the first line of the CSV",
			value.String,
			flag.Default("x,category,y"),
		),
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
		"--type": flag.New(
			"Type of graph to generate",
			value.StringEnum("point", "line", "grouped-bar", "stacked-bar"),
			flag.Default("line"),
			flag.Required(),
		),
		"--x-time-unit": flag.New(
			"X Time Unit only used for temporal types - see https://vega.github.io/vega-lite/docs/timeunit.html",
			value.StringEnum("utcyear", "utcyearmonthdate"), // TODO: more?
		),
		"--x-type": flag.New(
			"X type. See https://vega.github.io/vega-lite/docs/type.html",
			value.StringEnum("nominal", "quantitative", "temporal"),
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
			section.Flag(
				"--output",
				flag.HelpShort("Path to output file. Use DATEME as an alias for 'graph.<timestamp>'."),
				value.Path,
				flag.Default("DATEME.html"),
				flag.Required(),
			),
			section.Command(
				"graph",
				// "Graph from a 3-column CSV. First column is x-axis, and should be datetime. Second column is a string whose values are used to 'group'. Third column is a numeric column for the y-axis",
				"Various graphs from 3-column CSVs (x,category,y). Point, line, grouped-bar, stacked-bar ",
				graph,
				command.ExistingFlags(csvParseFlags),
				command.ExistingFlags(ioFlags),
				command.ExistingFlags(graphFlags),
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
