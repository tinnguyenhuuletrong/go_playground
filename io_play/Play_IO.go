package io_play

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

// word-count clone
//
//	main -f <filePath> - scan from file
//	main - scan from stdio
//
// Usage:
//   - gsutil cat gs://apache-beam-samples/shakespeare/kinglear.txt | ./bin/app -v
//   - ./bin/app -f <filePath>
//   - ./bin/app
//     typping. End with Ctrl + D
func Play_IO_WordCount() {
	var (
		r        io.Reader
		mode     string
		filePath string
		withLog  bool
	)

	flag.StringVar(&filePath, "f", "", "filePath to make a word scan. If not file it will sue stdin as an input")
	flag.BoolVar(&withLog, "v", false, "verbose log")
	flag.Parse()

	if filePath != "" {
		mode = "file"
	} else {
		mode = "stdin"
	}

	// create reader depend on mode
	if mode == "file" {
		fileReader, err := os.Open(filePath)
		if err != nil {
			log.Panic(err)
		}
		r = fileReader
	} else {
		r = os.Stdin
	}

	// enable/disable log
	if !withLog {
		log.SetFlags(0)
		log.SetOutput(io.Discard)
	}

	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanWords)

	count := 0
	for sc.Scan() {
		count += 1
		tmp := sc.Text()
		log.Printf("%d - %s\n", count, tmp)
	}

	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	fmt.Printf("%d\n", count)
}
