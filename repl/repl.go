// File: repl/repl.go

package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/shafik23/ys/lexer"
	"github.com/shafik23/ys/token"
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
		l := lexer.NewLexer(line)

		for {
			tok := l.NextToken()
			if tok.Type == token.EOF {
				break
			}

			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
