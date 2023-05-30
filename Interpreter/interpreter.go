package Interpreter

import (
	"fmt"
	"strconv"

	"github.com/AnshVM/golox/Environment"
	"github.com/AnshVM/golox/Error"
	"github.com/AnshVM/golox/Parser"
	"github.com/AnshVM/golox/Tokens"
)

type Interpreter struct {
	Env *Environment.Environment
}

func (i *Interpreter) Interpret(stmts []Parser.Stmt) error {
	for _, stmt := range stmts {
		err := i.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) Exec(stmt Parser.Stmt) error {
	switch s := stmt.(type) {
	case Parser.Expression:
		return i.ExecExpressionStmt(&s)
	case Parser.Print:
		return i.ExecPrintStmt(&s)
	case Parser.Var:
		return i.ExecVarStmt(&s)
	case Parser.Block:
		return i.ExecBlockStmt(&s)
	}
	return nil
}

func (i *Interpreter) ExecExpressionStmt(stmt *Parser.Expression) error {
	_, err := i.Eval(stmt.Expression)
	return err
}

func (i *Interpreter) ExecPrintStmt(stmt *Parser.Print) error {
	result, err := i.Eval(stmt.Expression)
	if err == nil {
		fmt.Printf("%v\n", result)
	}
	return err
}

func (i *Interpreter) ExecVarStmt(stmt *Parser.Var) error {
	if stmt.Initializer != nil {
		value, err := i.Eval(stmt.Initializer)
		if err == nil {
			i.Env.Define(stmt.Name, value)
		}
		return err
	}
	i.Env.Define(stmt.Name, nil)
	return nil
}

func (i *Interpreter) ExecBlockStmt(stmt *Parser.Block) error {
	var executeBlock = func(statements []Parser.Stmt, env *Environment.Environment) error {
		prev := i.Env
		defer func() {
			i.Env = prev
		}()
		i.Env = env
		var err error
		for _, stmt := range statements {
			err = i.Exec(stmt)
			if err != nil {
				return err
			}
		}
		return nil
	}
	err := executeBlock(stmt.Statements, &Environment.Environment{Enclosing: i.Env, Values: map[string]any{}})
	return err
}

func (i *Interpreter) Eval(expr Parser.Expr) (any, error) {
	switch e := expr.(type) {
	case *Parser.Literal:
		return i.EvalLiteral(e), nil
	case *Parser.Binary:
		return i.EvalBinary(e)
	case *Parser.Grouping:
		return i.EvalGrouping(e)
	case *Parser.Unary:
		return i.EvalUnary(e)
	case *Parser.Conditional:
		return i.EvalConditional(e)
	case *Parser.Variable:
		return i.EvalVariable(e), nil
	case *Parser.Assign:
		return i.EvalAssign(e)
	}
	return nil, Error.ErrRuntimeError
}

func (i *Interpreter) EvalAssign(expr *Parser.Assign) (any, error) {
	value, err := i.Eval(expr.Value)
	if err != nil {
		return nil, err
	}
	i.Env.Assign(expr.Name, value)
	return value, nil
}

func (i *Interpreter) EvalLiteral(expr *Parser.Literal) any {
	return expr.Value
}

func (i *Interpreter) EvalGrouping(expr *Parser.Grouping) (any, error) {
	evaluated, err := i.Eval(expr.Expression)
	if err != nil {
		return nil, err
	}
	return evaluated, nil
}

func (i *Interpreter) EvalUnary(expr *Parser.Unary) (any, error) {
	right, err := i.Eval(expr.Right)
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

func (i *Interpreter) EvalVariable(expr *Parser.Variable) any {
	return i.Env.Get(expr.Name)
}

func (i *Interpreter) EvalConditional(conditional *Parser.Conditional) (any, error) {
	cond, err := i.Eval(conditional.Condition)
	if err != nil {
		return nil, Error.ErrRuntimeError
	}
	if cond == true {
		return i.Eval(conditional.Then)
	} else {
		return i.Eval(conditional.Else)
	}
}

func (i *Interpreter) EvalBinary(binary *Parser.Binary) (any, error) {
	switch binary.Operator.Type {
	case Tokens.MINUS:
		left, right, err := i.EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left - right), nil

	case Tokens.SLASH:
		left, right, err := i.EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		if right == 0 {
			Error.ReportRuntimeError(binary.Operator, "Cannot divide by zero")
			return nil, Error.ErrRuntimeError
		}
		return (left / right), nil

	case Tokens.STAR:
		left, right, err := i.EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left * right), nil

	case Tokens.PLUS:
		left, right, err := i.EvalBinaryOperandsAny(binary)
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
		left, right, err := i.EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left > right), nil

	case Tokens.GREATER_EQUAL:
		left, right, err := i.EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left >= right), nil

	case Tokens.LESS:
		left, right, err := i.EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left < right), nil

	case Tokens.LESS_EQAUL:
		left, right, err := i.EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left <= right), nil

	case Tokens.EQUAL_EQUAL:
		left, right, err := i.EvalBinaryOperandsAny(binary)
		if err != nil {
			return nil, err
		}
		return isEqual(left, right), nil

	case Tokens.BANG_EQUAL:
		left, right, err := i.EvalBinaryOperandsAny(binary)
		if err != nil {
			return nil, err
		}
		return isEqual(left, right), nil
	}

	// unreachable
	return nil, nil
}

func (i *Interpreter) EvalBinaryOperandsAny(binary *Parser.Binary) (any, any, error) {
	right, err := i.Eval(binary.Right)
	if err != nil {
		return nil, nil, err
	}
	left, err := i.Eval(binary.Left)
	if err != nil {
		return nil, nil, err
	}
	return left, right, nil
}

func (i *Interpreter) EvalBinaryOperandsNumber(binary *Parser.Binary) (float32, float32, error) {
	evalRight, err := i.Eval(binary.Right)
	if err != nil {
		return 0, 0, err
	}
	evalLeft, err := i.Eval(binary.Left)
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
