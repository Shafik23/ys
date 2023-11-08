// File: repl/repl.go

package repl

const PROMPT = ">> "

// Start initiates the REPL loop.
// func Start(in io.Reader, out io.Writer) {
// 	scanner := bufio.NewScanner(in)

// 	for {
// 		fmt.Fprintf(out, PROMPT)
// 		scanned := scanner.Scan()
// 		if !scanned {
// 			return
// 		}

// 		line := scanner.Text()
// 		l := lexer.NewLexer(line)

// 		for tok, _ := l.NextToken(); tok.Type != token.EOF; tok, _ = l.NextToken() {
// 			fmt.Fprintf(out, "%+v\n", tok)
// 		}
// 	}
// }
