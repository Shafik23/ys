// File: repl/repl.go

package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/shafik23/ys/lexer"
	"github.com/shafik23/ys/token"
)

// Start launches the REPL.
func Start() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf(">> ")
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

			fmt.Printf("%+v\n", tok)
		}
	}
}
