# tablegraph

Create graphs from CSVs

The Go files are experimental. The Python version actually works.

Simple install via symlink:

```bash
ln -s $HOME/Git/tablegraph/tablegraph.py $HOME/bin/tablegraph
```

# Go + Vega-Lite Rewrite

## Problems with Python + Plotly

### Packaging.

I don't feel comfortable with Python's packaging ecoystem. I want to `brew install <tool>` reliably and I don't know how to get that with Python + dependencies.

Options: 

- learn how to do Python + dependencies, probably by following Simon Willerson's work on DataSette. This seems to involve a somewhat complicated set of tools, but he's done most of the work. If I do this, I'm worried that the components will change over time and I'll have to keep on top of that. That does get me the `plotly` package instead of manually writing the JSON.
- Use a single Python file and only the stdlib. This is what I do now (and it's honestly not too bad), but I want more complext help than `argparse` can provide. TODO: see how I can do that
- Switch to Go. Go has a great packaging system I understand, and I could use my CLI library warg.

### Program Correctness / Tooling

Typs + Lints

Python has optional type hints that I use, but because I don't really run MyPy consistently, they function more like docs... that goes for a lot of Python tooling. All the best editors for Python are proprietary (Pylance, PyCharm).

Due to my single file choice, it's getting more difficult to add different types of graphs

Options:

- Python + dependencies
- Switch to Go

### Plotly JSON

Plotly uses fairly complex JSON to build charts. I want to add more types of charts, and I'm not looking forward to learning the right JSON for histograms, etc.

Options:

- Python + dependencies. This gets me the plotly library, which would help, but I'd still have to learn it.
- Switch to Vega, which uses CSV but might not be powerful enough.

## Vega-Lite unknowns

Vega/Vega-Lite *already* use CSV to make graphs

Let's seen if I can recreate what plotly gives me in `cashflow`... multiline chart, zoom, show nubmer at point, space filling chart, isolate a line

Also want histograms, bar and column charts.
