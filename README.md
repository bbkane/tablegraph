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

Line chart: [Open the Chart in the Vega Editor](https://vega.github.io/editor/#/url/vega-lite/N4IgJAzgxgFgpgWwIYgFwhgF0wBwqgegIDc4BzJAOjIEtMYBXAI0poHsDp5kTykBaADZ04JAKyUAVhDYA7EABoQAEzjQATjRyZ289AGVMbKAGsABDk1Q1ZtgDMzYswBU4sMwGE2CHElk0bNlJ1FxoEOEpFFSRMFFRQBnVBNGjYziNTCEooCGIQAF8lZHUTFOFZOCi4WSg2ZRpZMjRQAA9mkDsAwWUU5RjKpUwATxxK9ExEHDZ1JGTCkCH2zrhulMsaayjh0ZSARwY-HVidUgKlWsFppa6e9AghhCY2ZMGRsZBZbwbZgvz8oA)

Let's see if I can get mouseover number, increase the chart width, isolate a line

### Chart Width

https://vega.github.io/vega-lite/docs/size.html#specifying-responsive-width-and-height

[Open the Chart in the Vega Editor](https://vega.github.io/editor/#/url/vega-lite/N4IgJAzgxgFgpgWwIYgFwhgF0wBwqgegIDc4BzJAOjIEtMYBXAI0poHsDp5kTykBaADZ04JAKyUAVhDYA7EABoQAEzjQATjRyZ289AGVMbKAGsABDk1Q1ZtgDMzYswBU4sMwGE2CHElk0bNlJ1FxoEOEpFFSRMFFRQBnVBNGjYziNTCEooCGIQAF8lZHUTFOFZOCj4GjIsFKg5WJoK9SiAdxplenrGpGa4VqU4WQblZrI0UAAPSZA7AMFlFOUYyqVMAE8cSvRMRBw2dSRkwpAN2fm4RZTLGmsoze2UgEcGPx1YnVICpQbBQ4uCyW6AgGwQTDYyXWWx2IFk3maxwK+XyQA)

There's also the "autosize" param, but I'll play with that later I guess :)

### View Number on mouseover

https://vega.github.io/vega-lite/docs/tooltip.html

```
  "mark": {"type": "line", "tooltip": true },
```

[Open the Chart in the Vega Editor](https://vega.github.io/editor/#/url/vega-lite/N4IgJAzgxgFgpgWwIYgFwhgF0wBwqgegIDc4BzJAOjIEtMYBXAI0poHsDp5kTykBaADZ04JAKyUAVhDYA7EABoQAEzjQATjRyZ289AGVMbKAGsABDk1Q1ZtgDMzYswBU4sMwGE2CHElk0bNlJ1FxoEOEpFFSRMFFRQBnVBNGjYziNTCEooCGIQAF8lZHUTNFBMAE8cOBThWRqlIzZBHRw0THUGOEKMOBoyLBSoOViaevUogHcaZXohkaQxuAmlOFlh5TGyMpAADx27AMFlFOUYhpBK6pTMRBw2dSRknoqDo5P0SxprKKua9AAjgw-DpYjpSAUlMNBA83nBjikIBUEExmr8qv8QLJvGMngV8vkgA)

This also shows the line names. It *does* seem to interact poorly with the zoom... might have to look up autosize again.

[Open the Chart in the Vega Editor](https://vega.github.io/editor/#/url/vega-lite/N4IgJAzgxgFgpgWwIYgFwhgF0wBwqgegIDc4BzJAOjIEtMYBXAI0poHsDp5kTykBaADZ04JAKyUAVhDYA7EABoQAEzjQATjRyZ289AGVMbKAGsABDk1Q1ZtgDMzYswBU4sMwGE2CHElk0bNlJ1FxoEOEpFFSRMFFRQBnVBNGjYziNTCEooCGIQAF8lZHUTNFBMAE8cOBThWRqlIzZBHRw0THUGOCUcNhpZTHbOuEKMOBoyLBSoOVj+uHUogHcaZXpp2aR5xaU4WRnlfrIykAAPE7sAwWUU5RiGkErqlMxEXvUkZNGKi6ub9EsNGsUSeNXQAEcGH4dLEdKQCkoZoI2It4iBLnBrikIBUEExmiCqmCQLJvP1PgV8vkgA)

### Isolate line

https://vega.github.io/vega-lite/examples/interactive_legend.html

https://vega.github.io/vega-lite/examples/interactive_line_hover.html

This makes it so when I click the line I get different opacity:

[Open the Chart in the Vega Editor](https://vega.github.io/editor/#/url/vega-lite/N4IgJAzgxgFgpgWwIYgFwhgF0wBwqgegIDc4BzJAOjIEtMYBXAI0poHsDp5kTykBaADZ04JAKyUAVhDYA7EABoQAEzjQATjRyZ289AGVMbKAGsABDk1Q1ZtgDMzYswBU4sMwGE2CHElk0bNlJ1FxoEOEpFFSRMFFRQBnVBNGjYziNTCEooCGIQAF8lZHUTNFBMAE8cOBThWRqlIzZBHRw0THUGOCUcNhpZTHbOuEKMOBoyLBSoOVj+uHUogHcaZXpp2aR5xaU4WRnlfrIykAAPE7sAwWUU5RiGkErqlMxEXvUkZNGKi6ub9EsNGsUSeNXQAEcGH4dLEdKQCkoZoI2It4iBLnBrikIBUEExmiCqmCQLJvP1PgiQGxfFA6D80TNZIcdHITr4PggUjAggsosRPl00ABGUb8wSC1AABkoItG7KQCAgaAA2qBZAridzglEIJi3IM0aCUr1+oMlBjrkrUMqQDi8QSALr5fJOoA)

If I add `"bind": "legend"`, then i can click the legend, but not the line anymore. This is worth it to me. I'd like to be able to select multiple lines like I can with plotly, but that's fine... IT WORKS if I use "shift" click to select multiple legend items.

[Open the Chart in the Vega Editor](https://vega.github.io/editor/#/url/vega-lite/N4IgJAzgxgFgpgWwIYgFwhgF0wBwqgegIDc4BzJAOjIEtMYBXAI0poHsDp5kTykBaADZ04JAKyUAVhDYA7EABoQAEzjQATjRyZ289AGVMbKAGsABDk1Q1ZtgDMzYswBU4sMwGE2CHElk0bNlJ1FxoEOEpFFSRMFFRQBnVBNGjYziNTCEooCGIQAF8lZHUTNFBMAE8cOBThWRqlIzZBHRw0THUGOCUcNhpZTHbOuEKMOBoyLBSoOVj+uHUogHcaZXpp2aR5xaU4WRnlfrIykAAPE7sAwWUU5RiGkErqlMxEXvUkZNGKi6ub9EsNGsUSeNXQAEcGH4dLEdKQCkoZoI2It4iBLnBrikIBUEExmiCqmCQLJvP1PgiQGxfFA6D80TNZIcdHITr4PggUjAggsosRPl00ABGUb8wSC1AABkoItG7KQCAgaAA2qBZAridzglEmP1-iBBOQ9jclBBMW5BmjQSlev1BkoMdclahlSAcXiCQBdfL5b1AA)

With title:

[Open the Chart in the Vega Editor](https://vega.github.io/editor/#/url/vega-lite/N4IgJAzgxgFgpgWwIYgFwhgF0wBwqgegIDc4BzJAOjIEtMYBXAI0poHsDp5kTykBaADZ04JAKyUAVhDYA7EABoQAEzjQATjRyZ289AGVMbKAGsABDk1Q1ZtgDMzYswBU4sMwGE2CHElk0bNlJ1FxoEOEpFFSRMFFRQBnVBNGjYziNTCEooCGIQAF8lZHUTNFBMAE8cOBThWRqlIzZBHRw0THUGOCUcNhpZTHbOuEKMOBoyLBSoOVj+uHUogHcaZXpp2aR5xaU4WRnlfrIykAAPE7sAwWUU5RiGkErqlMxEXvUkZNGKi6ub9EsNGsUSeNXQAEcGH4dLEdKQCkoZoI2It4iBLnBrikIBUEExmiCqmCQLJvP1PgiQGxfFA6D80TNZIcdHITr4PggUjAggsosRPl00ABGUb8wSC1AABkoItGMMEYPKcFOgwMGRMEEp7KQCE1qAA2qBZDridzglEmP1-iAFWQ9jclBBMW5VeUiSlev1BkoMdc9fqQDi8QSALr5fJhoA)

Embed plot: https://vega.github.io/vega-lite/usage/embed.html

Embedding works pretty well. I feel confident I can write the code to inline stuff in the chart.

Bar charts: https://vega.github.io/vega-lite/docs/bar.html

### Grouped Bar chart Over Time

grouped bar chart: https://vega.github.io/vega-lite/docs/bar.html#grouped-bar-chart-with-facet

For example, Git (or expenses) over time. Could also be stacked? Or an areachart
