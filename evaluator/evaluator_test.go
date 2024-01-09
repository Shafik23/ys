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

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	// Iterate over each test case.
	for _, tt := range tests {
		// Evaluate the input.
		evaluated := testEval(tt.input)

		// Check if the expected value is an integer.
		integer, ok := tt.expected.(int)
		if ok {
			// Compare the value of the integer.
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			// Compare the value of the null.
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, evaluated object.Object) bool {
	if evaluated != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", evaluated, evaluated)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
		if (10 > 1) {
			if (10 > 1) {
				return 10;
			}

			return 1;
		}
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input          string
		expectedErrMsg string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{`
		if (10 > 1) {
			if (10 > 1) {
				return true + false;
			}

			return 1;
		}
		`, "unknown operator: BOOLEAN + BOOLEAN"},
	}

	// Iterate over each test case.
	for _, tt := range tests {
		// Evaluate the input.
		evaluated := testEval(tt.input)

		// Cast the object to an error.
		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("no error object returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}

		// Compare the error message.
		if errObj.Message != tt.expectedErrMsg {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedErrMsg, errObj.Message)
		}
	}
}
