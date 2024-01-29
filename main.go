package main

import (
	"fmt"
	"os"
	"strings"
)

type PrinterFunc func(format string, a ...any) (n int, err error)

// Keys are strings to be more generic - but actually this is just a number: "1" is mask file, "2" is next file and so on
type Mapping map[string]*PerFile

func main() {
	const MainFile = "0"

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%s <mask file 1> [file 2] [file 3] ...", os.Args[0])
		os.Exit(1)
	}

	// parse flags
	cycle := true
	args := make([]string, 0)
	for _, v := range os.Args {
		if strings.HasPrefix(v, "--") {
			if strings.ToLower(v) == "--no-cycle" {
				cycle = false
			}
			continue
		}

		args = append(args, v)
	}

	mapping := InitMapping(args[1:], cycle)
	mapping.Printer(MainFile, fmt.Printf)
	mapping.Close()
}

func (m Mapping) Printer(mainFile string, printerFunc PrinterFunc) {
	for {
		m.Next()
		if m[mainFile].eof {
			break
		}

		token := m[mainFile].Get()

		if v, exists := m[token]; exists && !v.eof {
			printerFunc("%s\n", v.Get())
		} else {
			printerFunc("%s\n", token)
		}
	}
}

func InitMapping(args []string, shouldCycle bool) Mapping {
	ret := make(Mapping)

	cycle := shouldCycle

	for k, v := range args {
		if k == 0 {
			// Original file must never be cyclic or you would never stop.
			// Infinite sequence generation and lazy usage of first X lines using `head` can be a topic for later.
			cycle = false
		} else {
			cycle = shouldCycle
		}

		f, err := NewPerFile(v, cycle)
		if err != nil {
			if k == 0 {
				fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
				os.Exit(1)
			} else {
				fmt.Fprintf(os.Stderr, "Error opening file: %s, skipping\n", err)
				continue
			}
		}

		ret[fmt.Sprintf("%d", k)] = f
	}

	return ret
}

func (m Mapping) Close() {
	for _, v := range m {
		err := v.file.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error closing file: %s\n", err)
		}
	}
}

func (m Mapping) Next() {
	for _, v := range m {
		v.Next()
	}
}
