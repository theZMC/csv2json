package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strconv"

	"github.com/spf13/pflag"
)

func main() {
	inFile := pflag.String("in-file", "stdin", "Input file")
	outFile := pflag.String("out-file", "stdout", "Output file")
	pflag.Parse()

	if *inFile == "" || *outFile == "" {
		pflag.PrintDefaults()
		os.Exit(1)
	}

	var out io.Writer = os.Stdout
	var err error

	if *outFile != "stdout" {
		out, err = os.Create(*outFile)
		if err != nil {
			slog.Error("Could not create output file",
				slog.String("error", err.Error()),
				slog.String("output-file", *outFile),
			)
			os.Exit(1)
		}
	}

	var in io.Reader = os.Stdin
	if *inFile != "stdin" && *inFile != "-" {
		in, err = os.Open(*inFile)
		if err != nil {
			slog.Error("Could not open input file",
				slog.String("error", err.Error()),
				slog.String("input-file", *inFile),
			)
			os.Exit(1)
		}
	}

	writeJSON(out, toJSON(fromCSV(in)))

	os.Exit(0)
}

func fromCSV(r io.Reader) ([]string, <-chan []string) {
	reader := csv.NewReader(r)
	ch := make(chan []string, 1)
	headers, err := reader.Read()
	if err != nil {
		slog.Error("Could not read CSV file",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	go func() {
		defer close(ch)
		for record, err := reader.Read(); err != io.EOF; record, err = reader.Read() {
			if err != nil {
				slog.Error("Could not read CSV file",
					slog.String("error", err.Error()),
				)
				os.Exit(1)
			}
			ch <- record
		}
	}()

	return headers, ch
}

var numRGX = regexp.MustCompile(`^\d*(\.)?\d+$`)

func toJSON(headers []string, records <-chan []string) <-chan map[string]any {
	ch := make(chan map[string]any, 1)

	go func() {
		defer close(ch)
		for record := range records {
			recordMap := make(map[string]any)
			for i, header := range headers {
				if record[i] != "" {
					if numRGX.MatchString(record[i]) {
						num, _ := strconv.Atoi(record[i])
						recordMap[header] = num
					} else {
						recordMap[header] = record[i]
					}
				}
			}
			ch <- recordMap
		}
	}()

	return ch
}

func writeJSON(w io.Writer, ch <-chan map[string]any) {
	enc := json.NewEncoder(w)

	w.Write([]byte("["))
	firstRecord := true
	for record := range ch {
		if !firstRecord {
			w.Write([]byte(","))
		} else {
			firstRecord = false
		}
		if err := enc.Encode(record); err != nil {
			slog.Error("Could not write JSON to file",
				slog.String("error", err.Error()),
			)
			os.Exit(1)
		}
	}
	w.Write([]byte("]"))
}
