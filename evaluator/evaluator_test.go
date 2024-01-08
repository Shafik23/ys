package evaluator

import (
	"testing"

	"github.com/shafik23/ys/lexer"
	"github.com/shafik23/ys/object"
	"github.com/shafik23/ys/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		// {"5 + 5 + 5 + 5 - 10", 10},
		// {"2 * 2 * 2 * 2 * 2", 32},
		// {"-50 + 100 + -50", 0},
		// {"5 * 2 + 10", 20},
		// {"5 + 2 * 10", 25},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	// Create a new lexer and parser for each test case.
	l := lexer.New(input)
	p := parser.New(l)

	ast := p.ParseProgram()

	// Evaluate the program.
	return Eval(ast)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	// Cast the object to an integer.
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	// Compare the value of the integer.
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	// Return true if the test passed.
	return true
}
