# Ys

Ys (pronounced "Wise") is a simple language written in Go. 

The project is structured as follows:

- `main.go`: This is the entry point of the application.
- `ast/`: This directory contains files related to the Abstract Syntax Tree (AST) that represents the structure of Ys programs.
  - `ast.go`: Defines the structures of the AST.
  - `ast_test.go`: Contains unit tests for the AST.
- `evaluator/`: This directory contains files related to the evaluation of Ys programs.
  - `builtins.go`: Defines built-in functions.
  - `evaluator.go`: Contains the logic for evaluating nodes of the AST.
  - `evaluator_test.go`: Contains unit tests for the evaluator.
- `lexer/`: This directory contains files related to the lexical analysis of Ys programs.
  - `lexer.go`: Contains the logic for breaking down Ys programs into tokens.
  - `lexer_test.go`: Contains unit tests for the lexer.
- `object/`: This directory contains files related to the objects that Ys programs manipulate.
  - `environment.go`: Defines the environment in which Ys programs run.
  - `object.go`: Defines the structures of objects.
  - `object_test.go`: Contains unit tests for the objects.
- `parser/`: This directory contains files related to the parsing of Ys programs.
  - `parser.go`: Contains the logic for parsing tokens into an AST.
  - `parser_test.go`: Contains unit tests for the parser.
  - `parser_tracing.go`: Contains utility functions for tracing the parser's progress (useful for debugging).
- `repl/`: This directory contains files related to the Read-Eval-Print Loop (REPL) of Ys.
  - `repl.go`: Contains the logic for the REPL.
- `token/`: This directory contains files related to the tokens that the lexer produces.
  - `token.go`: Defines the types of tokens.

To build the project, run the `build.sh` script. This will produce an executable that you can run to start the REPL and interact with the Ys language.