package main

import (
	"os"

	"github.com/bbkane/warg"
	"github.com/bbkane/warg/command"
	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/section"
	"github.com/bbkane/warg/value"
)

func main() {

	htmlTitleFlag := flag.New(
		flag.HelpShort("HTML title. Flag ignored when --format != 'html'"),
		value.String,
		flag.Default("tablegraph output"),
	)

	graphFlags := flag.FlagMap{
		"--div-title": flag.New(
			flag.HelpShort("Title of div when --format is 'div'"),
			value.String,
		),
		"--fieldnames": flag.New(
			flag.HelpShort("Field names"),
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
				"html-bottom",
				command.HelpShort("Print bottom of generated HTML file to stdout"),
				command.DoNothing,
			),
			section.Command(
				"html-top",
				command.HelpShort("Print top of generated HTML file to stdout"),
				command.DoNothing,
				command.ExistingFlag("--html-title", htmlTitleFlag),
			),
			section.Command(
				"table",
				command.HelpShort("Make HTML table"),
				command.DoNothing,
				command.ExistingFlags(graphFlags),
				command.Flag(
					"--page-length",
					flag.HelpShort("Entries in table before needing to click next page"),
					value.Int,
				),
				command.HelpLong("NOTE: columns should not have a `.` in the title. See https://datatables.net/forums/discussion/69257/data-with-a-in-the-name-makes-table-creation-fail#latest"),
			),
			section.Command(
				"timechart",
				command.HelpShort("Line graph. First column is x-axis, and should be datetime. Second column can optionally be a string whose values are used to 'group' the numeric columns and create different lines in the chart. Other (numeric) columns are y-axes."),
				command.DoNothing,
				command.ExistingFlags(graphFlags),
				command.Flag(
					"--chart-title",
					flag.HelpShort("Chart title"),
					value.String,
				),
				command.Flag(
					"--xaxis-title",
					flag.HelpShort("X-Axis Title"),
					value.String,
				),
				command.Flag(
					"--yaxis-title",
					flag.HelpShort("Y-Axis Title"),
					value.String,
				),
			),
		),
	)
	app.MustRun(os.Args, os.LookupEnv)
}
