// File: repl/repl.go

package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/shafik23/ys/lexer"
	"github.com/shafik23/ys/parser"
)

const PROMPT = ">> "

// Start launches the REPL, taking input from an io.Reader and sending output to an io.Writer.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Something UnWise happened:\n")
	io.WriteString(out, " parser errors:\n")

	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
