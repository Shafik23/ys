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
		{"-5", -5},
		{"-10", -10},
		{"-100", -100},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
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

func TestEvailBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"5 < 10", true},
		{"5 > 10", false},
		{"5 == 5", true},
		{"5 != 5", false},
		{"5 == 10", false},
		{"5 != 10", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(5 < 10) == true", true},
		{"(5 < 10) == false", false},
		{"(5 > 10) == true", false},
		{"(5 > 10) == false", true},
	}

	// Iterate over each test case.
	for _, tt := range tests {
		// Evaluate the input.
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, evaluated object.Object, expected bool) {
	// Cast the object to a boolean.
	result, ok := evaluated.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", evaluated, evaluated)
		return
	}

	// Compare the value of the boolean.
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	// Iterate over each test case.
	for _, tt := range tests {
		// Evaluate the input.
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}
