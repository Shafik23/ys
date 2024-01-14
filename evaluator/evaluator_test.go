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
	env := object.NewEnvironment()

	// Evaluate the program.
	return Eval(ast, env)
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
		{`	let f = fn(x) {
				return x;
				x + 10;
			};
			f(10);`, 10},
		{`	let f = fn(x) {
				let result = x + 10;
				return result;
				return 10;
			};
			f(10);`, 20},
		{`	let f = fn(x) {
				let result = x + 10;
				return result;
				return 10;
			};
			f(10) + f(10);`, 40},
		{`	let f = fn(x) {
				let result = x + 10;
				return result;
				return 10;
			};
			let x = f(10);
			x;`, 20},
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
		{"foobar", "identifier not found: foobar"},
		{`"Hello" - "World!"`, "unknown operator: STRING - STRING"},
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

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		// Evaluate the input.
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	// Create a new function.
	input := "fn(x) { x + 2; };"

	// Evaluate the input.
	evaluated := testEval(input)

	// Cast the object to a function.
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	// Compare the parameters.
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	// Compare the body.
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("function has wrong body. Body=%q", fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		// Evaluate the input.
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	let newAdder = fn(x) {
		fn(y) { x + y };
	};

	let addTwo = newAdder(2);
	addTwo(2);
	`

	testIntegerObject(t, testEval(input), 4)
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!";`
	testStringObject(t, testEval(input), "Hello World!")
}

func testStringObject(t *testing.T, evaluated object.Object, expected string) {
	// Cast the object to a string.
	result, ok := evaluated.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", evaluated, evaluated)
		return
	}

	// Compare the value of the string.
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
		return
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got type INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
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
			// Cast the object to an error.
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			// Compare the error message.
			if errObj.Message != tt.expected {
				t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)

	// Cast the object to an array.
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	// Compare the length of the array.
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	// Compare the elements of the array.
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"let myArray = [1, 2, 3]; myArray[2];", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i];", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
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
