package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nsd20463/durationsort"
)

func main() {
	var err error
	// read in lines from stdin, duration-sort them, write them out
	scan := bufio.NewScanner(os.Stdin)
	var lines []string
	for scan.Scan() {
		lines = append(lines, scan.Text()) // TODO use Bytes() to avoid the alloc and mempcy. Or get more clever with the buffers and have strings which point into the []bytes
	}
	if err = scan.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}
	scan = nil

	durationsort.Strings(lines)

	out := bufio.NewWriter(os.Stdout)
	err = nil
	for _, line := range lines {
		_, err = out.WriteString(line)
		if err != nil {
			break
		}
		err = out.WriteByte('\n')
		if err != nil {
			break
		}
	}
	if err == nil {
		err = out.Flush()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing stdout: %v\n", err)
		os.Exit(1)
	}
}
