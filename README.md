# csv2json

This is a super simple utility used to convert a combined CSV of disparate data
types into an array of JSON objects containing the only the keys that were
populated in the CSV.

Really this was just so we could convert exported splunk CSVs to JSON to edit
and re-ingest via a HEC.

## Usage

```sh
❯ csv2json --help
Usage of csv2json:
      --in-file string    Input file (default "stdin")
      --out-file string   Output file (default "stdout")
```

## License

This project is [MIT-licensed](LICENSE). Use it as you'd like.

## Contributions

I can't possible imagine you'd want to contribute, but if you do, I'll happily
review any PRs.
