# covid-19

This site used to exist in https://covid-19.sandbox.avm99963.com/.

The copy in this archive is a copy of the files served by Apache to end users
(index.html + plus the output files). The original code that generated these
graphs daily can be browsed in the [project's original repository][git].

## Modifications

- In the live version of covid-19.sandbox.avm99963.com we already recovered a
  previous output image in SVG and PNG formats (I think from the Web Archive),
  since the cron service that updated the outputs was giving us blank graphs
  due to the data becoming unavailable.
    - The blank SVG and PNG images were left in the webserver with names
      `graph_20241006.svg` and `graph_20241006.png`, but have been deleted in
      this archive.
- The SVG file has been compressed with [svgomg][svgomg].

[git]: https://gerrit.avm99963.com/plugins/gitiles/covid19
[svgomg]: https://jakearchibald.github.io/svgomg/
