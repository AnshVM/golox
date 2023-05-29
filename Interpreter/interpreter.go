package Interpreter

import (
	"fmt"
	"strconv"

	"github.com/AnshVM/golox/Error"
	"github.com/AnshVM/golox/Parser"
	"github.com/AnshVM/golox/Tokens"
)

func Interpret(expr Parser.Expr) {
	result, err := Eval(expr)
	if err == nil {
		fmt.Printf("%v\n", result)
	}
}

func Eval(expr Parser.Expr) (any, error) {
	switch e := expr.(type) {
	case *Parser.Literal:
		return EvalLiteral(e), nil
	case *Parser.Binary:
		return EvalBinary(e)
	case *Parser.Grouping:
		return EvalGrouping(e)
	case *Parser.Unary:
		return EvalUnary(e)
	case *Parser.Conditional:
		return EvalConditional(e)
	}
	return nil, Error.ErrRuntimeError
}

func EvalLiteral(expr *Parser.Literal) any {
	return expr.Value
}

func EvalGrouping(expr *Parser.Grouping) (any, error) {
	evaluated, err := Eval(expr.Expression)
	if err != nil {
		return nil, err
	}
	return evaluated, nil
}

func EvalUnary(expr *Parser.Unary) (any, error) {
	right, err := Eval(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case Tokens.MINUS:
		rightVal, err := checkNumberOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}
		return -rightVal, nil

	case Tokens.BANG:
		return !isTruthy(right), nil
	}
	return nil, Error.ErrRuntimeError
}

func EvalConditional(conditional *Parser.Conditional) (any, error) {
	cond, err := Eval(conditional.Condition)
	if err != nil {
		return nil, Error.ErrRuntimeError
	}
	if cond == true {
		return Eval(conditional.Then)
	} else {
		return Eval(conditional.Else)
	}
}

func EvalBinary(binary *Parser.Binary) (any, error) {
	switch binary.Operator.Type {
	case Tokens.MINUS:
		left, right, err := EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left - right), nil

	case Tokens.SLASH:
		left, right, err := EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		if right == 0 {
			Error.ReportRuntimeError(binary.Operator, "Cannot divide by zero")
			return nil, Error.ErrRuntimeError
		}
		return (left / right), nil

	case Tokens.STAR:
		left, right, err := EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left * right), nil

	case Tokens.PLUS:
		left, right, err := EvalBinaryOperandsAny(binary)
		if err != nil {
			return nil, err
		}
		if isFloat32(right) && isFloat32(left) {
			return (left.(float32) + right.(float32)), nil
		}
		if isString(right) && isString(left) {
			return (left.(string) + right.(string)), nil
		}
		if isString(right) && isFloat32(left) {
			val, _ := left.(float32)
			return (strconv.FormatFloat(float64(val), 'f', 0, 32) + right.(string)), nil
		}
		// if isString(left) && isFloat32(right) {
		// 	val, _ := right.(float32)
		// 	return (left.(string) + strconv.FormatFloat(float64(val), 'f', 0, 32)), nil
		// }
		Error.ReportRuntimeError(binary.Operator, "Operands must strings or numbers")
		return nil, Error.ErrRuntimeError

	case Tokens.GREATER:
		left, right, err := EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left > right), nil

	case Tokens.GREATER_EQUAL:
		left, right, err := EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left >= right), nil

	case Tokens.LESS:
		left, right, err := EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left < right), nil

	case Tokens.LESS_EQAUL:
		left, right, err := EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left <= right), nil

	case Tokens.EQUAL_EQUAL:
		left, right, err := EvalBinaryOperandsAny(binary)
		if err != nil {
			return nil, err
		}
		return isEqual(left, right), nil

	case Tokens.BANG_EQUAL:
		left, right, err := EvalBinaryOperandsAny(binary)
		if err != nil {
			return nil, err
		}
		return isEqual(left, right), nil
	}

	// unreachable
	return nil, nil
}

func EvalBinaryOperandsAny(binary *Parser.Binary) (any, any, error) {
	right, err := Eval(binary.Right)
	if err != nil {
		return nil, nil, err
	}
	left, err := Eval(binary.Left)
	if err != nil {
		return nil, nil, err
	}
	return left, right, nil
}

func EvalBinaryOperandsNumber(binary *Parser.Binary) (float32, float32, error) {
	evalRight, err := Eval(binary.Right)
	if err != nil {
		return 0, 0, err
	}
	evalLeft, err := Eval(binary.Left)
	if err != nil {
		return 0, 0, err
	}
	right, err := checkNumberOperand(binary.Operator, evalRight)
	if err != nil {
		return 0, 0, err
	}
	left, err := checkNumberOperand(binary.Operator, evalLeft)
	if err != nil {
		return 0, 0, err
	}
	return left, right, nil
}

func checkNumberOperand(operator *Tokens.Token, right any) (float32, error) {
	if val, ok := right.(float32); ok {
		return val, nil
	} else {
		Error.ReportRuntimeError(operator, "Operand must be a number")
		return 0, Error.ErrRuntimeError
	}
}

// Handles the special case that one operand is boolean, and other is not
func isEqual(a any, b any) bool {
	if isBool(a) && !isBool(b) {
		b = isTruthy(b)
	}
	if isBool(b) && !isBool(a) {
		a = isTruthy(a)
	}
	return a == b
}

func isBool(x any) bool {
	return x == true || x == false
}

func isFloat32(val any) bool {
	_, ok := val.(float32)
	return ok
}

func isString(val any) bool {
	_, ok := val.(string)
	return ok
}

func isTruthy(val any) bool {
	if val == nil {
		return false
	}
	if boolVal, ok := val.(bool); ok {
		return boolVal
	}
	return true

}